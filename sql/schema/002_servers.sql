-- +goose Up
CREATE TABLE servers (
    creator_id UUID REFERENCES users(id) NOT NULL,
    server_id UUID NOT NULL,
	created_at TIMESTAMP NOT NULL,
	hostname TEXT NOT NULL,
	port TEXT NOT NULL,
	PRIMARY KEY(hostname, port)
);

-- +goose Down
DROP TABLE servers;