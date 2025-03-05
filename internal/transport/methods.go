package transport

import (
	"net/http"

	"github.com/chnmk/music-library-microservice/internal/config"
)

func libraryGet(w http.ResponseWriter, r *http.Request) {
	// TODO

	config.MusLib.GetSongs()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsPost(w http.ResponseWriter, r *http.Request) {
	// TODO

	/*
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
	*/

	config.MusLib.AddSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsGetLyrics(w http.ResponseWriter, r *http.Request) {
	// TODO

	config.MusLib.GetLyrics()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsPut(w http.ResponseWriter, r *http.Request) {
	// TODO

	config.MusLib.ChangeSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsDelete(w http.ResponseWriter, r *http.Request) {
	// TODO

	config.MusLib.DeleteSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}
