package handlers

import (
	"fmt"
	"log"
	"muttr_chat/internal/app/domain"
	"muttr_chat/internal/storage/repositories"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WebSocketsServer struct {
	AccountServerUrl     string
	Clients              map[uuid.UUID]*domain.User
	Threads              map[uuid.UUID]*domain.Thread
	ThreadRepository     *repositories.ThreadRepository
	ThreadUserRepository *repositories.ThreadUserRepository
	sync.Mutex
}

func NewWebSocketsServer(threadRepository *repositories.ThreadRepository, threadUserRepository *repositories.ThreadUserRepository, accountServerUrl string) (*WebSocketsServer, error) {
	server := &WebSocketsServer{AccountServerUrl: accountServerUrl}
	threads, err := threadRepository.GetAll()
	if err != nil {
		return nil, err
	}
	for _, t := range threads {
		thread, err := domain.OpenThread(server.ThreadRepository, server.ThreadUserRepository, t.Id)
		if err != nil {
			return nil, err
		}
		server.Threads[thread.Id] = thread
		go thread.Run(server.ThreadRepository, server.ThreadUserRepository)
	}
	return server, nil
}

func (server *WebSocketsServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	if auth := r.Header.Get("Authorization"); auth != "" {
		token, err := jwt.Parse(auth[6:], func(token *jwt.Token) (interface{}, error) {
			if err = token.Claims.Valid(); err != nil {
				return nil, err
			}
			return token, nil
		})
		if err != nil {
			err := fmt.Errorf("invalid token provided: %w", err)
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		claims := token.Claims.(*jwt.StandardClaims)
		if claims.Subject == "" {
			err := "subject field of token is blank"
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err))
			return
		}

		user := domain.User{
			Conn: conn,
			Send: make(chan []byte),
		}

		id, err := uuid.Parse(claims.Subject)
		if err != nil {
			err := fmt.Errorf("token subject field is not a valid uuid: %w", err)
			log.Println(err)
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte(err.Error()))
		}
		server.Clients[id] = &user

	}

}
