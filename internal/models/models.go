package models

import (
	"context"
)

type Database interface {
	Insert(ctx context.Context, song FullSongData) error
	SelectAll(ctx context.Context, params FullSongData) ([]PaginatedSongData, error)
	SelectLyrics(ctx context.Context, song NewSongData) (string, error)
	Update(ctx context.Context, song NewSongData, params FullSongData) error
	Delete(ctx context.Context, song NewSongData) error
}

type MusicLibrary interface {
	AddSong(song NewSongData) error
	GetSongs(params FullSongData) (songs []PaginatedSongData, err error)
	GetLyrics(song NewSongData) (lyrics []PaginatedLyrics, err error)
	ChangeSong(song NewSongData, params FullSongData) error
	DeleteSong(song NewSongData) error
}

// Данные, которые приходят в запросах.
type NewSongData struct {
	Artist string `json:"artist"`
	Song   string `json:"song"`
}

// Данные, которые возвращаются клиенту, либо параметры запроса.
type FullSongData struct {
	Artist      string `json:"artist"`
	Song        string `json:"song"`
	Lyrics      string `json:"lyrics"`
	ReleaseDate string `json:"releasedate"`
	Link        string `json:"link"`
}

// Пагинированные данные, которые возвращаются клиенту при GET-запросе к /library.
type PaginatedSongData struct {
	Page    int            `json:"page"`
	Entries []FullSongData `json:"entries"`
}

// Данные, которые приходят со строннего API.
type LyricsData struct {
	ReleaseDate string `json:"releasedate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Данные, которые возвращаются клиенту при GET-запросе к /songs.
type PaginatedLyrics struct {
	Page int    `json:"page"`
	Text string `json:"text"`
}
