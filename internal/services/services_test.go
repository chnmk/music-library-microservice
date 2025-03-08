package services

import (
	"strconv"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

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
