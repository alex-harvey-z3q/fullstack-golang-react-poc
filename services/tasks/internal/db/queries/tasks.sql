-- name: ListTasks :many
SELECT id, title, done, created_at, updated_at FROM tasks ORDER BY id;

-- name: GetTask :one
SELECT id, title, done, created_at, updated_at FROM tasks WHERE id = $1;

-- name: CreateTask :one
INSERT INTO tasks (title) VALUES ($1)
RETURNING id, title, done, created_at, updated_at;

-- name: UpdateTask :one
UPDATE tasks SET title = $2, done = $3, updated_at = now()
WHERE id = $1
RETURNING id, title, done, created_at, updated_at;

-- name: DeleteTask :exec
DELETE FROM tasks WHERE id = $1;
