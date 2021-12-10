package ploto

import (
	"context"
	"database/sql"
	"github.com/feiin/sqlstring"
)

type DB struct {
	*sql.DB
	LogSql bool
	logger LoggerInterface
}

type RowsResult struct {
	*sql.Rows
	LastError error
}

type RowResult struct {
	rows      *sql.Rows
	LastError error
}

//RawDB return the *sql.DB
func (db *DB) RawDB() *sql.DB {
	return db.DB
}

// Query executes a query that returns RowsResult, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) Query(query string, args ...interface{}) *RowsResult {
	if db.LogSql {
		db.logger.Info("query sql:%s", sqlstring.Format(query, args...))
	}
	rs, err := db.DB.Query(query, args...)
	return &RowsResult{rs, err}
}

// QueryContext executes a query that returns RowsResult, typically a SELECT.
// The args are for any placeholder parameters in the query.
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) *RowsResult {
	if db.LogSql {
		db.logger.Info("QueryContext sql:%s", sqlstring.Format(query, args...))
	}
	rs, err := db.DB.QueryContext(ctx, query, args...)
	return &RowsResult{rs, err}
}

// QueryRowContext executes a query that is expected to return at most one row.
// QueryRowContext always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *RowResult {
	if db.LogSql {
		db.logger.Info("QueryRowContext sql:%s", sqlstring.Format(query, args...))
	}
	rows, err := db.DB.QueryContext(ctx, query, args...)
	return &RowResult{rows: rows, LastError: err}
}

// QueryRow executes a query that is expected to return at most one row.
// QueryRow always returns a non-nil value. Errors are deferred until
// Row's Scan method is called.
// If the query selects no rows, the *Row's Scan will return ErrNoRows.
// Otherwise, the *Row's Scan scans the first selected row and discards
// the rest.
func (db *DB) QueryRow(query string, args ...interface{}) *RowResult {
	if db.LogSql {
		db.logger.Info("QueryRow sql:%s", sqlstring.Format(query, args...))
	}
	return db.QueryRowContext(context.Background(), query, args...)
}

//Close returns the connection to the connection pool
func (r RowsResult) Close() error {
	return r.Rows.Close()
}

//Scan
func (r *RowsResult) Scan(dest interface{}) error {

	if r.Err() != nil {
		return r.Err()
	}
	err := ScanResult(r.Rows, dest)
	return err
}

//Raw
func (r *RowsResult) Raw() (*sql.Rows, error) {
	return r.Rows, r.LastError
}

//RowResult return the error of RowResult
func (r *RowResult) Err() error {
	return r.LastError
}

//Scan RowResult's scan
func (r *RowResult) Scan(dest interface{}) error {

	if r.Err() != nil {
		return r.Err()
	}

	if r.rows.Err() != nil {
		return r.rows.Err()
	}
	defer r.rows.Close()

	if !r.rows.Next() {
		if err := r.rows.Err(); err != nil {
			return err
		}
		return sql.ErrNoRows
	}

	err := Scan(r.rows, dest)

	if err != nil {
		return err
	}

	// Make sure the query can be processed to completion with no errors.
	return r.rows.Close()
}

//Init init all the database clients
func Open(configs DialectConfig, log LoggerInterface) (*Dialect, error) {
	dialect := &Dialect{}
	dialect.Clients = make(map[string]*DB)
	dialect.Configs = configs

	if log == nil {
		dialect.logger = DefaultLogger{}
	} else {
		dialect.logger = log
	}

	for k := range configs.Clients {

		db, err := dialect.CreateClient(k)
		db.logger = dialect.logger
		if err != nil {
			return nil, err
		}
		dialect.Clients[k] = db
	}

	return dialect, nil
}
