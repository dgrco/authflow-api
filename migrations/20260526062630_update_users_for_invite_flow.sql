-- +goose Up
ALTER TABLE users
  ADD COLUMN invite_token TEXT,
  ADD COLUMN invite_expires_at TIMESTAMPTZ,
  ADD COLUMN invite_accepted_at TIMESTAMPTZ;

-- +goose Down
ALTER TABLE users
  DROP COLUMN invite_token,
  DROP COLUMN invite_expires_at,
  DROP COLUMN invite_accepted_at;
