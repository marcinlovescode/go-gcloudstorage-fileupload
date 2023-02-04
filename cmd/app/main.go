package main

import (
	"fmt"
	"log"
	"os"

	"github.com/marcinlovescode/go-gcloudstorage-fileupload/config"
	"github.com/marcinlovescode/go-gcloudstorage-fileupload/internal/app"
)

func main() {
	configPath := "./config/config.yml"
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatalf("main - Config error: %s", err)
	}
	if cfg.GCloudStorage.UseEmulator {
		err := os.Setenv("STORAGE_EMULATOR_HOST", fmt.Sprintf("http://localhost:%d", cfg.GCloudStorage.EmulatorPort))
		if err != nil {
			log.Fatalf("main - Can't set GCloud Emulator Port: %s", err)
		}
	}
	app.Run(cfg)
}
