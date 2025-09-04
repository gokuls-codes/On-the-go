-- name: CreateProject :one
INSERT INTO projects (
  name, description, github_url, repo_name, container_port, host_port
) VALUES (
  ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateRepoName :exec
UPDATE projects SET
  repo_name = ?, 
  updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: UpdateImageId :exec
UPDATE projects SET 
  image_id = ?, 
  updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: UpdateContainerId :exec 
UPDATE projects SET 
  container_id = ?, 
  updated_at = CURRENT_TIMESTAMP 
WHERE id = ?;

-- name: GetProjectByRepoName :one
SELECT * FROM projects WHERE repo_name = ?;