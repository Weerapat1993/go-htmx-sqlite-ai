package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db/queries"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func newRemoteDB(url, authToken string) (*sqlDB, error) {
	dsn := url
	if authToken != "" {
		sep := "?"
		if strings.Contains(url, "?") {
			sep = "&"
		}
		dsn += sep + "authToken=" + authToken
	}
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	db.SetMaxOpenConns(maxOpenConns)
	return &sqlDB{db: db, queries: queries.New(db)}, nil
}
