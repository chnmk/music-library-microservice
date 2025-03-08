package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

/*
	TODO: ошибки, коды, ответы, пагинация. Вынести повторяющиеся части кода.
*/

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
func libraryGet(w http.ResponseWriter, r *http.Request) {
	params := make(map[string]string)

	params["group"] = r.URL.Query().Get("group")
	params["song"] = r.URL.Query().Get("song")
	params["lyrics"] = r.URL.Query().Get("lyrics")

	page := r.URL.Query().Get("page")

	lib, err := config.MusLib.GetSongs(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if pageInt > len(lib)+1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("page not found"))
			return
		}

		// TODO: return result[page]
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("placeholder: result json"))
	}

	// TODO: return result
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: paginated result json WITH SONG ID"))
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

	if song.Group == "" || song.Song == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("song data invalid"))
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
	id_string := r.URL.Query().Get("id")
	if id_string == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := config.MusLib.GetLyrics(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	_ = result

	// TODO: paginate string

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: paginated result json"))
}

func songsPut(w http.ResponseWriter, r *http.Request) {
	id_string := r.URL.Query().Get("id")
	if id_string == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var song models.SongData
	var buf bytes.Buffer

	_, err = buf.ReadFrom(r.Body)
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

	if song.Group == "" || song.Song == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("song data invalid"))
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
	id_string := r.URL.Query().Get("id")
	if id_string == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("expected 'id' parameter"))
		return
	}

	id, err := strconv.Atoi(id_string)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	err = config.MusLib.DeleteSong(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder: success message"))
}
