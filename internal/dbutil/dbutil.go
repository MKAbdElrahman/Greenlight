package dbutil

import (
	"context"
	"database/sql"
	"time"
)

type DBConfig struct {
	DSN          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

func NewDBConnection(cfg DBConfig) (*sql.DB, error) {
	db, err := createDBConnection(cfg.DSN)
	if err != nil {
		return nil, err
	}

	if err := configureDBConnectionPool(db, cfg); err != nil {
		return nil, err
	}

	err = testDBConnection(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDBConnection(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	return db, err
}

func testDBConnection(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := db.PingContext(ctx)
	return err
}

func configureDBConnectionPool(db *sql.DB, cfg DBConfig) error {
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)

	duration, err := time.ParseDuration(cfg.MaxIdleTime)
	if err != nil {
		return err
	}

	db.SetConnMaxIdleTime(duration)
	return nil
}
