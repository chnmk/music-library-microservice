package transport

import (
	"net/http"

	"github.com/chnmk/music-library-microservice/internal/services"
)

func libraryGet(w http.ResponseWriter, r *http.Request) {
	// TODO

	services.MusLib.GetSongs()

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

	services.MusLib.AddSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsGetLyrics(w http.ResponseWriter, r *http.Request) {
	// TODO

	services.MusLib.GetLyrics()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsPut(w http.ResponseWriter, r *http.Request) {
	// TODO

	services.MusLib.ChangeSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}

func songsDelete(w http.ResponseWriter, r *http.Request) {
	// TODO

	services.MusLib.DeleteSong()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("placeholder"))
}
