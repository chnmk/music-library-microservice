package config

import (
	"fmt"
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
		fmt.Println(".env file not found")
	}

	var ok bool

	SERVER_PORT, ok = os.LookupEnv("SERVER_PORT")
	if ok {
		fmt.Println("var exists")
	} else {
		fmt.Println("var doesn't exist")
	}
}
