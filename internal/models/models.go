package models

type MusicLibrary interface {
	// TODO
	AddSong()

	GetSongs()
	GetLyrics()

	ChangeSong()
	DeleteSong()
}

type SongData struct {
	Group string
	Song  string
}
