package transport

import (
	"strconv"
	"sync"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

var once sync.Once

func TestPaginateData(t *testing.T) {
	once.Do(func() {
		config.SetConfig()
	})

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
