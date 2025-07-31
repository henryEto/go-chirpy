-- +goose Up
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY NOT NULL DEFAULT encode(gen_random_bytes(32), 'hex'),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '60 day',
    revoked_at TIMESTAMP DEFAULT NULL

);

-- +goose Down
DROP TABLE IF EXISTS refresh_tokens;
