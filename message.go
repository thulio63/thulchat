package main

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID uuid.UUID
	Body string
	TimeSent time.Time
	Sender uuid.UUID
}