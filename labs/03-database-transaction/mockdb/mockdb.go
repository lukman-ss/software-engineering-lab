// Package mockdb provides a minimal in-memory database/sql driver for educational purposes.
//
// Concurrency Model:
//   - Snapshot isolation for reads (BEGIN clones committed state)
//   - Atomic INSERT ... ON CONFLICT with global lock (checks committed state during INSERT)
//   - Commit without lock (tx already contains full snapshot with INSERT applied)
//
// This simulates PostgreSQL-like behavior where ON CONFLICT checks committed state.
package mockdb

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

type row = map[string]any

type tableState struct {
	rows []row
}

type DB struct {
	mu        sync.Mutex
	committed map[string]*tableState
}

func newDB() *DB {
	return &DB{
		committed: map[string]*tableState{
			"payments":            {},
			"orders":              {},
			"wallet_transactions": {},
			"invoices":            {},
			"outbox_events":       {},
			"processed_events":    {},
			"commissions":         {},
		},
	}
}

func cloneTables(m map[string]*tableState) map[string]*tableState {
	out := make(map[string]*tableState, len(m))
	for k, t := range m {
		rows := make([]row, len(t.rows))
		for i, r := range t.rows {
			nr := make(row, len(r))
			for c, v := range r {
				nr[c] = v
			}
			rows[i] = nr
		}
		out[k] = &tableState{rows: rows}
	}
	return out
}

type staticDriver struct{}

func (staticDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("use sql.OpenDB(mockdb.Connector())")
}

type connector struct{ db *DB }

func Connector() driver.Connector { return &connector{db: newDB()} }
func NewDB() *sql.DB            { return sql.OpenDB(Connector()) }

func (c *connector) Connect(context.Context) (driver.Conn, error) { return &conn{db: c.db}, nil }
func (c *connector) Driver() driver.Driver                        { return staticDriver{} }

type conn struct {
	db *DB
	tx map[string]*tableState
}

func (c *conn) Prepare(query string) (driver.Stmt, error)                 { return &stmt{c: c, query: query}, nil }
func (c *conn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	return &stmt{c: c, query: query}, nil
}
func (c *conn) Close() error               { return nil }
func (c *conn) Ping(context.Context) error { return nil }

func (c *conn) Begin() (driver.Tx, error) { return c.begin() }
func (c *conn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return c.begin()
}

func (c *conn) begin() (driver.Tx, error) {
	if c.tx != nil {
		return nil, errors.New("transaction already open on connection")
	}
	// Hold lock for entire transaction duration for correctness
	// This serializes transactions but ensures ON CONFLICT works correctly
	c.db.mu.Lock()
	c.tx = cloneTables(c.db.committed)
	return &txn{c: c, open: true}, nil
}

type txn struct {
	c    *conn
	open bool
}

func (t *txn) IsOpen() bool {
	return t.open
}

func (t *txn) Commit() error {
	c := t.c
	// Unlock AFTER commit
	defer c.db.mu.Unlock()
	for table, txState := range c.tx {
		c.db.committed[table] = &tableState{rows: append([]row{}, txState.rows...)}
	}
	t.open = false
	t.c.tx = nil
	return nil
}

func (t *txn) Rollback() error {
	c := t.c
	defer c.db.mu.Unlock()
	t.open = false
	t.c.tx = nil
	return nil
}

func (c *conn) CheckConstraintViolation(query string, vals []driver.Value) error {
	return nil
}

func (c *conn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.exec(query, toVals(args))
}

func (c *conn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.query(query, toVals(args))
}

func (c *conn) exec(query string, vals []driver.Value) (driver.Result, error) {
	upper := strings.ToUpper(strings.TrimSpace(query))

	// Within transaction: no lock needed (held from Begin)
	if c.tx != nil && strings.HasPrefix(upper, "INSERT") && strings.Contains(upper, "ON CONFLICT") {
		table := extractTableFromInsert(upper)
		if existing, ok := c.db.committed[table]; ok {
			for _, existingRow := range existing.rows {
				found := false
				for _, txRow := range c.tx[table].rows {
					if equalVal(existingRow, txRow) {
						found = true
						break
					}
				}
				if !found {
					c.tx[table].rows = append(c.tx[table].rows, existingRow)
				}
			}
		}
		return doInsertReturning(c.tx, query, vals)
	}

	if c.tx != nil {
		return doExec(c.tx, query, vals)
	}
	c.db.mu.Lock()
	defer c.db.mu.Unlock()
	return doExec(c.db.committed, query, vals)
}

