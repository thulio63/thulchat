-- +goose Up
CREATE TABLE messages (
    sender_id UUID REFERENCES users(id) NOT NULL,
    body TEXT NOT NULL,
	sent_at TIMESTAMP NOT NULL,
	hostname TEXT NOT NULL,
	port TEXT NOT NULL,
	FOREIGN KEY(hostname, port) REFERENCES servers(hostname, port) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE messages;