package database

import (
	"context"
	"testing"

	"github.com/chnmk/music-library-microservice/internal/models"
)

func TestDatabaseExample(t *testing.T) {
	mydb := models.MockDatabase{Data: make([]models.FullSongData, 0)}

	var song1 models.FullSongData
	var song2 models.FullSongData

	song1.Artist = "artist1"
	song2.Artist = "artist2"

	mydb.Insert(context.Background(), song1)
	mydb.Insert(context.Background(), song2)

	var params models.FullSongData

	params.Artist = "artist1"

	newData, err := mydb.SelectAll(context.Background(), params)
	if err != nil {
		t.Fatal("expected no error")
	}
	if len(newData) != 1 {
		t.Error("expected one entry")
	}
}
