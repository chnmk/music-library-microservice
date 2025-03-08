package services

import (
	"os"
	"strconv"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

func TestCRUDFunctions(t *testing.T) {
	os.Setenv("REQUEST_SERVER", "http://localhost:3001")
	config.SetConfig()
	config.Database = models.MockDatabase{Data: make(map[int]models.SongData)}
	lib := NewLibrary()

	var song models.SongData

	song.Group = "Group 1"
	song.Song = "Song 1"
	song.Lyrics = "Verse 1\n\nVerse 2"

	lib.AddSong(song)

	// Получение всей библиотеки
	everything, err := lib.GetSongs(make(map[string]string))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(everything) == 0 {
		t.Fatal("expected data")
	}
	if len(everything[0].Entries) == 0 {
		t.Fatal("expected data")
	}

	id := everything[0].Entries[0].ID
	if id != 0 {
		t.Fatalf("expected id 0, got %d", id)
	}

	// Получение текста
	lyr, err := lib.GetLyrics(id)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(lyr) != 2 {
		t.Fatalf("expected two verses, got: %v", len(lyr))
	}

	if lyr[0].Text != "Verse 1" {
		t.Fatalf("expected Verse 1, got: %s", lyr[0].Text)
	}

	// Изменение песни
	var song1 models.SongData

	song1.Group = "Group 111"
	song1.Song = "Song 111"
	song1.Lyrics = "Verse 1 \n\n Verse 2"

	err = lib.ChangeSong(id, song1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	everything1, err := lib.GetSongs(make(map[string]string))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(everything1) == 0 {
		t.Fatal("expected data")
	}
	if len(everything1[0].Entries) == 0 {
		t.Fatal("expected data")
	}

	gr := everything1[0].Entries[0].Group
	if gr != "Group 111" {
		t.Errorf("data change failed, expected Group 111, got %s", gr)
	}

	// Удаление песни
	err = lib.DeleteSong(id)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	_, err = lib.GetSongs(make(map[string]string))
	if err == nil {
		t.Fatalf("expected error")
	}

}

func TestFilterDataByGroup(t *testing.T) {
	data := make(map[int]models.SongData)
	params := make(map[string]string)

	var exampleSong models.SongData
	for i := range 23 {
		exampleSong.Group = "Group " + strconv.Itoa(i)
		exampleSong.Song = "Song " + strconv.Itoa(i)
		exampleSong.Lyrics = "Lyrics " + strconv.Itoa(i)
		data[i] = exampleSong
	}

	params["group"] = "18"

	filtered, err := filter(data, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(filtered) > 1 {
		t.Errorf("expected one entry, got: %d", len(filtered))
	}

	for k, v := range filtered {
		if k != 18 || v.Group != "Group 18" {
			t.Errorf("expected: key 18, value Group 18. got: key %d, value %s", k, v.Group)
		}
	}
}

func TestFilterDataBySongAndLyrics(t *testing.T) {
	data := make(map[int]models.SongData)
	params := make(map[string]string)

	var exampleSong models.SongData
	for i := range 23 {
		exampleSong.Group = "Group " + strconv.Itoa(i)
		exampleSong.Song = "Song " + strconv.Itoa(i)
		exampleSong.Lyrics = "Lyrics " + strconv.Itoa(i)
		data[i] = exampleSong
	}

	params["song"] = "19"
	params["lyrics"] = "19"

	filtered, err := filter(data, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(filtered) > 1 {
		t.Errorf("expected one entry, got: %d", len(filtered))
	}

	for k, v := range filtered {
		if k != 19 || v.Song != "Song 19" {
			t.Errorf("expected: key 18, value Song 18. got: key %d, value %s", k, v.Song)
		}
	}
}

func TestFilterDataMultiple(t *testing.T) {
	data := make(map[int]models.SongData)
	params := make(map[string]string)

	var exampleSong models.SongData
	for i := range 22 {
		exampleSong.Group = "Group " + strconv.Itoa(i)
		exampleSong.Song = "Song " + strconv.Itoa(i)
		exampleSong.Lyrics = "Lyrics " + strconv.Itoa(i)
		data[i] = exampleSong
	}

	params["song"] = "3"

	filtered, err := filter(data, params)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if len(filtered) != 2 { // 3, 13
		t.Errorf("expected two entries, got: %d", len(filtered))
	}

	if _, ok := filtered[3]; !ok {
		t.Errorf("expected data for key 3")
	}

	if _, ok := filtered[13]; !ok {
		t.Errorf("expected data for key 13")
	}

	if filtered[3].Group != "Group 3" || filtered[13].Group != "Group 13" {
		t.Errorf("expected Group 3 and Group 13, got %s and %s", filtered[3].Group, filtered[13].Group)
	}
}

func TestFilterDataInvalid(t *testing.T) {
	data := make(map[int]models.SongData)
	params := make(map[string]string)

	var exampleSong models.SongData
	for i := range 22 {
		exampleSong.Group = "Group " + strconv.Itoa(i)
		exampleSong.Song = "Song " + strconv.Itoa(i)
		exampleSong.Lyrics = "Lyrics " + strconv.Itoa(i)
		data[i] = exampleSong
	}

	filtered, err := filter(data, params)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	if len(filtered) != len(data) {
		t.Errorf("expected filtered data to be the same length")
	}

	params["song"] = "334134"

	_, err = filter(data, params)
	if err == nil {
		t.Errorf("expected no data")
	}
}

func TestPaginateData(t *testing.T) {
	config.SetConfig()

	data := make(map[int]models.SongData)

	var exampleSong models.SongData
	for i := range 23 {
		exampleSong.Group = "Group " + strconv.Itoa(i)
		exampleSong.Song = "Song " + strconv.Itoa(i)
		exampleSong.Lyrics = "Lyrics " + strconv.Itoa(i)
		data[i] = exampleSong
	}

	paginated := paginateLibrary(data)

	if len(paginated) != 3 {
		t.Errorf("got %d pages, expected 3", len(paginated))
	}

	if len(paginated[0].Entries) != config.PageSize || len(paginated[1].Entries) != config.PageSize {
		t.Error("wrong page length")
	}

	if paginated[0].CurrentPage != 1 || paginated[1].CurrentPage != 2 || paginated[2].CurrentPage != 3 {
		t.Error("wrong page numeration")
	}

	if paginated[0].Entries[0].Group != "Group 0" {
		t.Log(data)

		t.Log(paginated)
		t.Errorf("wrong data: %s, expected 'Group 0'", paginated[0].Entries[0].Group)
	}

	if paginated[1].Entries[0].Group != "Group 10" {
		t.Errorf("wrong data: %s, expected 'Group 10'", paginated[1].Entries[0].Group)
	}
}

func TestPaginateEmptyData(t *testing.T) {
	data := make(map[int]models.SongData)

	empty := paginateLibrary(data)
	if len(empty[0].Entries) != 0 {
		t.Error("expected empty result")
	}
}

func TestPaginateLyrics(t *testing.T) {
	lyrics := "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"

	paginated := paginateLyrics(lyrics)

	if len(paginated) != 2 {
		t.Fatalf("expected 2 pages, got %d", len(paginated))
	}

	if paginated[0].LyricsPage != 1 || paginated[1].LyricsPage != 2 {
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

	if paginatedEmpty[0].LyricsPage != 1 {
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

	if paginated[0].LyricsPage != 1 {
		t.Error("wrong page numeration")
	}

	if paginated[0].Text != "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?" {
		t.Error("wrong text for page 1")
	}

}

func TestClearData(t *testing.T) {
	os.Setenv("MAX_ENTRIES", "3")
	config.SetConfig()
	config.Database = models.MockDatabase{Data: make(map[int]models.SongData)}
	lib := NewLibrary()

	var song models.SongData

	song.Group = "Group 1"
	song.Song = "Song 1"
	song.Lyrics = "Verse 1\n\nVerse 2"

	lib.AddSong(song)

	var song2 models.SongData

	song2.Group = "Group 2"
	song2.Song = "Song 2"
	song2.Lyrics = "Verse 1\n\nVerse 2"

	lib.AddSong(song2)
	lib.AddSong(song2)
	lib.AddSong(song2)

	// Получение всей библиотеки
	everything, err := lib.GetSongs(make(map[string]string))
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(everything[0].Entries) != 3 {
		t.Errorf("data not cleared, current length: %d", len(everything[0].Entries))
	}
	if everything[0].Entries[0].Group == "Song 1" || everything[0].Entries[1].Group == "Song 1" || everything[0].Entries[2].Group == "Song 1" {
		t.Errorf("deleted wrong entry")
	}
}
