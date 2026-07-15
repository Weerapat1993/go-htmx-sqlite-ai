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

	m, err := migrate.New("file:///migrations", "sqlite://"+dbURL)
	if err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint: migrate init:", err)
		os.Exit(1)
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		fmt.Fprintln(os.Stderr, "entrypoint: migrate up:", err)
		os.Exit(1)
	}

	if err := syscall.Exec("/my-app", []string{"/my-app"}, os.Environ()); err != nil {
		fmt.Fprintln(os.Stderr, "entrypoint: exec server:", err)
		os.Exit(1)
	}
}
