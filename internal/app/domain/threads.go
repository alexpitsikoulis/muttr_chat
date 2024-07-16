package domain

import (
	"fmt"
	"muttr_chat/internal/storage/models"
	"muttr_chat/internal/storage/repositories"
	"sync"

	"github.com/google/uuid"
)

type Thread struct {
	Id           uuid.UUID
	ThreadType   models.ThreadType
	ServerId     *uuid.UUID
	Name         *string
	Users        map[uuid.UUID]*User
	VoiceEnabled bool
	Broadcast    chan []byte
	Register     chan *User
	Unregister   chan *User
	Connect      chan *User
	sync.Mutex
}

func (t *Thread) Run(threadRepo *repositories.ThreadRepository, threadUserRepo *repositories.ThreadUserRepository) {
	for {
		select {
		case user := <-t.Register:
			// TODO: add mutex and use channels for changes to thread struct to ensure thread safety
			err := t.AddUser(threadRepo, threadUserRepo, user)
			if err != nil {
				// TODO: implement error channel
			}
		case user := <-t.Unregister:
			// TODO: add mutex and use channels for changes to thread struct to ensure thread safety
			err := t.RemoveUser(threadUserRepo, user.Id, t.Id)
			if err != nil {
				// TODO: implement error channel
			}
		case msg := <-t.Broadcast:
			for _, user := range t.Users {
				select {
				case user.Send <- msg:
				default:
					close(user.Send)
					delete(t.Users, user.Id)
				}
			}
		}
	}
}

func StartNewDMThread(threadRepo *repositories.ThreadRepository, name *string, voiceEnabled bool) (*Thread, error) {
	// TODO: create thread user entities for both parties
	return startNewThread(threadRepo, models.DirectMessage, nil, name, voiceEnabled)
}

func StartNewServerThread(threadRepo *repositories.ThreadRepository, serverId uuid.UUID, name string, voiceEnabled bool) (*Thread, error) {
	// TODO: create thread user entities for all parties (TBD whether all server members are in by default)
	return startNewThread(threadRepo, models.ServerThread, &serverId, &name, voiceEnabled)
}

func StartNewGroupMessageThread(threadRepo *repositories.ThreadRepository, name *string, voiceEnabled bool) (*Thread, error) {
	// TODO: create thread user entities for all parties
	return startNewThread(threadRepo, models.GroupMessage, nil, name, voiceEnabled)
}

func OpenThread(threadRepo *repositories.ThreadRepository, threadUserRepo *repositories.ThreadUserRepository, id uuid.UUID) (*Thread, error) {
	thread, err := threadRepo.GetById(id)
	if err != nil {
		return nil, err
	}
	users, err := populateThreadUsers(threadUserRepo, id)
	if err != nil {
		return nil, err
	}
	return &Thread{
		Id:           thread.Id,
		ServerId:     thread.ServerId,
		Name:         thread.Name,
		VoiceEnabled: thread.VoiceEnabled,
		Users:        users,
		Broadcast:    make(chan []byte),
		Register:     make(chan *User),
		Unregister:   make(chan *User),
	}, nil
}

func (t *Thread) AddUser(threadRepo *repositories.ThreadRepository, threadUserRepo *repositories.ThreadUserRepository, user *User) error {
	tx, err := threadUserRepo.DB.Begin()
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	var newThread *models.Thread
	if t.ThreadType == models.DirectMessage {
		newThread = &models.Thread{
			Id:           uuid.New(),
			ThreadType:   models.GroupMessage,
			ServerId:     t.ServerId,
			Name:         t.Name,
			VoiceEnabled: t.VoiceEnabled,
		}
		err = threadRepo.Upsert(tx, newThread)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create group message from direct message thread: %v", err)
		}
	}

	err = threadUserRepo.Upsert(tx, &models.ThreadUser{UserId: user.Id, ThreadId: t.Id, UserRole: user.UserRole})
	if err != nil {
		tx.Rollback()
		return err
	}

	if newThread != nil {
		t.Id = newThread.Id
		t.ThreadType = models.GroupMessage
		t.Broadcast = make(chan []byte)
		t.Register = make(chan *User)
		t.Unregister = make(chan *User)
	}

	t.Users[user.Id] = user
	return nil
}

func (t *Thread) RemoveUser(threadUserRepo *repositories.ThreadUserRepository, userId, threadId uuid.UUID) error {
	tx, err := threadUserRepo.DB.Begin()
	defer tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	err = threadUserRepo.Delete(tx, userId, t.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	delete(t.Users, userId)
	return nil
}

func startNewThread(threadRepo *repositories.ThreadRepository, threadType models.ThreadType, serverId *uuid.UUID, name *string, voiceEnabled bool) (*Thread, error) {
	thread := &Thread{
		Id:           uuid.New(),
		ThreadType:   threadType,
		ServerId:     serverId,
		Name:         name,
		Users:        make(map[uuid.UUID]*User),
		VoiceEnabled: voiceEnabled,
		Broadcast:    make(chan []byte),
		Register:     make(chan *User),
		Unregister:   make(chan *User),
	}
	tx, err := threadRepo.DB.Begin()
	defer tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %v", err)
	}
	err = threadRepo.Upsert(tx, models.NewThread(
		thread.Id,
		thread.ServerId,
		thread.Name,
		thread.VoiceEnabled,
	))
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return thread, nil
}

func populateThreadUsers(threadUserRepo *repositories.ThreadUserRepository, threadId uuid.UUID) (map[uuid.UUID]*User, error) {
	users, err := threadUserRepo.GetManyByThreadId(threadId)
	if err != nil {
		return nil, err
	}
	userMap := make(map[uuid.UUID]*User)
	for _, u := range users {
		userMap[u.UserId] = nil
	}
	return userMap, nil
}
