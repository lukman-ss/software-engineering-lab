// Package mockdb provides a minimal in-memory database/sql driver so the
// transaction labs can run deterministically without a real database,
// credentials, or network access. It supports the small subset of SQL used by
// the lab examples: INSERT, UPDATE, and SELECT (COUNT(*) and single-row reads).
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

// --- in-memory state ---

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

// --- driver plumbing ---

type staticDriver struct{}

func (staticDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("use sql.OpenDB(mockdb.Connector())")
}

type connector struct{ db *DB }

// Connector returns a driver.Connector backed by a fresh, isolated in-memory
// database. Each call yields independent state, so concurrent tests never share
// rows.
func Connector() driver.Connector { return &connector{db: newDB()} }

// NewDB returns a *sql.DB wired to a private in-memory database.
func NewDB() *sql.DB { return sql.OpenDB(Connector()) }

func (c *connector) Connect(context.Context) (driver.Conn, error) { return &conn{db: c.db}, nil }
func (c *connector) Driver() driver.Driver                        { return staticDriver{} }

type conn struct {
	db *DB
	tx map[string]*tableState // non-nil while inside a transaction
}

func (c *conn) Prepare(query string) (driver.Stmt, error) { return &stmt{c: c, query: query}, nil }
func (c *conn) PrepareContext(_ context.Context, query string) (driver.Stmt, error) {
	return &stmt{c: c, query: query}, nil
}
func (c *conn) Close() error               { return nil }
func (c *conn) Ping(context.Context) error { return nil }

func (c *conn) Begin() (driver.Tx, error) { return c.begin() }
func (c *conn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	return c.begin()
}

// begin clones the committed state. All subsequent statements on this
// connection operate on the clone until Commit or Rollback.
func (c *conn) begin() (driver.Tx, error) {
	if c.tx != nil {
		return nil, errors.New("transaction already open on connection")
	}
	c.db.mu.Lock()
	c.tx = cloneTables(c.db.committed)
	c.db.mu.Unlock()
	return &txn{c: c}, nil
}

type txn struct{ c *conn }

func (t *txn) Commit() error {
	c := t.c
	c.db.mu.Lock()
	c.db.committed = c.tx
	c.db.mu.Unlock()
	c.tx = nil
	return nil
}

func (t *txn) Rollback() error {
	t.c.tx = nil // discard the clone; committed state is untouched
	return nil
}

func (c *conn) ExecContext(_ context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	return c.exec(query, toVals(args))
}

func (c *conn) QueryContext(_ context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	return c.query(query, toVals(args))
}

func (c *conn) exec(query string, vals []driver.Value) (driver.Result, error) {
	if c.tx != nil {
		return doExec(c.tx, query, vals)
	}
	c.db.mu.Lock()
	defer c.db.mu.Unlock()
	return doExec(c.db.committed, query, vals)
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

// --- statement execution ---

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
	// Normalize whitespace: replace newlines with spaces for easier parsing
	normalized := strings.ReplaceAll(q, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\t", " ")

	switch {
	case strings.HasPrefix(strings.ToUpper(strings.TrimSpace(normalized)), "INSERT"):
		return doInsert(tables, normalized, vals)
	case strings.HasPrefix(strings.ToUpper(strings.TrimSpace(normalized)), "UPDATE"):
		return doUpdate(tables, normalized, vals)
	default:
		return nil, fmt.Errorf("unsupported statement: %s", q)
	}
}

func doInsert(tables map[string]*tableState, q string, vals []driver.Value) (driver.Result, error) {
	// Parse: INSERT INTO <table> (col1, col2) VALUES ($1, $2)
	upper := strings.ToUpper(q)
	idxInto := strings.Index(upper, "INSERT INTO ")
	idxVals := strings.Index(upper, " VALUES ")
	idxOpen := strings.Index(q[idxVals:], "(")
	if idxInto < 0 || idxVals < 0 || idxOpen < 0 {
		return nil, fmt.Errorf("malformed INSERT: %s", q)
	}

	// Table name is the first word after "INSERT INTO "
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

	// Find opening paren for VALUES
	valsParenOpen := idxVals + idxOpen
	valsParenClose := findMatchingClose(q, valsParenOpen)
	if valsParenClose < 0 {
		return nil, fmt.Errorf("malformed VALUES parens: %s", q)
	}
	valsInner := q[valsParenOpen+1 : valsParenClose]

	// Find the columns paren before VALUES
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

	t := tables[table]
	if t == nil {
		return nil, fmt.Errorf("unknown table %q", table)
	}
	t.rows = append(t.rows, r)
	return &result{affected: 1}, nil
}

// findMatchingClose finds the closing ')' matching the opening at startIdx.
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

	eq := strings.Index(setPart, "=")
	setCol := strings.TrimSpace(setPart[:eq])
	setVal := resolveVal(strings.TrimSpace(setPart[eq+1:]), vals)

	weq := strings.Index(wherePart, "=")
	wCol := strings.TrimSpace(wherePart[:weq])
	wVal := resolveVal(strings.TrimSpace(wherePart[weq+1:]), vals)

	t := tables[table]
	if t == nil {
		return nil, fmt.Errorf("unknown table %q", table)
	}
	var affected int64
	for _, r := range t.rows {
		if equalVal(r[wCol], wVal) {
			r[setCol] = setVal
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
		var wCol string
		var wVal any
		if whereIdx < 0 {
			table = strings.TrimSpace(rest)
		} else {
			table = strings.TrimSpace(rest[:whereIdx])
			wherePart := rest[whereIdx+len("WHERE"):]
			weq := strings.Index(wherePart, "=")
			wCol = strings.TrimSpace(wherePart[:weq])
			wVal = resolveVal(strings.TrimSpace(wherePart[weq+1:]), vals)
		}

		t := tables[table]
		n := int64(0)
		if t != nil {
			for _, r := range t.rows {
				if wCol == "" || equalVal(r[wCol], wVal) {
					n++
				}
			}
		}
		return &resultRows{cols: []string{"count"}, data: [][]driver.Value{{n}}}, nil
	}

	// SELECT <cols> FROM <table> WHERE <col> = <val>
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
	weq := strings.Index(wherePart, "=")
	wCol := strings.TrimSpace(wherePart[:weq])
	wVal := resolveVal(strings.TrimSpace(wherePart[weq+1:]), vals)

	cols := splitList(colsPart)
	t := tables[table]
	var data [][]driver.Value
	if t != nil {
		for _, r := range t.rows {
			if equalVal(r[wCol], wVal) {
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

// --- helpers ---

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
