package router

import (
	"net/http"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/server/handler"
)

// newAPIMux builds the "/api/*" route table.
func newAPIMux(h *handler.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc(newPath(http.MethodGet, "/api/health"), h.Health)

	return mux
}
