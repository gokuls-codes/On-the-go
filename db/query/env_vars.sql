-- name: CreateEnvVar :one
INSERT INTO env_vars (
  project_id, key, value
) VALUES (
  ?, ?, ?
)
RETURNING *;