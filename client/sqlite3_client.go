package client

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type SQLite3Config struct {
	DataSourceName string

	MaxIdleConn     int
	MaxOpenConn     int
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

func NewSQLite3(cfg *SQLite3Config) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", cfg.DataSourceName)
	if err != nil {
		return nil, fmt.Errorf("connect sqlite3: %w", err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxLifetime(cfg.MaxConnLifetime)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping sqlite3: %w", err)
	}

	return db, nil
}
