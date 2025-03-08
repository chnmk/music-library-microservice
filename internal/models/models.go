package models

import (
	"context"
)

type Database interface {
	AddSong(ctx context.Context, id int, song SongData) error
	RestoreData(ctx context.Context) (id int, songs map[int]SongData, err error)
}

type MusicLibrary interface {
	AddSong(song SongData)
	GetSongs(params map[string]string) (songs []PaginatedSongData, err error)
	GetLyrics(id int) (lyrics []PaginatedLyrics, err error)
	ChangeSong(id int, song SongData) (err error)
	DeleteSong(id int) (err error)
}

// Данные, которые приходят в запросах. Поле lyrics может быть пустым.
type SongData struct {
	Group  string `json:"group"`
	Song   string `json:"song"`
	Lyrics string `json:"lyrics"`
}

// Структура данных, которые возвращаются клиенту при GET-запросе к /library.
type SongDataWithID struct {
	ID     int    `json:"id"`
	Group  string `json:"group"`
	Song   string `json:"song"`
	Lyrics string `json:"lyrics"`
}

// Пагинированные данные, которые возвращаются клиенту при GET-запросе к /library.
type PaginatedSongData struct {
	CurrentPage int              `json:"currentPage"`
	Entries     []SongDataWithID `json:"entries"`
}

// Данные, которые приходят со строннего API.
type LyricsData struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Данные, которые возвращаются клиенту при GET-запросе к /songs.
type PaginatedLyrics struct {
	LyricsPage int    `json:"lyricsPage"`
	Text       string `json:"text"`
}
