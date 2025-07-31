-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (
  user_id
) VALUES ( sqlc.arg(userID)::uuid )
RETURNING *;


-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
  SET revoked_at = NOW(), updated_at = NOW()
  WHERE token = sqlc.arg(token)::text
  AND revoked_at IS NULL
  AND expires_at > NOW()
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.*
  FROM users
  JOIN refresh_tokens
  ON users.id = refresh_tokens.user_id
  WHERE refresh_tokens.token = sqlc.arg(token)::text
  AND revoked_at IS NULL
  AND expires_at > NOW();
