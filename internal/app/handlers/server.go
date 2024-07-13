package handlers

import (
	"fmt"
	"muttr_chat/internal/app"
	"net/http"
)

type AppServer struct {
	Port uint
}

func (appServer *AppServer) ListenAndServe() error {
	return http.ListenAndServe(fmt.Sprintf("localhost:%d", appServer.Port), nil)
}

type Route struct {
	Method string
	URL    string
}

func NewAppServer(cfg *app.Config) (*AppServer, error) {
	if cfg == nil {
		return &AppServer{}, nil
	}
	buildRoutingTable()
	return &AppServer{Port: cfg.Port}, nil
}

func buildRoutingTable() {
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Println("hello world")
	})
}
