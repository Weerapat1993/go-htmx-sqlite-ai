package handler

import (
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/core"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/home"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/layout"
	"net/http"
	"sync/atomic"
)

var counter atomic.Int64 //nolint:gochecknoglobals // demo counter for template

// Home handles the home page.
func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	page := layout.Shell("home", home.Page(int(counter.Load())))
	h.html(r.Context(), w, core.HTML("go-htmx-sqlite-ai — No JS build step. Just Go.", page))
}

// Count increments the counter and returns the updated Counter fragment.
func (h *Handler) Count(w http.ResponseWriter, r *http.Request) {
	newCount := counter.Add(1)
	h.html(r.Context(), w, home.Counter(int(newCount)))
}
