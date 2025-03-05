package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func libraryGet(w http.ResponseWriter, r *http.Request) {
	lib, err := config.MusLib.GetSongs()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	lyrics := r.URL.Query().Get("lyrics")

	if group != "" || song != "" || lyrics != "" {
		result := make(map[string]models.SongData)

		for k, v := range lib {
			if group != "" && !strings.Contains(v.Group, group) {
				break
			}

			if song != "" && !strings.Contains(v.Song, group) {
				break
			}

			if lyrics != "" && !strings.Contains(v.Lyrics, group) {
				break
			}

			result[k] = v
		}
	}

	// TODO
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: paginated result json"))
}

// Добавление новой песни
func songsPost(w http.ResponseWriter, r *http.Request) {
	var song models.SongData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = config.MusLib.AddSong(song)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: success message"))
}

// Получение текста песни с пагинацией по куплетам
func songsGetLyrics(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	result, err := config.MusLib.GetLyrics(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_ = result

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: paginated result json"))
}

func songsPut(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	var song models.SongData
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = config.MusLib.ChangeSong(id, song)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: success message"))
}

func songsDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
	id := r.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	err := config.MusLib.DeleteSong(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: success message"))
}
