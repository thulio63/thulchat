package main

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Body string `json:"body"`
	TimeSent time.Time `json:"timeSent"`
	Sender uuid.UUID `json:"sender"`
}