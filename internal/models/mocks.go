package models

import (
	"context"
	"errors"
	"fmt"
)

// Пример моковой базы данных для использования в тестах.
type MockDatabase struct {
	Data []FullSongData
}

func (db *MockDatabase) Insert(ctx context.Context, song FullSongData) error {
	db.Data = append(db.Data, song)

	return nil
}

func (db MockDatabase) SelectAll(ctx context.Context, params FullSongData) ([]PaginatedSongData, error) {
	var result []PaginatedSongData
	var currentPage int

	result = append(result, PaginatedSongData{Page: currentPage + 1})

	for _, v := range db.Data {
		if len(result[currentPage].Entries) >= 10 {
			currentPage++
			result = append(result, PaginatedSongData{Page: currentPage + 1})
		}

		if params.Artist != "" && params.Song == "" {
			if params.Artist == v.Artist && params.Song == v.Song {
				result[currentPage].Entries = append(result[currentPage].Entries, v)
			}
		} else {
			result[currentPage].Entries = append(result[currentPage].Entries, v)
		}
	}

	if len(result[0].Entries) == 0 {
		fmt.Println(result)
		return nil, errors.New("no data")
	}

	return result, nil
}

func (db MockDatabase) SelectLyrics(ctx context.Context, song NewSongData) (string, error) {
	for _, v := range db.Data {
		if v.Artist == song.Artist && v.Song == song.Song {
			return v.Lyrics, nil
		}
	}

	return "", errors.New("something is wrong")
}

func (db *MockDatabase) Update(ctx context.Context, song NewSongData, params FullSongData) error {
	var exists bool

	for i, v := range db.Data {
		if v.Artist == song.Artist && v.Song == song.Song {
			exists = true

			if params.Artist != "" {
				db.Data[i].Artist = params.Artist
			}
			if params.Song != "" {
				db.Data[i].Song = params.Song
			}
			if params.Lyrics != "" {
				db.Data[i].Lyrics = params.Lyrics
			}
			if params.ReleaseDate != "" {
				db.Data[i].ReleaseDate = params.ReleaseDate
			}
			if params.Link != "" {
				db.Data[i].Link = params.Link
			}
		}
	}

	if !exists {
		return errors.New("something is wrong")
	}

	return nil
}

func (db *MockDatabase) Delete(ctx context.Context, song NewSongData) error {
	id := -1

	for i, v := range db.Data {
		if v.Artist == song.Artist && v.Song == song.Song {
			id = i
		}
	}

	if id == -1 {
		return errors.New("something is wrong")
	}

	db.Data = append(db.Data[:id], db.Data[id+1:]...)

	return nil
}
