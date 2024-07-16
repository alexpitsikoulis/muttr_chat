package handlers

import (
	"fmt"
	"muttr_chat/internal/app"
	"muttr_chat/internal/storage/repositories"
	"net/http"
)

type AppServer struct {
	Port                 uint
	AccountServer        app.ServerConfig
	WebSocketsServer     WebSocketsServer
	ThreadRepository     *repositories.ThreadRepository
	ThreadUserRepository *repositories.ThreadUserRepository
}

func (appServer *AppServer) ListenAndServe() error {
	fmt.Printf("Server is listening on port %d...\n", appServer.Port)
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", appServer.Port), nil)
}

type Route struct {
	Method string
	URL    string
}

func NewAppServer(cfg *app.Config, threadRepo *repositories.ThreadRepository, threadUserRepo *repositories.ThreadUserRepository) (*AppServer, error) {
	appServer := &AppServer{}
	if cfg == nil {
		return appServer, nil
	}
	accountServerUrl := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	wsServer, err := NewWebSocketsServer(threadRepo, threadUserRepo, accountServerUrl)
	if err != nil {
		return nil, err
	}
	appServer.buildRoutingTable(wsServer)
	return &AppServer{Port: cfg.Port, AccountServer: cfg.Server}, nil
}

func (server *AppServer) buildRoutingTable(wsServer *WebSocketsServer) {
	http.HandleFunc("/health-check", func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
	})
	http.HandleFunc("/ws", wsServer.HandleConnections)
}
