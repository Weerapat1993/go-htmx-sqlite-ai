package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db/queries"
)

const maxOpenConns = 4

type Database interface {
	DB() *sql.DB
	Queries() *queries.Queries
	Close() error
}

type sqlDB struct {
	db      *sql.DB
	queries *queries.Queries
}

var _ Database = (*sqlDB)(nil)

func (d *sqlDB) DB() *sql.DB {
	return d.db
}

func (d *sqlDB) Queries() *queries.Queries {
	return d.queries
}

func (d *sqlDB) Close() error {
	if err := d.db.Close(); err != nil {
		return fmt.Errorf("closing database: %w", err)
	}
	return nil
}

// New opens a database connection. A "libsql://" URL connects to a remote
// Turso database (auth token read from TURSO_AUTH_TOKEN); any other URL is
// treated as a local SQLite file path.
func New(url string) (Database, error) {
	var (
		db  *sqlDB
		err error
	)
	if strings.HasPrefix(url, "libsql://") {
		db, err = newRemoteDB(url, os.Getenv("TURSO_AUTH_TOKEN"))
	} else {
		db, err = newLocalDB(url)
	}
	if err != nil {
		return nil, err
	}
	if err = db.DB().PingContext(context.Background()); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}
	return db, nil
}

// NewFromRawDB creates a Database from an existing *sql.DB. Useful for testing.
func NewFromRawDB(rawDB *sql.DB) Database {
	rawDB.SetMaxOpenConns(maxOpenConns)
	return &sqlDB{db: rawDB, queries: queries.New(rawDB)}
}
