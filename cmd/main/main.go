package main

import (
	"log"

	"github.com/dorik33/medods_test/internal/api"
	"github.com/dorik33/medods_test/internal/config"

)

func main() {
	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	apiServer := api.New(cfg)
	if err := apiServer.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
