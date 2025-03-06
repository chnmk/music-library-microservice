package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"

	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/joho/godotenv"
)

var (
	MusLib    models.MusicLibrary
	Database  models.Database
	SlogLevel slog.LevelVar

	ServerPort string

	RequestServer string
	RequestPort   string

	DBProtocol       string
	DBHost           string
	DBPort           string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	LogLevel   string
	MaxEntries int
)

func SetConfig() {
	// Значения по умолчанию
	ServerPort = "3000"
	RequestServer = "localhost"
	RequestPort = "3001"
	DBProtocol = "postgres"
	DBHost = "postgres"
	DBPort = "5432"
	PostgresUser = "user"
	PostgresPassword = "12345"
	PostgresDB = "orders"
	LogLevel = "debug"
	MaxEntries = 10000

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found")
	} else {
		slog.Info("reading from .env file...")

		value, ok := os.LookupEnv("LOG_LEVEL")
		if ok {
			if value == "info" || value == "Info" || value == "INFO" || value == "I" || value == "i" || value == "0" {
				SlogLevel.Set(slog.LevelInfo)
				slog.Info("logging level set to 'info' (0)")
			}

			if value == "debug" || value == "Debug" || value == "DEBUG" || value == "D" || value == "d" || value == "-4" {
				SlogLevel.Set(slog.LevelDebug)
				slog.Info("logging level set to 'debug' (-4)")
			}
		}

		// TODO: Если значений станет слишком много, можно использовать мапу и получать данные в цикле.
		ServerPort = getEnv("SERVER_PORT")
		RequestServer = getEnv("REQUEST_SERVER")
		RequestPort = getEnv("REQUEST_PORT")
		DBProtocol = getEnv("DB_PROTOCOL")
		DBHost = getEnv("DB_HOST")
		DBPort = getEnv("DB_PORT")
		PostgresUser = getEnv("POSTGRES_USER")
		PostgresPassword = getEnv("POSTGRES_PASSWORD")
		PostgresDB = getEnv("POSTGRES_DB")

		MaxEntries, err = strconv.Atoi(getEnv("MAX_ENTRIES"))
		if err != nil {
			slog.Debug(
				fmt.Sprintf("error converting env variable to int: %v", err),
				"name", "MAX_ENTRIES",
			)
		}
	}
}

func getEnv(name string) string {
	value, ok := os.LookupEnv(name)
	if ok {
		slog.Debug(
			"env variable exists",
			"name", name,
			"value", value,
		)
	} else {
		slog.Debug(
			"env variable doesn't exists",
			"name", name,
		)
	}

	return value
}
