package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

type musicLibrary struct {
	mu    sync.Mutex
	maxId int
	songs map[int]models.SongData
}

func (l *musicLibrary) AddSong(song models.SongData) {
	songWithLyrics, err := l.requestLyrics(song)
	if err != nil {
		slog.Info(
			"couldn't get lyrics",
			"err", err.Error(),
		)
	}

	id := l.maxId

	// Полагаю, песни стоит добавлять даже без текста
	l.mu.Lock()
	l.songs[id] = songWithLyrics
	l.maxId++
	l.mu.Unlock()

	err = config.Database.AddSong(context.Background(), id, songWithLyrics)
	if err != nil {
		slog.Info(
			"couldn't add song to database",
			"id", id,
			"err", err.Error(),
		)
	}
}

func (l *musicLibrary) requestLyrics(song models.SongData) (models.SongData, error) {
	/*
		При добавлении сделать запрос в АПИ, описанного сваггером. Апи,
		описанный сваггером, будет поднят при проверке тестового задания.
		Реализовывать его отдельно не нужно
	*/

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

func (l *musicLibrary) GetSongs(params map[string]string) ([]models.PaginatedSongData, error) {
	l.mu.Lock()
	if len(l.songs) == 0 {
		return nil, errors.New("no songs found")
	}
	l.mu.Unlock()

	filtered, err := filter(l.songs, params)
	if err != nil {
		return nil, errors.New("no songs found")
	}

	return paginateLibrary(filtered), nil
}

func (l *musicLibrary) GetLyrics(id int) ([]models.PaginatedLyrics, error) {
	l.mu.Lock()
	song, ok := l.songs[id]
	l.mu.Unlock()

	if !ok {
		return nil, errors.New("song not found")
	}

	paginatedLyrics := paginateLyrics(song.Lyrics)

	return paginatedLyrics, nil
}

func (l *musicLibrary) ChangeSong(id int, song models.SongData) error {
	l.mu.Lock()
	_, ok := l.songs[id]
	if !ok {
		l.mu.Unlock()
		return errors.New("song not found")
	}

	l.songs[id] = song
	l.mu.Unlock()

	return nil
}

func (l *musicLibrary) DeleteSong(id int) error {
	l.mu.Lock()
	_, ok := l.songs[id]
	if !ok {
		l.mu.Unlock()
		return errors.New("song not found")
	}

	delete(l.songs, id)
	l.mu.Unlock()

	return nil
}

func NewLibrary() models.MusicLibrary {
	if config.RestoreData {
		max, data, err := config.Database.RestoreData(context.Background())
		if err != nil {
			slog.Info(
				"failed to restore data from DB",
				"err", err,
			)
		} else {
			return &musicLibrary{maxId: max, songs: data}
		}
	}

	return &musicLibrary{songs: make(map[int]models.SongData)}
}
