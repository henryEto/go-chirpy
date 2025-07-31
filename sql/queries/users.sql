-- name: UpgradeUser :exec
UPDATE users
  SET is_chirpy_red = true
  WHERE id = sqlc.arg(userID)::uuid;

-- name: GetUserByID :one
SELECT *
  FROM users
  WHERE id = sqlc.arg(userID)::uuid;

-- name: UpdateUser :one
UPDATE users
  SET hashed_password = sqlc.arg(hashedPassword)::text,
  email = sqlc.arg(email)::text,
  updated_at = NOW()
  WHERE id = sqlc.arg(userID)::uuid
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
  FROM users
  WHERE email = sqlc.arg(email)::text;

-- name: CreateUser :one
INSERT INTO users (
  email,
  hashed_password
) VALUES ( sqlc.arg(email)::text, sqlc.arg(hashedPassword)::text )
RETURNING *;

-- name: ResetUsers :exec
DELETE FROM users;
