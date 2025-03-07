package transport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

/*
	TODO: ошибки, коды, ответы, пагинация. Вынести повторяющиеся части кода.
*/

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
	page := r.URL.Query().Get("page")

	var filtered map[int]models.SongData

	if group != "" || song != "" || lyrics != "" {
		filtered = make(map[int]models.SongData)

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

			filtered[k] = v
		}
	}

	// TODO: paginate map
	var result []models.PaginatedSongData
	var currentPage int

	result = append(result, models.PaginatedSongData{CurrentPage: 1}) // Нумерация страниц будет начинаться с единицы

	for k, v := range filtered {
		if len(result[currentPage].Entries) > config.PageSize {
			currentPage++
			result[currentPage].CurrentPage = currentPage + 1
		}

		var valueWithID models.SongDataWithID
		valueWithID.ID = k
		valueWithID.Group = v.Group
		valueWithID.Song = v.Song
		valueWithID.Lyrics = v.Lyrics
		result[currentPage].Entries = append(result[currentPage].Entries, valueWithID)
	}

	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if pageInt > len(result)+1 {
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
