package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/joho/godotenv"
)

var (
	// Глобальные объекты.
	MusLib   models.MusicLibrary
	Database models.Database
	EnvVars  map[string]string

	// Объекты для завершения работы.
	Exit    context.CancelFunc
	ExitCtx context.Context

	// Порт сервера, строка подключения к БД, информация о API с текстами песен.
	ServerPort         string
	DBConnectionString string
	RequestServer      string
	RequestTimeout     int

	/*
		// Переменные, которые не используются отдельно от строки подключения к БД.

		DBProtocol       string
		DBHost           string
		DBPort           string
		PostgresUser     string
		PostgresPassword string
		PostgresDB       string
	*/

	// Уровень логирования.
	LogLevel  string
	SlogLevel slog.LevelVar

	// Настройки сервиса.
	MaxEntries  int
	RestoreData bool
	PageSize    int
)

// Устанавливает глобальные значения на основе переменных окружения.
// TODO: если переменных станет слишком много, можно разбить на несколько функций или файлов.
func SetConfig() {
	// Стандартные значения переменных окружения.
	EnvVars = make(map[string]string)
	EnvVars["SERVER_PORT"] = "3000"
	EnvVars["REQUEST_SERVER"] = "http://localhost:3001"
	EnvVars["REQUEST_TIMEOUT"] = "1"
	EnvVars["POSTGRES_USER"] = "user"
	EnvVars["POSTGRES_PASSWORD"] = "12345"
	EnvVars["POSTGRES_DB"] = "muslib"
	EnvVars["DB_PROTOCOL"] = "postgres"
	EnvVars["DB_HOST"] = "postgres"
	EnvVars["DB_PORT"] = "5432"
	EnvVars["SSL_MODE"] = "disable"
	EnvVars["LOG_LEVEL"] = "debug"
	EnvVars["MAX_ENTRIES"] = "10000"
	EnvVars["RESTORE_FROM_DB"] = "true"
	EnvVars["PAGE_SIZE"] = "10"

	// Поиск .env файла.
	err := godotenv.Load()
	if err != nil {
		slog.Info(".env file not found")
		LogLevel = "debug"
	} else {
		value, ok := os.LookupEnv("LOG_LEVEL")
		if ok {
			LogLevel = value
		}
	}

	// Установка уровня логгирования.
	if LogLevel == "debug" || LogLevel == "Debug" || LogLevel == "DEBUG" || LogLevel == "D" || LogLevel == "d" || LogLevel == "-4" {
		SlogLevel.Set(slog.LevelDebug)
		slog.Info("logging level set to 'debug' (-4)")
	} else {
		SlogLevel.Set(slog.LevelInfo)
		slog.Info("logging level set to 'info' (0)")
	}

	slog.Info("reading environment variables...")

	// Чтение переменных окружения.
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

	// Запись переменных окружения в переменные config.
	ServerPort = EnvVars["SERVER_PORT"]
	RequestServer = EnvVars["REQUEST_SERVER"]

	/*
		DBProtocol = EnvVars["DB_PROTOCOL"]
		DBHost = EnvVars["DB_HOST"]
		DBPort = EnvVars["DB_PORT"]
		PostgresUser = EnvVars["POSTGRES_USER"]
		PostgresPassword = EnvVars["POSTGRES_PASSWORD"]
		PostgresDB = EnvVars["POSTGRES_DB"]
	*/

	// Запись переменных int и bool. Если станет слишком много, можно вынести в отдельную функцию.
	val, err := strconv.Atoi(EnvVars["REQUEST_TIMEOUT"])
	if err != nil {
		slog.Debug(
			"error converting env var to int, using default",
			"name", "REQUEST_TIMEOUT",
			"value", EnvVars["REQUEST_TIMEOUT"],
		)
		RequestTimeout = 1
	} else {
		RequestTimeout = val
	}

	val, err = strconv.Atoi(EnvVars["REQUEST_TIMEOUT"])
	if err != nil {
		slog.Debug(
			"error converting env var to int, using default",
			"name", "REQUEST_TIMEOUT",
			"value", EnvVars["REQUEST_TIMEOUT"],
		)
		RequestTimeout = 1
	} else {
		RequestTimeout = val
	}

	val, err = strconv.Atoi(EnvVars["MAX_ENTRIES"])
	if err != nil {
		slog.Debug(
			"error converting env var to int, using default",
			"name", "MAX_ENTRIES",
			"value", EnvVars["MAX_ENTRIES"],
		)
		MaxEntries = 10000
	} else {
		MaxEntries = val
	}

	val, err = strconv.Atoi(EnvVars["PAGE_SIZE"])
	if err != nil {
		slog.Debug(
			"error converting env var to int, using default",
			"name", "PAGE_SIZE",
			"value", EnvVars["PAGE_SIZE"],
		)
		PageSize = 10
	} else {
		PageSize = val
	}

	valBool, err := strconv.ParseBool(EnvVars["RESTORE_FROM_DB"])
	if err != nil {
		slog.Debug(
			"error converting env var to bool, using default",
			"name", "RESTORE_FROM_DB",
			"value", EnvVars["RESTORE_FROM_DB"],
		)
		RestoreData = false
	} else {
		RestoreData = valBool
	}

	// Создание строки подключения к БД.
	DBConnectionString = fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		EnvVars["DB_PROTOCOL"],
		EnvVars["POSTGRES_USER"],
		EnvVars["POSTGRES_PASSWORD"],
		EnvVars["DB_HOST"],
		EnvVars["DB_PORT"],
		EnvVars["POSTGRES_DB"],
	)

	if EnvVars["SSL_MODE"] == "disable" {
		DBConnectionString = DBConnectionString + "?sslmode=disable"
	}

	// Создание контекста для завершения работы.
	ExitCtx, Exit = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
}
