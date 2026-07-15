package handler_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/server/handler"
)

func TestHealth(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		validate func(*testing.T, *httptest.ResponseRecorder)
	}{
		{
			name: "returns 200 status code",
			validate: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				assert.Equal(t, http.StatusOK, rec.Code)
			},
		},
		{
			name: "returns JSON content type",
			validate: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
			},
		},
		{
			name: "returns correct JSON body",
			validate: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				body, err := io.ReadAll(rec.Body)
				require.NoError(t, err)
				assert.JSONEq(t, `{"status":"ok","version":"dev"}`, string(body))
			},
		},
		{
			name: "returns valid JSON structure with status and version fields",
			validate: func(t *testing.T, rec *httptest.ResponseRecorder) {
				t.Helper()
				var result map[string]string
				err := json.NewDecoder(rec.Body).Decode(&result)
				require.NoError(t, err)
				assert.Equal(t, "ok", result["status"])
				assert.Equal(t, "dev", result["version"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := handler.New(slog.New(slog.DiscardHandler), nil)
			req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/health", nil)
			rec := httptest.NewRecorder()

			h.Health(rec, req)

			tt.validate(t, rec)
		})
	}
}

func TestHealth_NoDatabase(t *testing.T) {
	t.Parallel()
	// Handler with nil Database should still work for health check
	h := handler.New(nil, nil)

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	// Should not panic or error even without database
	assert.NotPanics(t, func() {
		h.Health(rec, req)
	})

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestHealth_DatabaseDown(t *testing.T) {
	t.Parallel()

	rawDB, err := sql.Open("sqlite", "file::memory:")
	require.NoError(t, err)
	require.NoError(t, rawDB.Close()) // force ping failures

	h := handler.New(slog.New(slog.DiscardHandler), db.NewFromRawDB(rawDB))

	req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	h.Health(rec, req)

	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
	assert.JSONEq(t, `{"status":"degraded","version":"dev"}`, rec.Body.String())
}

func TestHealth_MultipleRequests(t *testing.T) {
	t.Parallel()
	h := handler.New(slog.New(slog.DiscardHandler), nil)

	// Health endpoint should be idempotent
	for range 5 {
		req := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/api/health", nil)
		rec := httptest.NewRecorder()

		h.Health(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.JSONEq(t, `{"status":"ok","version":"dev"}`, rec.Body.String())
	}
}
