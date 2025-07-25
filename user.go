package main

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Email string
	Nickname string
	Password string
}