package transport

import "net/http"

// Возможно стоит сделать два хендлера, на /lib и /song
func LibHandler(w http.ResponseWriter, r *http.Request) {

	// Получение данных библиотеки с фильтрацией по всем полям и пагинацией
	if r.Method == "" {

	}

	// Получение текста песни с пагинацией по куплетам
	if r.Method == "" {

	}

	// Удаление песни
	if r.Method == "" {

	}

	// Изменение данных песни
	if r.Method == "" {

	}

	// Добавление новой песни в формате
	if r.Method == "" {

	}
}

/*
JSON

{
 "group": "Muse",
 "song": "Supermassive Black Hole"
}
*/
