-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1; 

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET updated_at = $1, revoked_at = $2
WHERE token = $3
RETURNING token;