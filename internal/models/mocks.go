package models

import "context"

// Моковая база данных для использования в тестах.
type MockDatabase struct {
	Data map[int]SongData
}

func (db MockDatabase) AddSong(ctx context.Context, id int, song SongData) error {
	db.Data[id] = song

	return nil
}

func (db MockDatabase) RestoreData(ctx context.Context) (int, map[int]SongData, error) {
	var maxID int
	for k, _ := range db.Data {
		if k > maxID {
			maxID = k
		}
	}

	return maxID, db.Data, nil
}
