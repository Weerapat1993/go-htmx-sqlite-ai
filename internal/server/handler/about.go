package handler

import (
	"net/http"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/about"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/core"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/layout"
)

// About handles the about page.
func (h *Handler) About(w http.ResponseWriter, r *http.Request) {
	page := layout.Shell("about", about.Page())
	h.html(r.Context(), w, http.StatusOK, core.HTML("About — go-htmx-sqlite-ai", page))
}
