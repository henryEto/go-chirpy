-- name: DeleteChirp :exec
DELETE FROM chirps
  WHERE id = sqlc.arg(chirpID)::uuid;

-- name: GetChirpByID :one
SELECT *
  FROM chirps
  WHERE id = sqlc.arg(id)::uuid;

-- name: GetAllChirps :many
SELECT *
  FROM chirps
  ORDER BY created_at ASC;

-- name: PostChirp :one
INSERT INTO chirps (
  body, user_id
) VALUES ( sqlc.arg(body)::text, sqlc.arg(userID)::uuid )
RETURNING *;

-- name: ResetChirps :exec
DELETE FROM chirps;