// extractTableFromInsert extracts table name from INSERT query
func extractTableFromInsert(q string) string {
	upper := strings.ToUpper(q)
	idxInto := strings.Index(upper, "INSERT INTO ")
	if idxInto < 0 {
		return ""
	}
	afterInsert := strings.TrimSpace(q[idxInto+len("INSERT INTO "):])
	tableEndIdx := strings.Index(afterInsert, " ")
	parenIdx := strings.Index(afterInsert, "(")
	if parenIdx >= 0 && (tableEndIdx < 0 || parenIdx < tableEndIdx) {
		tableEndIdx = parenIdx
	}
	if tableEndIdx < 0 {
		tableEndIdx = len(afterInsert)
	}
	return strings.TrimSpace(afterInsert[:tableEndIdx])
}

func (c *conn) query(query string, vals []driver.Value) (driver.Rows, error) {
	if c.tx != nil {
		return doQuery(c.tx, query, vals)
	}
	c.db.mu.Lock()
	defer c.db.mu.Unlock()
	return doQuery(c.db.committed, query, vals)
}

type stmt struct {
	c     *conn
	query string
}

func (s *stmt) Close() error  { return nil }
func (s *stmt) NumInput() int { return -1 }
func (s *stmt) Exec(args []driver.Value) (driver.Result, error) {
	return s.c.exec(s.query, args)
}
func (s *stmt) Query(args []driver.Value) (driver.Rows, error) {
	return s.c.query(s.query, args)
}

func toVals(args []driver.NamedValue) []driver.Value {
	v := make([]driver.Value, len(args))
	for i, a := range args {
		v[i] = a.Value
	}
	return v
}

type result struct{ affected int64 }

func (r *result) LastInsertId() (int64, error) { return 0, errors.New("last insert id not supported") }
func (r *result) RowsAffected() (int64, error) { return r.affected, nil }

type resultRows struct {
	cols []string
	data [][]driver.Value
	pos  int
}

func (r *resultRows) Columns() []string { return r.cols }
func (r *resultRows) Close() error      { return nil }
func (r *resultRows) Next(dest []driver.Value) error {
	if r.pos >= len(r.data) {
		return io.EOF
	}
	copy(dest, r.data[r.pos])
	r.pos++
	return nil
}

