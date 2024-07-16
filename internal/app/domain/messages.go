package domain

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Id        uuid.UUID
	ThreadId  uuid.UUID
	SenderId  uuid.UUID
	Message   string
	ReadBy    []uuid.UUID
	Reactions []Reaction
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

type Reaction struct {
	UserId   uuid.UUID
	Reaction [4]byte
}
