package app

import (
	"log/slog"
	"os"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/database"
	"github.com/chnmk/music-library-microservice/internal/services"
	"github.com/chnmk/music-library-microservice/internal/transport"
)

// Запускает логгер, конфиг, базу данных, бизнес-логику и http-сервер.
func Run() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: &config.SlogLevel})))

	slog.Info("initialization start...")

	config.SetConfig()
	config.Database = database.NewDatabase(config.ExitCtx)
	config.MusLib = services.NewLibrary()

	slog.Info("initialization complete, starting server...")

	transport.StartServer(config.ExitCtx)

	slog.Info("shutdown complete")
}
