package app

import (
	"context"
	"fmt"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/database"
	"github.com/chnmk/music-library-microservice/internal/memory"
	"github.com/chnmk/music-library-microservice/internal/transport"
)

func Run() {
	// TODO: Покрыть код debug- и info-логами

	config.MusLib = memory.NewLibrary()
	config.Database = database.NewDatabase()

	transport.StartServer(context.Background())
	fmt.Println("shutdown successful")
}
