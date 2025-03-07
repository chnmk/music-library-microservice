package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/joho/godotenv"
)

var (
	MusLib   models.MusicLibrary
	Database models.Database
	EnvVars  map[string]string

	ServerPort         string
	DBConnectionString string

	RequestServer string
	RequestPort   string

	DBProtocol       string
	DBHost           string
	DBPort           string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string

	LogLevel  string
	SlogLevel slog.LevelVar

	MaxEntries int
	PageSize   int
)

// TODO: почистить (MaxEntries, непонятно что с LogLevel, нужны ли части строки)
func SetConfig() {
	EnvVars = make(map[string]string)
	EnvVars["SERVER_PORT"] = "3000"
	EnvVars["REQUEST_SERVER"] = "localhost"
	EnvVars["REQUEST_PORT"] = "3001"
	EnvVars["DB_PROTOCOL"] = "postgres"
	EnvVars["DB_HOST"] = "postgres"
	EnvVars["DB_PORT"] = "5432"
	EnvVars["POSTGRES_USER"] = "user"
	EnvVars["POSTGRES_PASSWORD"] = "12345"
	EnvVars["POSTGRES_DB"] = "orders"
	EnvVars["LOG_LEVEL"] = "debug"
	EnvVars["MAX_ENTRIES"] = "10000"
	EnvVars["PAGE_SIZE"] = "10"

	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found")
	} else {
		value, ok := os.LookupEnv("LOG_LEVEL")
		if ok {
			LogLevel = value
		}
	}

	if LogLevel == "debug" || LogLevel == "Debug" || LogLevel == "DEBUG" || LogLevel == "D" || LogLevel == "d" || LogLevel == "-4" {
		SlogLevel.Set(slog.LevelDebug)
		slog.Info("logging level set to 'debug' (-4)")
	} else {
		SlogLevel.Set(slog.LevelInfo)
		slog.Info("logging level set to 'info' (0)")
	}

	slog.Info("reading environment variables...")

	for name, def := range EnvVars {
		value, exists := os.LookupEnv(name)
		if exists {
			slog.Debug(
				"env variable found",
				"name", name,
				"value", value,
			)
			EnvVars[name] = value
		} else {
			slog.Debug(
				"env variable not found, using default",
				"name", name,
				"value", def,
			)
		}
	}

	ServerPort = EnvVars["SERVER_PORT"]
	RequestServer = EnvVars["REQUEST_SERVER"]
	RequestPort = EnvVars["REQUEST_PORT"]
	DBProtocol = EnvVars["DB_PROTOCOL"]
	DBHost = EnvVars["DB_HOST"]
	DBPort = EnvVars["DB_PORT"]
	PostgresUser = EnvVars["POSTGRES_USER"]
	PostgresPassword = EnvVars["POSTGRES_PASSWORD"]
	PostgresDB = EnvVars["POSTGRES_DB"]
	// MaxEntries = EnvVars["MAX_ENTRIES"]
	PageSize = 10 // TODO

	DBConnectionString = fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		EnvVars["DB_PROTOCOL"],
		EnvVars["POSTGRES_USER"],
		EnvVars["POSTGRES_PASSWORD"],
		EnvVars["DB_HOST"],
		EnvVars["DB_PORT"],
		EnvVars["POSTGRES_DB"],
	)

}
