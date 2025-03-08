package transport

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/chnmk/music-library-microservice/internal/services"
)

func TestInvalidLibraryMethods(t *testing.T) {
	req := httptest.NewRequest("POST", "/library", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got: %d, error: %s", rec.Code, rec.Body)
	}
}

func TestInvalidSongsMethods(t *testing.T) {
	req := httptest.NewRequest("PATCH", "/songs", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got: %d, error: %s", rec.Code, rec.Body)
	}
}

func TestRequestNoData(t *testing.T) {
	config.SetConfig()
	config.Database = models.MockDatabase{Data: make(map[int]models.SongData)}
	config.MusLib = services.NewLibrary()

	req := httptest.NewRequest("GET", "/library", nil)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got: %d, error: %s", rec.Code, rec.Body)
	}
}

func TestCRUDRequests(t *testing.T) {
	config.SetConfig()
	config.Database = models.MockDatabase{Data: make(map[int]models.SongData)}
	config.MusLib = services.NewLibrary()

	// Добавление песни
	body := strings.NewReader(`{
		"group": "Muse", 
		"song": "Supermassive Black Hole"
	}`)

	req := httptest.NewRequest("POST", "/songs", body)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(SongsHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	// Получение песен
	req = httptest.NewRequest("GET", "/library", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if !strings.Contains(rec.Body.String(), "Supermassive Black Hole") {
		t.Errorf("data not found, got %s", rec.Body.String())
	}

	// Изменение песен
	body = strings.NewReader(`{
		"group": "Muzlo", 
		"song": "A Very Tiny Black Hole",
		"lyrics": "Verse 1\n\nVerse2"
	}`)

	req = httptest.NewRequest("PUT", "/songs?id=0", body)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(SongsHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	// Получение текста песен
	req = httptest.NewRequest("GET", "/songs?id=0", body)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(SongsHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if !strings.Contains(rec.Body.String(), "[{\"lyricsPage\":1,\"text\":\"Verse 1\"},") {
		t.Errorf("data not found, got %s", rec.Body.String())
	}

	// Удаление песен
	req = httptest.NewRequest("DELETE", "/songs?id=0", body)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(SongsHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	// Проверка, что песня удалилась
	req = httptest.NewRequest("GET", "/library", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got: %d, error: %s", rec.Code, rec.Body)
	}
}

func TestAPIRequests(t *testing.T) {
	config.SetConfig()
	config.Database = models.MockDatabase{Data: make(map[int]models.SongData)}
	config.MusLib = services.NewLibrary()

	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		var l models.LyricsData
		l.Link = "example"
		l.ReleaseDate = "example"
		l.Text = "Verse 1\n\nVerse 2"

		resp, err := json.Marshal(l)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	})

	// Запуск сервера
	server := &http.Server{Addr: ":3001", Handler: nil}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(1 * time.Second)
		server.ListenAndServe()
		wg.Done()
	}()

	time.Sleep(1 * time.Second)

	// Добавление песни
	body := strings.NewReader(`{
	"group": "Muse", 
	"song": "Supermassive Black Hole"
}`)

	req := httptest.NewRequest("POST", "/songs", body)
	rec := httptest.NewRecorder()

	handler := http.HandlerFunc(SongsHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	// Выключение сервера
	server.Shutdown(context.Background())

	wg.Wait()

	// Проверка результата
	req = httptest.NewRequest("GET", "/library", nil)
	rec = httptest.NewRecorder()

	handler = http.HandlerFunc(LibraryHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got: %d, error: %s", rec.Code, rec.Body)
	}

	if !strings.Contains(rec.Body.String(), "Verse 1") {
		t.Errorf("data not found, got %s", rec.Body.String())
	}
}
