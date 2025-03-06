package transport

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/chnmk/music-library-microservice/internal/config"
)

// Запускает сервер, обрабатывает завершение его работы.
func StartServer(ctx context.Context) {
	http.HandleFunc("/library", LibraryHandler)
	http.HandleFunc("/songs", SongsHandler)

	server := &http.Server{Addr: ":" + config.SERVER_PORT, Handler: nil}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.Shutdown(ctx); err != nil {
			// log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		// log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

	slog.Info("shutting down...")
}

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией.
func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		libraryGet(w, r)
	}
}

// Добавление, получение текста, изменение, удаление конкретной песни.
func SongsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Добавление новой песни
		songsPost(w, r)
	case http.MethodGet:
		// Получение текста песни с пагинацией по куплетам
		songsGetLyrics(w, r)
	case http.MethodPut:
		// Изменение данных песни
		songsPut(w, r)
	case http.MethodDelete:
		// Удаление песни
		songsDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("expected POST, GET, PUT or DELETE"))
		return
	}
}
