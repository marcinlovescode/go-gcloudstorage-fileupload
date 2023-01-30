package main

import (
	"log"

	"github.com/marcinlovescode/go-clean-fileupload/config"
	"github.com/marcinlovescode/go-clean-fileupload/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	app.Run(cfg)
}
