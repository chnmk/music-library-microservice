package services

import (
	"errors"
	"sort"
	"strings"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

func filter(lib map[int]models.SongData, params map[string]string) (map[int]models.SongData, error) {
	if len(params) > 0 {
		filtered := make(map[int]models.SongData)

		for k, v := range lib {
			if _, ok := params["group"]; ok && !strings.Contains(v.Group, params["group"]) {
				continue
			}

			if _, ok := params["song"]; ok && !strings.Contains(v.Song, params["song"]) {
				continue
			}

			if _, ok := params["lyrics"]; ok && !strings.Contains(v.Lyrics, params["lyrics"]) {
				continue
			}

			filtered[k] = v
		}
		if len(filtered) == 0 {
			return nil, errors.New("no songs found")
		}

		return filtered, nil
	}

	return lib, nil
}

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
