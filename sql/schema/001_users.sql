-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	username TEXT NOT NULL UNIQUE,
	password BYTEA NOT NULL,
	nickname TEXT
);

-- +goose Down
DROP TABLE users;