package ploto

import (
	// "database/sql"
	"github.com/DATA-DOG/go-sqlmock"

	"testing"
	// "fmt"
	// "time"
)

func BenchmarkScanStruct(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
			AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id=?", 1)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		var user Users
		if rows.Next() {
			Scan(rows, &user)
		}
		if user.Id <= 0 {
			b.Fatalf("Scan error %+v", user)
		}

		rows.Close()

	}
}

func BenchmarkScanStructs(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
			AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id=?", 1)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		var users []Users
		ScanSlice(rows, &users)

		if len(users) != 1 {
			b.Fatalf("Scan error %+v", users)
		}

		rows.Close()

	}
}

func BenchmarkScanStructs2(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
			AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id=?", 1)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		var users []Users
		for rows.Next() {
			var user Users
			Scan(rows, &user)
			users = append(users, user)
		}

		if len(users) != 1 {
			b.Fatalf("Scan error %+v", users)
		}

		rows.Close()

	}
}

func BenchmarkScanOneField(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	type Users struct {
		Id int64 `db:"id"`
	}

	for i := 0; i < b.N; i++ {

		dataRows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id FROM users WHERE id=?", 1)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		var users Users
		if rows.Next() {
			Scan(rows, &users)
		}

		if users.Id != 1 {
			b.Fatalf("Scan error %+v", users)
		}

		rows.Close()

	}
}

func BenchmarkScanScalar(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		dataRows := sqlmock.NewRows([]string{"id"}).
			AddRow(1)
		mock.ExpectQuery("SELECT id FROM users WHERE id=?").WithArgs(1).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id FROM users WHERE id=?", 1)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		var id int64
		if rows.Next() {
			Scan(rows, &id)
		}

		if id != 1 {
			b.Fatalf("Scan error %+v", id)
		}

		rows.Close()

	}
}

func TestSqlMock(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
		AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(2, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(3, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(4, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
		AddRow(5, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
	mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(5).WillReturnRows(dataRows)

	rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id < ?", 5)
	if err != nil {
		t.Errorf("error '%s' was not expected while retrieving mock rows", err)
	}
	defer rows.Close()
	var users []Users

	err = ScanSlice(rows, &users)

	t.Logf("users %+v %+v", users, err)

}
func BenchmarkScanRows(b *testing.B) {

	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
			AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(2, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(3, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(4, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(5, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(100).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id<?", 100)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}

		b.StartTimer()
		var users []Users
		err = ScanSlice(rows, &users)

		rows.Close()

	}
}

func BenchmarkScanRows2(b *testing.B) {

	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		b.StopTimer()
		dataRows := sqlmock.NewRows([]string{"id", "name", "created_time", "updated_time"}).
			AddRow(1, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(2, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(3, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(4, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00").
			AddRow(5, "1111", "2021-10-01 00:00:00", "2021-10-01 00:00:00")
		mock.ExpectQuery("SELECT (.+) FROM users WHERE id<?").WithArgs(100).WillReturnRows(dataRows)

		rows, err := db.Query("SELECT id,name,created_time,updated_time FROM users WHERE id<?", 100)
		if err != nil {
			b.Fatalf("sql.Query: Error: %+v\n", err)
		}
		b.StartTimer()

		var users []Users
		for rows.Next() {
			var user Users
			Scan(rows, &user)
			users = append(users, user)
		}
		// err = ScanSlice(rows, &users)
		if len(users) != 5 {
			b.Fatalf("scan error %+v", users)
		}

		rows.Close()

	}
}
