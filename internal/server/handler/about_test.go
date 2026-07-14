package handler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbout_StatusCode(t *testing.T) {
	t.Parallel()

	h := newHomeHandler()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()

	h.About(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestAbout_ContentType(t *testing.T) {
	t.Parallel()

	h := newHomeHandler()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()

	h.About(rec, req)

	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
}

func TestAbout_ContainsHeadline(t *testing.T) {
	t.Parallel()

	h := newHomeHandler()
	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/about", nil)
	rec := httptest.NewRecorder()

	h.About(rec, req)

	assert.Contains(t, rec.Body.String(), "Why this stack")
}
