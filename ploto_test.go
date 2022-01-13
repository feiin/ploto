package ploto

import (
	"github.com/DATA-DOG/go-sqlmock"

	"database/sql"
	// "errors"
	"context"
	"testing"
)

func TestRowResultErr(t *testing.T) {
	rowResult := RowResult{
		rows:      nil,
		LastError: sql.ErrNoRows,
	}

	if rowResult.Err() != sql.ErrNoRows {
		t.Fatalf("should be ErrNoRows")
	}
}

func TestRowsResultErr(t *testing.T) {
	rowResult := RowsResult{

		LastError: sql.ErrNoRows,
	}

	if _, err := rowResult.Raw(); err != sql.ErrNoRows {
		t.Fatalf("should be ErrNoRows")
	}

}

func TestRowsResultRaw(t *testing.T) {
	rowsResult := RowsResult{
		LastError: sql.ErrNoRows,
	}

	rows, err := rowsResult.Raw()
	if err != sql.ErrNoRows {
		t.Fatalf("should be ErrNoRows")
	}

	if rows != nil {
		t.Fatalf("rows should be nil")
	}
}

func TestQueryEmptyRows(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"})
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var users []Users
	err = db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id=?", 1).Scan(&users)

	if err != nil {
		t.Fatalf("query with error %+v", err)
	}
	if len(users) > 0 {
		t.Fatalf("should be empty")
	}

}

func TestQuery5Rows(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
		AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(2, "2222", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(3, "3333", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(4, "4444", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(5, "5555", "2021-10-01 00:00:00", "2021-10-01 00:00:00")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(6).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var users []Users
	err = db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id<?", 6).Scan(&users)

	if err != nil {
		t.Fatalf("query with error %+v", err)
	}
	if len(users) != 5 {
		t.Fatalf("should be 5 users")
	}

}

func TestQueryEmptyRow(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"})
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var user Users
	err = db.QueryRow("SELECT id,name,created_time,updated_time FROM users WHERE id=?", 1).Scan(&user)

	if err != sql.ErrNoRows {
		t.Fatalf("should return the ErrNoRows: %+v", err)
	}

}

func TestQueryOneRow(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
		AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(2, "2222", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(3, "3333", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(4, "4444", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(5, "5555", "2021-10-01 00:00:00", "2021-10-01 00:00:00")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(6).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var user Users
	err = db.QueryRow("SELECT id,name,created_time,updated_time FROM users WHERE id<?", 6).Scan(&user)

	if err != nil {
		t.Fatalf("should return the one users: %+v", err)
	}

	if user.Id == 0 {
		t.Fatalf("should return the one users")

	}

}

func TestQueryContextRows(t *testing.T) {

	ctx := context.Background()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
		AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(2, "2222", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(3, "3333", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(4, "4444", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(5, "5555", "2021-10-01 00:00:00", "2021-10-01 00:00:00")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(6).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var users []Users
	err = db.QueryContext(ctx, "SELECT id,name,created_time,updated_time FROM users WHERE id<?", 6).Scan(&users)

	if err != nil {
		t.Fatalf("query with error %+v", err)
	}
	if len(users) != 5 {
		t.Fatalf("should be 5 users")
	}
	t.Logf("users:%v", users)

}

func TestExecContext(t *testing.T) {

	ctx := context.Background()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("update users").WillReturnResult(sqlmock.NewResult(1, 1))

	db := &DB{DB: mockDB}

	result, err := db.ExecContext(ctx, "update users set name='xxxx' WHERE id=1")

	if err != nil {
		t.Fatalf("query with error %+v", err)
	}

	affected, err := result.RowsAffected()

	t.Logf("update:%d", affected)

}

func TestExec(t *testing.T) {

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	mock.ExpectExec("update users").WillReturnResult(sqlmock.NewResult(1, 1))

	db := &DB{DB: mockDB}

	result, err := db.Exec("update users set name='xxxx' WHERE id=1")

	if err != nil {
		t.Fatalf("query with error %+v", err)
	}

	affected, err := result.RowsAffected()

	t.Logf("update:%d", affected)

}

func TestQueryContextOneRow(t *testing.T) {

	ctx := context.Background()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockDB.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
		AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(2, "2222", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(3, "3333", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(4, "4444", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(5, "5555", "2021-10-01 00:00:00", "2021-10-01 00:00:00")

	mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(6).WillReturnRows(dataRows)

	db := &DB{DB: mockDB}

	var user Users
	err = db.QueryContext(ctx, "SELECT id,name,created_time,updated_time FROM users WHERE id<?", 6).Scan(&user)

	if err != nil {
		t.Fatalf("should return the one users: %+v", err)
	}

	if user.Id == 0 {
		t.Fatalf("should return the one users")

	}

}
