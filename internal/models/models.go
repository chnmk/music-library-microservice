package models

type MusicLibrary interface {
	// TODO
	AddSong(song SongData) (err error)
	GetSongs() (allSongs map[int]SongData, err error)
	GetLyrics(id int) (lyrics string, err error)
	ChangeSong(id int, song SongData) (err error)
	DeleteSong(id int) (err error)
}

type SongData struct {
	Group  string
	Song   string
	Lyrics string
}