func doExec(tables map[string]*tableState, q string, vals []driver.Value) (driver.Result, error) {
	normalized := strings.ReplaceAll(q, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\t", " ")

	upper := strings.ToUpper(strings.TrimSpace(normalized))

	if strings.HasPrefix(upper, "INSERT") && (strings.Contains(upper, "RETURNING") || strings.Contains(upper, "ON CONFLICT")) {
		return doInsertReturning(tables, normalized, vals)
	}

	switch {
	case strings.HasPrefix(upper, "INSERT"):
		return doInsert(tables, normalized, vals)
	case strings.HasPrefix(upper, "UPDATE"):
		return doUpdate(tables, normalized, vals)
	default:
		return nil, fmt.Errorf("unsupported statement: %s", q)
	}
}

func doInsert(tables map[string]*tableState, q string, vals []driver.Value) (driver.Result, error) {
	upper := strings.ToUpper(q)
	idxInto := strings.Index(upper, "INSERT INTO ")
	idxVals := strings.Index(upper, " VALUES ")
	idxOpen := strings.Index(q[idxVals:], "(")
	if idxInto < 0 || idxVals < 0 || idxOpen < 0 {
		return nil, fmt.Errorf("malformed INSERT: %s", q)
	}

	afterInsert := strings.TrimSpace(q[idxInto+len("INSERT INTO "):])
	spaceIdx := strings.Index(afterInsert, " ")
	parenIdx := strings.Index(afterInsert, "(")
	endIdx := spaceIdx
	if parenIdx >= 0 && (spaceIdx < 0 || parenIdx < spaceIdx) {
		endIdx = parenIdx
	}
	if endIdx < 0 {
		endIdx = len(afterInsert)
	}
	table := strings.TrimSpace(afterInsert[:endIdx])

	valsParenOpen := idxVals + idxOpen
	valsParenClose := findMatchingClose(q, valsParenOpen)
	if valsParenClose < 0 {
		return nil, fmt.Errorf("malformed VALUES parens: %s", q)
	}
	valsInner := q[valsParenOpen+1 : valsParenClose]

	colsParenClose := strings.LastIndex(q[:idxVals], ")")
	colsParenOpen := strings.LastIndex(q[:colsParenClose], "(")
	if colsParenOpen < 0 {
		return nil, fmt.Errorf("malformed columns parens: %s", q)
	}
	colsPart := q[colsParenOpen+1 : colsParenClose]

	cols := splitList(colsPart)
	rawVals := splitList(valsInner)
	if len(cols) != len(rawVals) {
		return nil, fmt.Errorf("column/value count mismatch: %d cols vs %d vals in %s", len(cols), len(rawVals), q)
	}

	r := make(row, len(cols))
	for i, c := range cols {
		r[c] = resolveVal(rawVals[i], vals)
	}

	tbl := tables[table]
	if tbl == nil {
		return nil, fmt.Errorf("unknown table %q", table)
	}
	tbl.rows = append(tbl.rows, r)
	return &result{affected: 1}, nil
}

func doInsertReturning(tables map[string]*tableState, q string, vals []driver.Value) (driver.Result, error) {
	upper := strings.ToUpper(q)

	idxInto := strings.Index(upper, "INSERT INTO ")
	if idxInto < 0 {
		return nil, fmt.Errorf("malformed INSERT: %s", q)
	}

	afterInsert := strings.TrimSpace(q[idxInto+len("INSERT INTO "):])

	// Find table name
	tableEndIdx := strings.Index(afterInsert, " ")
	parenIdx := strings.Index(afterInsert, "(")
	if parenIdx >= 0 && (tableEndIdx < 0 || parenIdx < tableEndIdx) {
		tableEndIdx = parenIdx
	}
	if tableEndIdx < 0 {
		tableEndIdx = len(afterInsert)
	}
	table := strings.TrimSpace(afterInsert[:tableEndIdx])

	// Find columns and values
	colonIdx := parenIdx
	if colonIdx < 0 {
		colonIdx = strings.Index(afterInsert, "(")
	}
	if colonIdx < 0 {
		return nil, fmt.Errorf("malformed INSERT columns: %s", q)
	}
	colsEndIdx := strings.Index(afterInsert[colonIdx:], ")")
	if colsEndIdx < 0 {
		return nil, fmt.Errorf("malformed INSERT columns: %s", q)
	}
	colsPart := afterInsert[colonIdx+1 : colonIdx+colsEndIdx]
	cols := splitList(colsPart)

	// Find VALUES
	idxVals := strings.Index(upper, " VALUES ")
	if idxVals < 0 {
		return nil, fmt.Errorf("malformed VALUES: %s", q)
	}
	valsParenOpen := strings.Index(q[idxVals:], "(") + idxVals
	valsParenClose := findMatchingClose(q, valsParenOpen)
	if valsParenClose < 0 {
		return nil, fmt.Errorf("malformed VALUES: %s", q)
	}
	valsInner := q[valsParenOpen+1 : valsParenClose]
	rawVals := splitList(valsInner)
	if len(cols) != len(rawVals) {
		return nil, fmt.Errorf("column/value count mismatch")
	}

	r := make(row, len(cols))
	for i, c := range cols {
		r[c] = resolveVal(rawVals[i], vals)
	}

	// Find conflict columns
	onConflictIdx := strings.Index(upper, "ON CONFLICT")
	doNothingIdx := strings.Index(upper, "DO NOTHING")
	if onConflictIdx < 0 || doNothingIdx < 0 || doNothingIdx < onConflictIdx {
		return nil, fmt.Errorf("malformed ON CONFLICT clause: %s", q)
	}

	conflictPart := q[onConflictIdx+len("ON CONFLICT "):doNothingIdx]
	conflictColStr := strings.TrimSpace(conflictPart)
	if strings.HasPrefix(conflictColStr, "(") {
		conflictColStr = conflictColStr[1:]
	}
	if strings.HasSuffix(conflictColStr, ")") {
		conflictColStr = conflictColStr[:len(conflictColStr)-1]
	}
	conflictCols := splitList(conflictColStr)

	tbl := tables[table]
	if tbl == nil {
		return nil, fmt.Errorf("unknown table %q", table)
	}

	// Check for conflict: compare against existing rows in the table
	for _, existingRow := range tbl.rows {
		allMatch := true
		for _, cc := range conflictCols {
			if existingRow[cc] != r[cc] {
				allMatch = false
				break
			}
		}
		if allMatch {
			// Conflict detected - DO NOTHING (return affected=0)
			return &result{affected: 0}, nil
		}
	}

	tbl.rows = append(tbl.rows, r)
	return &result{affected: 1}, nil
}

func findMatchingClose(s string, startIdx int) int {
	depth := 0
	for i := startIdx; i < len(s); i++ {
		if s[i] == '(' {
			depth++
		} else if s[i] == ')' {
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func doUpdate(tables map[string]*tableState, q string, vals []driver.Value) (driver.Result, error) {
	setIdx := strings.Index(strings.ToUpper(q), " SET ")
	whereIdx := strings.Index(strings.ToUpper(q), " WHERE ")
	if setIdx < 0 || whereIdx < 0 {
		return nil, fmt.Errorf("malformed UPDATE: %s", q)
	}
	table := strings.Fields(q[len("UPDATE "):setIdx])[0]

	setPart := q[setIdx+len(" SET ") : whereIdx]
	wherePart := q[whereIdx+len(" WHERE "):]

	// Support multiple column updates: col1=val1, col2=val2
	setAssignments := strings.Split(setPart, ",")
	type assignment struct {
		col string
		val any
	}
	var updates []assignment

	for _, p := range setAssignments {
		eq := strings.Index(p, "=")
		if eq < 0 {
			continue
		}
		setCol := strings.TrimSpace(p[:eq])
		setVal := resolveVal(strings.TrimSpace(p[eq+1:]), vals)
		updates = append(updates, assignment{col: setCol, val: setVal})
	}

	weq := strings.Index(wherePart, "=")
	if weq < 0 {
		return nil, fmt.Errorf("malformed WHERE clause: %s", wherePart)
	}
	wCol := strings.TrimSpace(wherePart[:weq])
	wVal := resolveVal(strings.TrimSpace(wherePart[weq+1:]), vals)

	tbl := tables[table]
	if tbl == nil {
		return nil, fmt.Errorf("unknown table %q", table)
	}
	var affected int64
	for _, r := range tbl.rows {
		if equalVal(r[wCol], wVal) {
			for _, u := range updates {
				r[u.col] = u.val
			}
			affected++
		}
	}
	return &result{affected: affected}, nil
}

func doQuery(tables map[string]*tableState, q string, vals []driver.Value) (driver.Rows, error) {
	upper := strings.ToUpper(q)

	if strings.Contains(upper, "COUNT(*)") {
		fromIdx := strings.Index(upper, "FROM")
		rest := q[fromIdx+len("FROM"):]
		whereIdx := strings.Index(strings.ToUpper(rest), "WHERE")

		var table string
		var conditions []struct {
			col string
			val any
		}
		if whereIdx < 0 {
			table = strings.TrimSpace(rest)
		} else {
			table = strings.TrimSpace(rest[:whereIdx])
			wherePart := rest[whereIdx+len("WHERE"):]
			conditions = parseWhereConditions(wherePart, vals)
		}

		tbl := tables[table]
		n := int64(0)
		if tbl != nil {
			for _, r := range tbl.rows {
				if len(conditions) == 0 || rowMatchesConditions(r, conditions) {
					n++
				}
			}
		}
		return &resultRows{cols: []string{"count"}, data: [][]driver.Value{{n}}}, nil
	}

	// SELECT <cols> FROM <table> WHERE <col> = <val> [LIMIT X] [FOR UPDATE SKIP LOCKED]
	sel := q[len("SELECT "):]
	fromIdx := strings.Index(strings.ToUpper(sel), "FROM")
	colsPart := strings.TrimSpace(sel[:fromIdx])
	rest := sel[fromIdx+len("FROM"):]
	whereIdx := strings.Index(strings.ToUpper(rest), "WHERE")
	if whereIdx < 0 {
		return nil, fmt.Errorf("unsupported SELECT: %s", q)
	}
	table := strings.TrimSpace(rest[:whereIdx])

	wherePart := rest[whereIdx+len("WHERE"):]
	conditions := parseWhereConditions(wherePart, vals)
	limitIdx := strings.Index(strings.ToUpper(wherePart), "LIMIT")
	forUpdateIdx := strings.Index(strings.ToUpper(wherePart), "FOR UPDATE")

	endWhereIdx := len(wherePart)
	var limit int64 = -1

	if limitIdx >= 0 {
		endWhereIdx = limitIdx
		limitPart := strings.TrimSpace(wherePart[limitIdx+len("LIMIT"):])
		if forUpdateIdx >= 0 && forUpdateIdx > limitIdx {
			limitPart = strings.TrimSpace(wherePart[limitIdx+len("LIMIT") : forUpdateIdx])
		} else {
			spaceIdx := strings.Index(limitPart, " ")
			if spaceIdx >= 0 {
				limitPart = limitPart[:spaceIdx]
			}
		}
		if l, err := strconv.ParseInt(limitPart, 10, 64); err == nil {
			limit = l
		}
	} else if forUpdateIdx >= 0 {
		endWhereIdx = forUpdateIdx
	}

	wherePart = wherePart[:endWhereIdx]

	cols := splitList(colsPart)
	tbl := tables[table]
	var data [][]driver.Value
	if tbl != nil {
		for _, r := range tbl.rows {
			if limit >= 0 && int64(len(data)) >= limit {
				break
			}
			if rowMatchesConditions(r, conditions) {
				rowVals := make([]driver.Value, len(cols))
				for i, c := range cols {
					rowVals[i] = r[c]
				}
				data = append(data, rowVals)
			}
		}
	}
	return &resultRows{cols: cols, data: data}, nil
}

func parseWhereConditions(wherePart string, vals []driver.Value) []struct {
	col string
	val any
} {
	parts := strings.Split(wherePart, " AND ")
	var conditions []struct{ col string; val any }

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || strings.HasPrefix(strings.ToUpper(part), "LIMIT") || strings.HasPrefix(strings.ToUpper(part), "FOR UPDATE") {
			continue
		}
		eq := strings.Index(part, "=")
		if eq < 0 {
			continue
		}
		col := strings.TrimSpace(part[:eq])
		valStr := strings.TrimSpace(part[eq+1:])
		conditions = append(conditions, struct{ col string; val any }{col: col, val: resolveVal(valStr, vals)})
	}
	return conditions
}

func rowMatchesConditions(r row, conditions []struct{ col string; val any }) bool {
	for _, c := range conditions {
		if !equalVal(r[c.col], c.val) {
			return false
		}
	}
	return true
}

func splitList(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}

func resolveVal(raw string, vals []driver.Value) any {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "$") {
		n, err := strconv.Atoi(raw[1:])
		if err == nil && n-1 < len(vals) {
			return vals[n-1]
		}
		return nil
	}
	if strings.ToUpper(raw) == "NULL" {
		return nil
	}
	if len(raw) >= 2 && strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'") {
		return raw[1 : len(raw)-1]
	}
	if i, err := strconv.ParseInt(raw, 10, 64); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(raw, 64); err == nil {
		return f
	}
	return raw
}

func equalVal(a, b any) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}