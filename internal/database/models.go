// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	SenderID uuid.UUID
	Body     string
	SentAt   time.Time
	Hostname string
	Port     string
}

type Server struct {
	CreatorID uuid.UUID
	ServerID  uuid.UUID
	CreatedAt time.Time
	Hostname  string
	Port      string
}

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string
	Password  []byte
	Nickname  sql.NullString
}
