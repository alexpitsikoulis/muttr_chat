package main

import (
	"log"
	"muttr_chat/internal/app"

	"muttr_chat/internal/app/handlers"
)

func main() {
	cfg, err := app.NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	appServer, err := handlers.NewAppServer(cfg)
	if err != nil {
		log.Fatal(err)
	}
	appServer.ListenAndServe()
}
