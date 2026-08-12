-- +goose Up
CREATE TABLE chirps (
  id uuid PRIMARY KEY,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  body text NOT NULL,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;

