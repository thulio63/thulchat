-- +goose Up
CREATE TABLE messages (
    sender_id UUID REFERENCES users(id) NOT NULL,
    body TEXT NOT NULL,
	sent_at TIMESTAMP NOT NULL,
	hostname TEXT NOT NULL,
	port TEXT NOT NULL,
	server_id UUID NOT NULL UNIQUE REFERENCES servers (server_id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE messages;