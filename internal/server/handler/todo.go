package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/core"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/layout"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/components/todo"
	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/db/queries"
)

// TodoList handles the todo list page.
func (h *Handler) TodoList(w http.ResponseWriter, r *http.Request) {
	todos, err := h.database.Queries().ListTodos(r.Context())
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}

	page := layout.Shell("todo-list-db", todo.Page(todos))
	h.html(r.Context(), w, core.HTML("go-htmx-sqlite-ai — Task queue", page))
}

// TodoCreate creates a todo and returns the refreshed list fragment.
func (h *Handler) TodoCreate(w http.ResponseWriter, r *http.Request) {
	title := strings.TrimSpace(r.FormValue("title"))
	if title != "" {
		if _, err := h.database.Queries().CreateTodo(r.Context(), title); err != nil {
			h.logger.Error("failed to create todo", "error", err)
			http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
	h.renderTodoList(w, r)
}

// TodoToggle flips a todo's done state and returns the refreshed list fragment.
func (h *Handler) TodoToggle(w http.ResponseWriter, r *http.Request) {
	id, ok := h.todoID(w, r)
	if !ok {
		return
	}
	if _, err := h.database.Queries().ToggleTodo(r.Context(), id); err != nil {
		h.logger.Error("failed to toggle todo", "error", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.renderTodoList(w, r)
}

// TodoEdit returns the inline-edit form for a todo.
func (h *Handler) TodoEdit(w http.ResponseWriter, r *http.Request) {
	id, ok := h.todoID(w, r)
	if !ok {
		return
	}
	t, err := h.database.Queries().GetTodo(r.Context(), id)
	if err != nil {
		h.handleGetTodoError(w, err)
		return
	}
	h.html(r.Context(), w, todo.EditRow(t))
}

// TodoGet returns the display row for a todo; used to cancel an in-progress edit.
func (h *Handler) TodoGet(w http.ResponseWriter, r *http.Request) {
	id, ok := h.todoID(w, r)
	if !ok {
		return
	}
	t, err := h.database.Queries().GetTodo(r.Context(), id)
	if err != nil {
		h.handleGetTodoError(w, err)
		return
	}
	h.html(r.Context(), w, todo.Row(t))
}

// TodoUpdate saves an edited title and returns the refreshed list fragment.
func (h *Handler) TodoUpdate(w http.ResponseWriter, r *http.Request) {
	id, ok := h.todoID(w, r)
	if !ok {
		return
	}
	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "422 title required", http.StatusUnprocessableEntity)
		return
	}
	if _, err := h.database.Queries().UpdateTodoTitle(r.Context(), queries.UpdateTodoTitleParams{Title: title, ID: id}); err != nil {
		h.logger.Error("failed to update todo", "error", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.renderTodoList(w, r)
}

// TodoDelete deletes a todo and returns the refreshed list fragment.
func (h *Handler) TodoDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := h.todoID(w, r)
	if !ok {
		return
	}
	if err := h.database.Queries().DeleteTodo(r.Context(), id); err != nil {
		h.logger.Error("failed to delete todo", "error", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.renderTodoList(w, r)
}

func (h *Handler) renderTodoList(w http.ResponseWriter, r *http.Request) {
	todos, err := h.database.Queries().ListTodos(r.Context())
	if err != nil {
		h.logger.Error("failed to list todos", "error", err)
		http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
		return
	}
	h.html(r.Context(), w, todo.List(todos))
}

func (h *Handler) handleGetTodoError(w http.ResponseWriter, err error) {
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "404 Not Found", http.StatusNotFound)
		return
	}
	h.logger.Error("failed to get todo", "error", err)
	http.Error(w, "500 Internal Server Error", http.StatusInternalServerError)
}

func (h *Handler) todoID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "400 Bad Request", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}
