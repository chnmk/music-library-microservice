package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

// Фильтрация данных по группе, названию песни или тексту.
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

// Пагинация всех данных на сервере.
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

	for _, k := range keys {
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

// Пагинация текста песни по куплетам.
func paginateLyrics(lyrics string) []models.PaginatedLyrics {
	var result []models.PaginatedLyrics

	verses := strings.Split(lyrics, "\n\n")

	for idx, v := range verses {
		result = append(result, models.PaginatedLyrics{
			LyricsPage: idx + 1,
			Text:       v,
		})
	}

	return result
}

// Запрос текста со стороннего API.
func requestLyrics(song models.SongData) (models.SongData, error) {
	client := http.Client{Timeout: time.Duration(config.RequestTimeout) * time.Second}
	url := config.RequestServer + "/info?group=" + song.Group + "&song=" + song.Song
	url = strings.ReplaceAll(url, " ", "%20")

	slog.Debug(
		"requesting lyrics",
		"url", url,
	)

	resp, err := client.Get(url)
	if err != nil {
		return song, err
	}

	if resp.StatusCode != http.StatusOK {
		return song, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var lyricsData models.LyricsData
	var buf bytes.Buffer

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return song, err
	}

	if err = json.Unmarshal(buf.Bytes(), &lyricsData); err != nil {
		return song, err
	}

	song.Lyrics = lyricsData.Text

	return song, nil
}

// Очистка данных из памяти при превышении MAX_ENTRIES. Оставляет 75% самых новых записей.
func clearSongsData(songs map[int]models.SongData) map[int]models.SongData {
	result := make(map[int]models.SongData)

	slog.Info(
		"max entires limit reached, clearing data...",
		"entries", len(songs),
	)

	// Сортировка ключей в мапе.
	keys := make([]int, 0, len(songs))
	for k := range songs {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// Оставляем 75% самых новых записей.
	max := int(math.Floor(float64(len(keys)) * 0.25))
	newKeys := keys[max:]

	for _, k := range newKeys {
		result[k] = songs[k]
	}

	slog.Info(
		"data clear complete",
		"entries", len(result),
	)

	return result
}
