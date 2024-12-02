package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	DB, err := sql.Open("sqlite3", "api.db")

	if err != nil {
		panic(err)
	}
	err = DB.Ping()
	if err != nil {
		panic(err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	err = createTableSQL(DB)

	return DB, err
}

func createTableSQL(db *sql.DB) error {
	createTableSQLusers := `CREATE TABLE IF NOT EXISTS users (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	);`

	statementuser, err := db.Prepare(createTableSQLusers)

	if err != nil {
		return err
	}
	defer statementuser.Close()
	_, err = statementuser.Exec()
	if err != nil {
		return err
	}

	return nil
}
