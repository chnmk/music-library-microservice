package services

import "github.com/chnmk/music-library-microservice/internal/models"

type MusicLibrary interface {
	// TODO
	AddSong()
	addSong()

	GetSongs()
	GetLyrics()

	ChangeSong()
	DeleteSong()
}

type musicLibrary struct {
	// Дополнительные требования не указаны, пока пусть будет срез
	// TODO: на самом деле менять и удалять можно будет только по ID, так что...
	songs []models.SongData
}

var MusLib MusicLibrary

func (l *musicLibrary) AddSong() {
	// TODO
}

func (l *musicLibrary) addSong() {
	/*
		TODO

		При добавлении сделать запрос в АПИ, описанного сваггером. Апи,
		описанный сваггером, будет поднят при проверке тестового задания.
		Реализовывать его отдельно не нужно
	*/
}

func (l *musicLibrary) GetSongs() {
	// TODO
}

func (l *musicLibrary) GetLyrics() {
	// TODO
}

func (l *musicLibrary) ChangeSong() {
	// TODO
}

func (l *musicLibrary) DeleteSong() {
	// TODO
}

func NewLibrary() MusicLibrary {
	return &musicLibrary{}
}
