package models

import (
	"context"
)

type Database interface {
	// TODO Дополнительно: возможность восстанавливать из БД при запуске (.env)
	AddSong(ctx context.Context, id int, song SongData) error
}

type MusicLibrary interface {
	AddSong(song SongData)
	GetSongs(params map[string]string) (songs []PaginatedSongData, err error)
	GetLyrics(id int) (lyrics []PaginatedLyrics, err error)
	ChangeSong(id int, song SongData) (err error)
	DeleteSong(id int) (err error)
}

type SongData struct {
	Group  string
	Song   string
	Lyrics string
}

type SongDataWithID struct {
	ID     int    `json:"id"`
	Group  string `json:"group"`
	Song   string `json:"song"`
	Lyrics string `json:"lyrics"`
}

type PaginatedSongData struct {
	CurrentPage int              `json:"currentPage"`
	Entries     []SongDataWithID `json:"entries"`
}

type LyricsData struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type PaginatedLyrics struct {
	LyricsPage int    `json:"lyricsPage"`
	Text       string `json:"text"`
}
