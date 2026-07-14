package router

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/dist"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/server/handler"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/server/middleware"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/version"
)

// New creates a new router with the given context, logger, database, and rate limit.
func New(ctx context.Context, logger *slog.Logger, database db.Database, rateLimit int) http.Handler {
	h := handler.New(logger, database)

	ipCfg := middleware.IPConfig{
		TrustProxyHeaders: version.Value != "dev",
	}

	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc(newPath(http.MethodGet, "/health"), h.Health)
	mux.Handle(newPath(http.MethodGet, "/assets/"), middleware.CacheMiddleware(http.FileServer(http.FS(dist.AssetsDir))))
	mux.HandleFunc(newPath(http.MethodGet, "/{$}"), h.Home)
	mux.HandleFunc(newPath(http.MethodGet, "/about"), h.About)
	mux.HandleFunc(newPath(http.MethodPost, "/count"), h.Count)

	// Middleware chain
	hdlr := http.Handler(mux)
	hdlr = middleware.Chain(
		middleware.Recovery(logger, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		})),
		middleware.Logging(logger, ipCfg),
		middleware.Security(logger, ipCfg),
		middleware.RateLimit(ctx, logger, rateLimit, middleware.DefaultMaxEntries, ipCfg),
		middleware.CSRF(logger, ipCfg),
	)(hdlr)

	return hdlr
}

func newPath(method string, path string) string {
	return method + " " + path
}
