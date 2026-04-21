package main

import (
	"aliceSpeakerService/internal/config"
	apphttp "aliceSpeakerService/internal/http"
	"aliceSpeakerService/internal/store"
	"aliceSpeakerService/pkg/db"
	"log"
)

func main() {
	cfg := config.LoadConfig()
	repo := db.GetRepository(cfg.DBPath)
	if err := store.InitStore(repo); err != nil {
		log.Fatal(err)
	}

	server := apphttp.GetServer(cfg.Port)
	if err := server.StartHandle(); err != nil {
		log.Fatal(err)
	}
}
