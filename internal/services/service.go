package services

import (
	"errors"
	"sync"

	"github.com/chnmk/music-library-microservice/internal/models"
)

type musicLibrary struct {
	mu    sync.Mutex
	maxId int
	songs map[int]models.SongData
}

func (l *musicLibrary) AddSong(song models.SongData) error {
	var err error
	song.Lyrics, err = l.addSong(song)
	if err != nil {
		return err
	}

	l.mu.Lock()
	l.songs[l.maxId] = song
	l.maxId++
	l.mu.Unlock()

	return nil
}

func (l *musicLibrary) addSong(song models.SongData) (string, error) {
	/*
		TODO

		При добавлении сделать запрос в АПИ, описанного сваггером. Апи,
		описанный сваггером, будет поднят при проверке тестового задания.
		Реализовывать его отдельно не нужно
	*/
	return song.Song + "placeholder", nil
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

func (l *musicLibrary) GetLyrics(id int) (string, error) {
	l.mu.Lock()
	song, ok := l.songs[id]
	l.mu.Unlock()

	if !ok {
		return "", errors.New("song not found")
	}

	return song.Lyrics, nil
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
	return &musicLibrary{songs: make(map[int]models.SongData)}
}
