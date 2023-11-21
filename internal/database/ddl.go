package database

import (
	"database/sql"
	"fmt"
)

var ddlCreateUserTable = `CREATE TABLE IF NOT EXISTS user (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	nickname VARCHAR(20) NOT NULL,
	username VARCHAR(50) NOT NULL,
	password VARCHAR(250) NOT NULL,
    created_at TIMESTAMP(6) NOT NULL,
	updated_at TIMESTAMP(6) NOT NULL,
	deleted_at TIMESTAMP(6) NULL DEFAULT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS ux_username on user (username);`

func CreateTables(db *sql.DB) error {
	if _, err := db.Exec(ddlCreateUserTable); err != nil {
		return fmt.Errorf("create user table: %w", err)
	}

	return nil
}
