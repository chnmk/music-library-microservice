package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

// Пагинация текста песни по куплетам.
func paginateLyrics(lyrics string) []models.PaginatedLyrics {
	var result []models.PaginatedLyrics

	verses := strings.Split(lyrics, "\n\n")

	for idx, v := range verses {
		result = append(result, models.PaginatedLyrics{
			Page: idx + 1,
			Text: v,
		})
	}

	return result
}

// Запрос текста со стороннего API.
func requestLyrics(song models.NewSongData) (models.FullSongData, error) {
	var result models.FullSongData
	result.Artist = song.Artist
	result.Song = song.Song

	client := http.Client{Timeout: time.Duration(config.RequestTimeout) * time.Second}
	url := config.RequestServer + "/info?artist=" + song.Artist + "&song=" + song.Song
	url = strings.ReplaceAll(url, " ", "%20")

	slog.Debug(
		"requesting lyrics",
		"url", url,
	)

	resp, err := client.Get(url)
	if err != nil {
		return result, err
	}

	if resp.StatusCode != http.StatusOK {
		return result, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	var lyricsData models.LyricsData
	var buf bytes.Buffer

	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return result, err
	}

	if err = json.Unmarshal(buf.Bytes(), &lyricsData); err != nil {
		return result, err
	}

	result.Lyrics = lyricsData.Text
	result.ReleaseDate = lyricsData.ReleaseDate
	result.Link = lyricsData.Link

	return result, nil
}
