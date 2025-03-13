package transport

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

// Получение данных библиотеки с фильтрацией по всем полям и пагинацией.
//
//	@Summary		Get music library
//	@Description	get paginated song data with optional flitration by artist, song, lyrics or page
//	@Produce		json
//	@Param			artist	query	string	false	"fliter by artist"
//	@Param			song	query	string	false	"fliter by song"
//	@Param			lyrics	query	string	false	"fliter by lyrics"
//	@Param			page	query	string	false	"fliter by page"
//	@Success		200		{array}	models.PaginatedSongData
//	@Failure		404
//	@Failure		500
//	@Router			/library [get]
func libraryGet(w http.ResponseWriter, r *http.Request) {
	var params models.FullSongData
	var page string

	// Получение параметров запроса.
	params.Artist = r.URL.Query().Get("artist")
	params.Song = r.URL.Query().Get("song")
	params.Lyrics = r.URL.Query().Get("lyrics")
	params.ReleaseDate = r.URL.Query().Get("releasedate")
	params.Link = r.URL.Query().Get("Link")
	page = r.URL.Query().Get("page")

	// Получение всех песен.
	lib, err := config.MusLib.GetSongs(params)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	// Получение конкретной страницы при необходимости.
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil {
			slog.Error("request failed", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if pageInt > len(lib)+1 {
			slog.Error("request failed", "err", "page not found")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("page not found"))
			return
		}

		// Запись ответа с конкретной страницей.
		resp, err := json.Marshal([]models.PaginatedSongData{lib[pageInt]})
		if err != nil {
			slog.Error("request failed", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

	// Запись ответа.
	resp, err := json.Marshal(lib)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

// Добавление новой песни.
//
//	@Summary		Add song
//	@Description	add song to the server
//	@Accept			json
//	@Param			song	body	models.SongData	true	"artist and song title, lyrics are optional"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/songs [post]
func songsPost(w http.ResponseWriter, r *http.Request) {
	var song models.NewSongData
	var buf bytes.Buffer

	// Получение данных из тела запроса.
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if song.Artist == "" || song.Song == "" {
		slog.Error("request failed", "err", "song data invalid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("song data invalid"))
		return
	}

	// Добавление песни.
	err = config.MusLib.AddSong(song)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Получение текста песни с пагинацией по куплетам.
//
//	@Summary		Get lyrics
//	@Description	get paginated lyrics data for a song
//	@Accept			json
//	@Param			id	query	string	true	"id of a song"
//	@Success		200	{array}	models.PaginatedLyrics
//	@Failure		400
//	@Failure		500
//	@Router			/songs [get]
func songsGetLyrics(w http.ResponseWriter, r *http.Request) {
	var song models.NewSongData

	// Получение параметра id.
	song.Artist = r.URL.Query().Get("artist")
	song.Song = r.URL.Query().Get("song")

	// Получение пагинированного текста песни.
	result, err := config.MusLib.GetLyrics(song)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	// Запись ответа.
	resp, err := json.Marshal(result)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

// Изменение данных песни.
//
//	@Summary		Edit song
//	@Description	edit song data
//	@Accept			json
//	@Param			id		query	string			true	"id of a song"
//	@Param			song	body	models.SongData	true	"artist and song title, lyrics are optional"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/songs [put]
func songsPut(w http.ResponseWriter, r *http.Request) {
	var oldSong models.NewSongData

	// Получение параметра id.
	oldSong.Artist = r.URL.Query().Get("artist")
	oldSong.Song = r.URL.Query().Get("song")

	var song models.FullSongData
	var buf bytes.Buffer

	// Получение данных из тела запроса.
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &song); err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if oldSong.Artist == "" || oldSong.Song == "" {
		slog.Error("request failed", "err", "song data invalid")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("song data invalid"))
		return
	}

	// Редактирование песни.
	err = config.MusLib.ChangeSong(oldSong, song)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}

// Удаление песни.
//
//	@Summary		Delete song
//	@Description	delete song data
//	@Accept			json
//	@Param			id	query	string	true	"id of a song"
//	@Success		200
//	@Failure		400
//	@Failure		500
//	@Router			/songs [delete]
func songsDelete(w http.ResponseWriter, r *http.Request) {
	var song models.NewSongData

	// Получение параметра id.
	song.Artist = r.URL.Query().Get("artist")
	song.Song = r.URL.Query().Get("song")

	// Удаление песни.
	err := config.MusLib.DeleteSong(song)
	if err != nil {
		slog.Error("request failed", "err", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("success"))
}
