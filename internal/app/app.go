package app

import (
	"net/http"

	"github.com/chnmk/music-library-microservice/internal/transport"
)

// Запуск сервера
func Run() {
	// TODO: Вынести конфигурационные данные в .env-файл

	// TODO: запустить внутренние сервисы
	// muslib := new(services.MusicLibrary) // По-хорошему там абстракт билдер должен быть...

	// Запуск сервера
	http.HandleFunc("/library", transport.LibraryHandler)
	http.HandleFunc("/songs", transport.SongsHandler)
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		panic(err)
	}

	// TODO: Покрыть код debug- и info-логами

	// TODO: Shutdown
}
