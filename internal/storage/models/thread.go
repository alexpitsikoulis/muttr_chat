package models

import (
	"time"

	"github.com/google/uuid"
)

type ThreadType int

const (
	DirectMessage ThreadType = iota
	ServerThread
	GroupMessage
)

type Thread struct {
	Id           uuid.UUID  `db:"id"`
	ThreadType   ThreadType `db:"thread_type"`
	ServerId     *uuid.UUID `db:"server_id"`
	Name         *string    `db:"name"`
	VoiceEnabled bool       `db:"voice_enabled"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

func NewThread(
	id uuid.UUID,
	serverId *uuid.UUID,
	name *string,
	voiceEnabled bool,
) *Thread {
	now := time.Now()
	return &Thread{
		Id:           id,
		ServerId:     serverId,
		Name:         name,
		VoiceEnabled: voiceEnabled,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}
