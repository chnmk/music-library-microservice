package services

import (
	"os"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

func TestPaginateLyrics(t *testing.T) {
	lyrics := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"

	paginated := paginateLyrics(lyrics)

	if len(paginated) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(paginated))
	}

	if paginated[0].Page != 1 || paginated[1].Page != 2 {
		t.Error("wrong page numeration")
	}

	if paginated[0].Text != "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?" {
		t.Error("wrong text for page 1")
	}

	if paginated[1].Text != "Ooh\nYou set my soul alight\nOoh\nYou set my soul alight" {
		t.Error("wrong text for page 2")
	}
}

func TestPaginateLyricsEmpty(t *testing.T) {
	lyricsEmpty := ""

	paginatedEmpty := paginateLyrics(lyricsEmpty)

	if len(paginatedEmpty) != 1 {
		t.Fatalf("expected 1 page, got %d", len(paginatedEmpty))
	}

	if paginatedEmpty[0].Page != 1 {
		t.Error("wrong page numeration")
	}

	if paginatedEmpty[0].Text != "" {
		t.Error("wrong text for page 1")
	}

	lyrics := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?"

	paginated := paginateLyrics(lyrics)

	if len(paginated) != 1 {
		t.Fatalf("expected 1 page, got %d", len(paginated))
	}

	if paginated[0].Page != 1 {
		t.Error("wrong page numeration")
	}

	if paginated[0].Text != "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?" {
		t.Error("wrong text for page 1")
	}

}

func TestCRUDFunctions(t *testing.T) {
	os.Setenv("REQUEST_SERVER", "http://localhost:3001")
	config.SetConfig()
	config.Database = &models.MockDatabase{}
	config.MusLib = NewLibrary()

	var song models.NewSongData

	song.Artist = "Artist 1"
	song.Song = "Song 1"

	config.MusLib.AddSong(song)

	song.Artist = "Artist 2"
	song.Song = "Song 2"

	config.MusLib.AddSong(song)

	// Получение всей библиотеки
	var requestSong models.FullSongData
	requestSong.Artist = "Artist 1"
	requestSong.Song = "Song 1"

	everything, err := config.MusLib.GetSongs(requestSong)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(everything) == 0 {
		t.Fatal("expected data")
	}
	if len(everything[0].Entries) == 0 {
		t.Fatal("expected data")
	}
	if everything[0].Entries[0].Artist != "Artist 1" {
		t.Errorf("expected artist 1, got %s", everything[0].Entries[0].Artist)
	}

	// Изменение песни
	var changedSong models.FullSongData

	changedSong.Artist = "Artist 111"
	changedSong.Song = "Song 111"
	changedSong.Lyrics = "Verse 1 \n\n Verse 2"

	err = config.MusLib.ChangeSong(song, changedSong)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	everything1, err := config.MusLib.GetSongs(changedSong)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(everything1) == 0 {
		t.Fatal("expected data")
	}
	if len(everything1[0].Entries) == 0 {
		t.Fatal("expected data")
	}
	if everything1[0].Entries[1].Artist != "Artist 111" {
		t.Errorf("data change failed, expected Artist 111, got %s", everything1[0].Entries[1].Artist)
	}

	// Удаление песни
	song.Artist = "Artist 1"
	song.Song = "Song 1"

	err = config.MusLib.DeleteSong(song)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}
