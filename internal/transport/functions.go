package transport

import (
	"sort"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

func paginateLibrary(data map[int]models.SongData) []models.PaginatedSongData {
	// Сортировка ключей в мапе
	keys := make([]int, 0, len(data))
	for k := range data {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Пагинация мапы
	var result []models.PaginatedSongData
	var currentPage int

	result = append(result, models.PaginatedSongData{CurrentPage: currentPage + 1}) // Нумерация страниц будет начинаться с единицы

	for k := range keys {
		v := data[k]

		if len(result[currentPage].Entries) == config.PageSize {
			currentPage++
			result = append(result, models.PaginatedSongData{CurrentPage: currentPage + 1})
		}

		var valueWithID models.SongDataWithID
		valueWithID.ID = k
		valueWithID.Group = v.Group
		valueWithID.Song = v.Song
		valueWithID.Lyrics = v.Lyrics
		result[currentPage].Entries = append(result[currentPage].Entries, valueWithID)
	}
	return result
}
