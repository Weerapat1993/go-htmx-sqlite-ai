// Command entrypoint applies pending DB migrations, then execs the server
// binary. It exists because the server's deploy image (gcr.io/distroless/static)
// has no shell, so migrations can't be chained via a shell script before startup —
// this replaces that shell step with a static Go binary. It is only used as the
// container CMD; local dev still runs migrate.sh directly.
package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "./db.sqlite3"
	}

	// golang-migrate's sqlite:// DSN parses the segment after "//" as the URL
	// host, so a relative path like "./db.sqlite3" is misread and the driver
	// returns SQLITE_CANTOPEN (error 14). Resolve to an absolute path first,
	// mirroring the workaround already documented in migrate.sh.
	dbURL = strings.TrimPrefix(dbURL, "file:")
	absDBPath, err := filepath.Abs(dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint: resolve db path:", err)
		os.Exit(1)
	}

	m, err := migrate.New("file:///migrations", "sqlite://"+absDBPath)
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
