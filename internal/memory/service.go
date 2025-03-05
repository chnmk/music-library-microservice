package memory

import "github.com/chnmk/music-library-microservice/internal/models"

type musicLibrary struct {
	songs map[string]models.SongData
}

func (l *musicLibrary) AddSong(song models.SongData) error {
	// TODO
	return nil
}

func (l *musicLibrary) addSong() {
	/*
		TODO

		При добавлении сделать запрос в АПИ, описанного сваггером. Апи,
		описанный сваггером, будет поднят при проверке тестового задания.
		Реализовывать его отдельно не нужно
	*/
}

func (l *musicLibrary) GetSongs() (map[string]models.SongData, error) {
	// TODO
	return nil, nil
}

func (l *musicLibrary) GetLyrics(id string) (string, error) {
	// TODO
	return "", nil
}

func (l *musicLibrary) ChangeSong(id string, song models.SongData) error {
	// TODO
	return nil
}

func (l *musicLibrary) DeleteSong(id string) error {
	// TODO
	return nil
}

func NewLibrary() models.MusicLibrary {
	return &musicLibrary{songs: make(map[string]models.SongData)}
}
