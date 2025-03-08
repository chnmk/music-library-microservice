package services

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

type musicLibrary struct {
	mu    sync.Mutex
	maxId int
	songs map[int]models.SongData
}

func (l *musicLibrary) AddSong(song models.SongData) {
	songWithLyrics, err := requestLyrics(song)
	if err != nil {
		slog.Info(
			"couldn't get lyrics",
			"err", err.Error(),
		)
	}

	id := l.maxId

	slog.Debug(
		"adding new song...",
		"id", id,
	)

	// Полагаю, песни стоит добавлять даже без текста
	l.mu.Lock()
	l.songs[id] = songWithLyrics
	l.maxId++
	if len(l.songs) > config.MaxEntries {
		l.songs = clearSongsData(l.songs)
	}
	l.mu.Unlock()

	err = config.Database.AddSong(context.Background(), id, songWithLyrics)
	if err != nil {
		slog.Info(
			"couldn't add song to database",
			"id", id,
			"err", err.Error(),
		)
	}

	slog.Debug(
		"song successfully added",
		"id", id,
	)
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
