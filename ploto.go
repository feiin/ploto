package ploto

import (
	"context"
	"database/sql"
)

type DB struct {
	*sql.DB
}

type RowsResult struct {
	*sql.Rows
	LastError error
}

type RowResult struct {
	*sql.Row
	LastError error
}

//RawDB return the *sql.DB
func (db *DB) RawDB() *sql.DB {
	return db.DB
}

//Query
func (db *DB) Query(query string, args ...interface{}) *RowsResult {
	rs, err := db.DB.Query(query, args...)
	return &RowsResult{rs, err}
}

//QueryContext
func (db *DB) QueryContext(ctx context.Context, query string, args ...interface{}) *RowsResult {
	rs, err := db.DB.QueryContext(ctx, query, args...)
	return &RowsResult{rs, err}
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
	defer r.Rows.Close()
	err := ScanResult(r.Rows, dest)
	return err
}

//Raw
func (r *RowsResult) Raw() (*sql.Rows, error) {
	return r.Rows, r.LastError
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
		if err != nil {
			return nil, err
		}
		dialect.Clients[k] = db
	}

	return dialect, nil
}
