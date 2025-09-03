-- name: CreateProject :one
INSERT INTO projects (
  name, description, "githubURL"
) VALUES (
  ?, ?, ?
)
RETURNING *;