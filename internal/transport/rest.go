package transport

import "net/http"

// TODO: подумать ещё...

func LibraryHandler(w http.ResponseWriter, r *http.Request) {
	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	if r.Method == http.MethodGet {
		libraryGet(w, r)
	}
}

func SongsHandler(w http.ResponseWriter, r *http.Request) {
	// Добавление новой песни в формате
	if r.Method == http.MethodPost {
		songsPost(w, r)
	}

	// Получение текста песни с пагинацией по куплетам
	if r.Method == http.MethodGet {
		songsGet(w, r)
	}

	// Изменение данных песни
	if r.Method == http.MethodPut {
		songsPut(w, r)
	}

	// Удаление песни
	if r.Method == http.MethodDelete {
		songsDelete(w, r)
	}
}

func libraryGet(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func songsPost(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func songsGet(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func songsPut(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func songsDelete(w http.ResponseWriter, r *http.Request) {
	// TODO
}
