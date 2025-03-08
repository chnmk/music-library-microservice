package database

import (
	"context"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/models"
)

func TestDatabaseExample(t *testing.T) {
	mydb := models.MockDatabase{Data: make(map[int]models.SongData)}

	var song1 models.SongData
	var song2 models.SongData

	song1.Group = "group1"
	song2.Group = "group2"

	mydb.AddSong(context.Background(), 1, song1)
	mydb.AddSong(context.Background(), 2, song2)

	maxID, newData, _ := mydb.RestoreData(context.Background())
	if maxID != 2 || newData[1].Group != "group1" || newData[2].Group != "group2" {
		t.Fatalf("expected the same data")
	}
}
