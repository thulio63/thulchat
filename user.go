package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Username string
	Nickname string
	Password string
}