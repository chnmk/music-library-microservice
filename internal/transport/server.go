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

	server := &http.Server{Addr: ":" + config.ServerPort, Handler: nil}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := server.Shutdown(ctx); err != nil {
			slog.Info(
				"HTTP server Shutdown",
				"err", err,
			)
		}
		close(idleConnsClosed)
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		slog.Info(
			"HTTP server ListenAndServe",
			"err", err,
		)
	}

	<-idleConnsClosed

	slog.Info("shutting down...")
}

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией.
func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(
		"new request incoming...",
		"path", "/library",
		"method", r.Method,
	)

	if r.Method == http.MethodGet {
		libraryGet(w, r)
	} else {
		slog.Error("request denied", "err", "expected GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("expected GET"))
		return
	}

	slog.Debug(
		"request finished",
		"path", "/library",
		"method", r.Method,
	)
}

// Добавление, получение текста, изменение, удаление конкретной песни.
func SongsHandler(w http.ResponseWriter, r *http.Request) {
	slog.Debug(
		"new request incoming...",
		"path", "/songs",
		"method", r.Method,
	)

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
		slog.Error("request denied", "err", "expected POST, GET, PUT or DELETE")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("expected POST, GET, PUT or DELETE"))
		return
	}

	slog.Debug(
		"request finished",
		"path", "/songs",
		"method", r.Method,
	)
}
