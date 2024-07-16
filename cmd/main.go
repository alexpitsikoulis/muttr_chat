package main

import (
	"fmt"
	"log"
	"muttr_chat/internal/app"
	"muttr_chat/internal/storage/db"
	"muttr_chat/internal/storage/repositories"

	"muttr_chat/internal/app/handlers"
)

func main() {
	cfg, err := app.NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.InitDB(fmt.Sprintf("user=%s password=%s sslmode=disable", cfg.DB.Username, cfg.DB.Password)); err != nil {
		log.Fatal(err)
	}

	threadRepo := &repositories.ThreadRepository{DB: db.DB}
	threadUserRepo := &repositories.ThreadUserRepository{DB: db.DB}

	appServer, err := handlers.NewAppServer(cfg, threadRepo, threadUserRepo)
	if err != nil {
		log.Fatal(err)
	}
	appServer.ListenAndServe()
}
