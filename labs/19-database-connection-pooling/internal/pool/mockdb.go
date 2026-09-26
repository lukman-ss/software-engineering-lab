package pool

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ErrServerOverloaded = errors.New("server max connections exceeded")
)

type MockDriver struct {
	mu             sync.Mutex
	maxConnections int32
	activeConns    int32
	totalCreated   int32
	connectDelay   time.Duration
}

func NewMockDriver(maxConns int32, connectDelay time.Duration) *MockDriver {
	return &MockDriver{
		maxConnections: maxConns,
		connectDelay:   connectDelay,
	}
}

func (d *MockDriver) Open(name string) (driver.Conn, error) {
	if d.connectDelay > 0 {
		time.Sleep(d.connectDelay)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	if d.maxConnections > 0 && d.activeConns >= d.maxConnections {
		return nil, ErrServerOverloaded
	}

	atomic.AddInt32(&d.activeConns, 1)
	atomic.AddInt32(&d.totalCreated, 1)

	return &mockConn{driver: d}, nil
}

func (d *MockDriver) ActiveConnections() int32 {
	return atomic.LoadInt32(&d.activeConns)
}

func (d *MockDriver) TotalCreated() int32 {
	return atomic.LoadInt32(&d.totalCreated)
}

type mockConn struct {
	driver *MockDriver
	closed bool
	mu     sync.Mutex
}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	return &mockStmt{conn: c}, nil
}

func (c *mockConn) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !c.closed {
		c.closed = true
		atomic.AddInt32(&c.driver.activeConns, -1)
	}
	return nil
}

func (c *mockConn) Begin() (driver.Tx, error) {
	return &mockTx{conn: c}, nil
}

type mockStmt struct {
	conn *mockConn
}

func (s *mockStmt) Close() error {
	return nil
}

func (s *mockStmt) NumInput() int {
	return -1
}

func (s *mockStmt) Exec(args []driver.Value) (driver.Result, error) {
	return driver.RowsAffected(1), nil
}

func (s *mockStmt) Query(args []driver.Value) (driver.Rows, error) {
	return &mockRows{}, nil
}

type mockTx struct {
	conn *mockConn
}

func (t *mockTx) Commit() error {
	return nil
}

func (t *mockTx) Rollback() error {
	return nil
}

type mockRows struct {
	done bool
}

func (r *mockRows) Columns() []string {
	return []string{"id"}
}

func (r *mockRows) Close() error {
	return nil
}

func (r *mockRows) Next(dest []driver.Value) error {
	if r.done {
		return io.EOF
	}
	r.done = true
	dest[0] = int64(1)
	return nil
}

// Connector implementation to avoid global driver registration
type MockConnector struct {
	driver *MockDriver
}

func NewMockConnector(d *MockDriver) *MockConnector {
	return &MockConnector{driver: d}
}

func (c *MockConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return c.driver.Open("")
}

func (c *MockConnector) Driver() driver.Driver {
	return c.driver
}

func OpenDB(d *MockDriver) *sql.DB {
	return sql.OpenDB(NewMockConnector(d))
}
