-- name: GetAuthor :one
SELECT * FROM authors
WHERE id = ? LIMIT 1;

-- name: ListAuthors :many
SELECT * FROM authors
ORDER BY name;

-- name: CreateAuthor :one
INSERT INTO authors (
  name, bio
) VALUES (
  ?, ?
)
RETURNING *;

-- name: UpdateAuthor :exec
UPDATE authors
SET name = ?,
bio = ?
WHERE id = ?;

-- name: DeleteAuthor :exec
DELETE FROM authors
WHERE id = ?;

-- name: ListTodos :many
SELECT * FROM todos
ORDER BY done ASC, created_at DESC;

-- name: CountTodos :one
SELECT COUNT(*) FROM todos;

-- name: GetTodo :one
SELECT * FROM todos
WHERE id = ? LIMIT 1;

-- name: CreateTodo :one
INSERT INTO todos (
  title
) VALUES (
  ?
)
RETURNING *;

-- name: UpdateTodoTitle :one
UPDATE todos
SET title = ?, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: ToggleTodo :one
UPDATE todos
SET done = NOT done, updated_at = CURRENT_TIMESTAMP
WHERE id = ?
RETURNING *;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = ?;
