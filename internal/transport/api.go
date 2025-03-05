package transport

import (
	"net/http"
)

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	if r.Method == http.MethodGet {
		libraryGet(w, r)
	}
}

func SongsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Добавление новой песни в формате
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

func paginate() {
	// TODO
}

func validate() {
	// TODO
}
