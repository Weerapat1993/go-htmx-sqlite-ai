// Command entrypoint applies pending DB migrations, then execs the server
// binary. It exists because the server's deploy image (gcr.io/distroless/static)
// has no shell, so migrations can't be chained via a shell script before startup —
// this replaces that shell step with a static Go binary. It is only used as the
// container CMD; local dev still runs migrate.sh directly.
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "./db.sqlite3"
	}

	var (
		m   *migrate.Migrate
		err error
	)
	if strings.HasPrefix(dbURL, "libsql://") {
		m, err = newRemoteMigrator(dbURL, os.Getenv("TURSO_AUTH_TOKEN"))
	} else {
		m, err = newLocalMigrator(dbURL)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint: migrate init:", err)
		os.Exit(1)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintln(os.Stderr, "entrypoint: migrate up:", err)
		os.Exit(1)
	}

	// #nosec G204 -- path and args are hardcoded literals, not user input.
	if err := syscall.Exec("/my-app", []string{"/my-app"}, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint: exec server:", err)
		os.Exit(1)
	}
}

func newLocalMigrator(dbURL string) (*migrate.Migrate, error) {
	// golang-migrate's sqlite:// DSN parses the segment after "//" as the URL
	// host, so a relative path like "./db.sqlite3" is misread and the driver
	// returns SQLITE_CANTOPEN (error 14). Resolve to an absolute path first,
	// mirroring the workaround already documented in migrate.sh.
	path := strings.TrimPrefix(dbURL, "file:")
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}
	m, err := migrate.New("file:///migrations", "sqlite://"+absPath)
	if err != nil {
		return nil, fmt.Errorf("new migrator: %w", err)
	}
	return m, nil
}

func newRemoteMigrator(dbURL, authToken string) (*migrate.Migrate, error) {
	dsn := dbURL
	if authToken != "" {
		sep := "?"
		if strings.Contains(dbURL, "?") {
			sep = "&"
		}
		dsn += sep + "authToken=" + authToken
	}
	db, err := sql.Open("libsql", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	// libsql speaks the SQLite wire protocol/dialect, so golang-migrate's
	// sqlite driver works against it directly via an existing *sql.DB —
	// there is no dedicated "libsql" scheme registered with golang-migrate.
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("libsql migrate driver: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file:///migrations", "libsql", driver)
	if err != nil {
		return nil, fmt.Errorf("new migrator with database instance: %w", err)
	}
	return m, nil
}
