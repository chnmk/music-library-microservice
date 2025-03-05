package app

import (
	"net/http"

	"github.com/chnmk/music-library-microservice/internal/services"
	"github.com/chnmk/music-library-microservice/internal/transport"
)

var MusLib services.MusicLibrary

// Запуск сервера
func Run() {
	// TODO: Вынести конфигурационные данные в .env-файл

	// TODO: запустить внутренние сервисы
	services.MusLib = services.NewLibrary()

	// Запуск сервера (TODO: можно вынести в internal/transport)
	http.HandleFunc("/library", transport.LibraryHandler)
	http.HandleFunc("/songs", transport.SongsHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}

	// TODO: Покрыть код debug- и info-логами

	// TODO: Shutdown
}
