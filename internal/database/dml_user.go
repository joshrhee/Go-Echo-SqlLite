package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"time"
)

var ErrDuplicateUsername = errors.New("duplicate username")

const (
	dmlUserInsertStmt = "INSERT INTO `user` " + `(
			nickname, username, password, created_at, updated_at
		) VALUES (
			?, ?, ?, ?, ?
		)
	`
)

func CreateUser(ctx context.Context, db *sql.DB, nickname string, username string, password string, createTime time.Time) (int64, error) {
	result, err := db.ExecContext(
		ctx,
		dmlUserInsertStmt,
		nickname,
		username,
		password,
		createTime,
		createTime,
	)
	if err != nil {
		var tErr sqlite3.Error
		if errors.As(err, &tErr) {
			if tErr.Code == sqlite3.ErrConstraint && tErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, fmt.Errorf("insert user: %w", ErrDuplicateUsername)
			}
		}
		return 0, fmt.Errorf("insert user: %w", err)
	}

	return result.LastInsertId()
}
