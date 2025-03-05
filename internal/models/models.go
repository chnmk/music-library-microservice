package models

type MusicLibrary interface {
	// TODO
	AddSong(song SongData) (err error)
	GetSongs() (allSongs map[string]SongData, err error)
	GetLyrics(id string) (lyrics string, err error)
	ChangeSong(id string, song SongData) (err error)
	DeleteSong(id string) (err error)
}

type SongData struct {
	Group  string
	Song   string
	Lyrics string
}
