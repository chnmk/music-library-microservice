package config

import (
	"log/slog"
	"os"

	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/joho/godotenv"
)

var (
	MusLib      models.MusicLibrary
	Database    models.Database
	SERVER_PORT string
)

func SetConfig() {
	// Значения по умолчанию
	SERVER_PORT = "3000"

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found")
	}

	var value string
	var ok bool

	value, ok = os.LookupEnv("SERVER_PORT")
	if ok {
		slog.Debug("var exists")
		SERVER_PORT = value
	} else {
		slog.Debug("var doesn't exist")
	}
}
