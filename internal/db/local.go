package db

import (
	"database/sql"
	"fmt"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db/queries"

	_ "modernc.org/sqlite"
)

func newLocalDB(path string) (*sqlDB, error) {
	db, err := sql.Open("sqlite", "file:"+path+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)")
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	return &sqlDB{db: db, queries: queries.New(db)}, nil
}
