package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/database"
	"github.com/chnmk/music-library-microservice/internal/memory"
	"github.com/chnmk/music-library-microservice/internal/transport"
)

func Run() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	slog.Info("initialization start...")

	config.SetConfig()
	config.MusLib = memory.NewLibrary()
	config.Database = database.NewDatabase()

	slog.Info("initialization complete, starting server...")

	transport.StartServer(context.Background())

	slog.Info("shutdown complete")
}
