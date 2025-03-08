package models

import (
	"context"
)

type Database interface {
	// TODO
	// Создать БД (миграции)
	// Положить в БД
	// Дополнительно: возможность восстанавливать из БД при запуске (.env)
	DeleteMe(ctx context.Context)
}

type MusicLibrary interface {
	// TODO
	AddSong(song SongData) (err error)
	GetSongs(params map[string]string) (songs []PaginatedSongData, err error)
	GetLyrics(id int) (lyrics string, err error)
	ChangeSong(id int, song SongData) (err error)
	DeleteSong(id int) (err error)
}

type SongData struct {
	Group  string
	Song   string
	Lyrics string
}

type SongDataWithID struct {
	ID     int
	Group  string
	Song   string
	Lyrics string
}

type PaginatedSongData struct {
	CurrentPage int
	Entries     []SongDataWithID
}
