package domain

import (
	"fmt"
	"log"
	"muttr_chat/internal/storage/models"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type User struct {
	Id       uuid.UUID
	UserRole models.Role
	Conn     *websocket.Conn
	Send     chan []byte
}

func (u *User) ReadMessages() {
	defer u.Conn.Close()

	for {
		_, msg, err := u.Conn.ReadMessage()
		if err != nil {
			log.Println("failed to read messaged: ", err)
			break
		}
		fmt.Println(msg)
	}
}

func (u *User) WriteMessages() {
	defer u.Conn.Close()

	for msg := range u.Send {
		err := u.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("failed to write message: ", err)
			break
		}
		// TODO: store messages
	}
}
