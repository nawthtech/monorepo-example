package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nawthtech/nawthtech/backend/internal/config"
)

var SQLDB *sql.DB

// InitializeSQLDB initializes local SQL (sqlite) for dev or other SQL drivers if configured.
func InitializeSQLDB(cfg *config.Config) error {
	// If D1 is requested, return nil and let user implement D1 adapter.
	if cfg.D1.UseD1 {
		// The app will still run; any D1 calls should go through the D1 adapter which is not auto-implemented here.
		return nil
	}

	// local sqlite file
	dbPath := os.Getenv("SQLITE_PATH")
	if dbPath == "" {
		dbPath = "./data/nawthtech.db"
	}

	// ensure dir exists
	// (omitted: create dir if needed)

	dsn := fmt.Sprintf("%s?_foreign_keys=1", dbPath)
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	// quick ping
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	SQLDB = db
	return nil
}

func Close() error {
	if SQLDB != nil {
		return SQLDB.Close()
	}
	return nil
}

// Exec executes statement
func Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if SQLDB == nil {
		return nil, fmt.Errorf("sql db not initialized")
	}
	return SQLDB.ExecContext(ctx, query, args...)
}

// QueryRow returns single row
func QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	if SQLDB == nil {
		return nil
	}
	return SQLDB.QueryRowContext(ctx, query, args...)
}

// Query returns rows
func Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if SQLDB == nil {
		return nil, fmt.Errorf("sql db not initialized")
	}
	return SQLDB.QueryContext(ctx, query, args...)
}