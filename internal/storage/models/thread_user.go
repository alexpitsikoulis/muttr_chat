package models

import "github.com/google/uuid"

type Role int

const (
	User Role = iota
	Moderator
	Admin
)

type ThreadUser struct {
	UserId   uuid.UUID `db:"user_id"`
	ThreadId uuid.UUID `db:"thread_id"`
	UserRole Role      `db:"user_role"`
}

func NewTreadUser(threadId uuid.UUID, userId uuid.UUID, role Role) *ThreadUser {
	return &ThreadUser{
		ThreadId: threadId,
		UserId:   userId,
		UserRole: role,
	}
}
