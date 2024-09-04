package implementations

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	ConnectionString string
	Db               *sql.DB
}

func (s *Storage) Open(connection string) bool {
	db, err := sql.Open("postgres", s.ConnectionString)
	if err != nil {

		fmt.Println(err)
		return false
	}

	s.Db = db
	return true
}

func (s *Storage) Single(sql string, params []interface{}) *sql.Row {
	row := s.Db.QueryRow(sql, params...)

	return row
}

func (s *Storage) Exec(sql string, params []interface{}) (*sql.Result, error) {
	row, err := s.Db.Exec(sql, params...)

	return &row, err
}

func (s *Storage) Where(sql string, params []interface{}) *sql.Rows {

	var err error

	rows, err := s.Db.Query("SELECT * from public.accounts")

	if err != nil {
		fmt.Println("Failed to get accounts", err)
	}
	return rows
}

func (s *Storage) Close() bool {
	return true
}
