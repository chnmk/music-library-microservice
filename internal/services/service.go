package services

import (
	"log/slog"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
)

// Имплементация интерфейса models.MusicLibrary.
type musicLibrary struct {
}

// Запрашивает текст песни со стороннего API, добавляет песню в БД.
func (l musicLibrary) AddSong(song models.NewSongData) error {
	slog.Debug(
		"adding song...",
		"artist", song.Artist,
		"title", song.Song,
	)

	// Запрашивает текст песни со стороннего API. Ошибки не являются критическими.
	songWithLyrics, err := requestLyrics(song)
	if err != nil {
		slog.Error(
			"couldn't get lyrics for song",
			"artist", song.Artist,
			"title", song.Song,
			"err", err.Error(),
		)
	}

	// Добавляет песню в БД.
	err = config.Database.Insert(config.ExitCtx, songWithLyrics)
	if err != nil {
		slog.Error(
			"adding song: error",
			"artist", song.Artist,
			"title", song.Song,
			"err", err.Error(),
		)

		return err
	}

	slog.Debug(
		"adding song: success",
		"artist", song.Artist,
		"title", song.Song,
	)

	return nil
}

// Получение всех песен из БД. С пагинацией и фильтрацией.
func (l musicLibrary) GetSongs(params models.FullSongData) ([]models.PaginatedSongData, error) {
	slog.Debug(
		"getting songs...",
	)

	data, err := config.Database.SelectAll(config.ExitCtx, params)
	if err != nil {
		slog.Error(
			"getting songs: error",
			"err", err.Error(),
		)

		return nil, err
	}

	slog.Debug(
		"getting songs: success",
	)

	return data, nil
}

// Получение текста песни с пагинацией по куплетам.
func (l musicLibrary) GetLyrics(song models.NewSongData) ([]models.PaginatedLyrics, error) {
	slog.Debug(
		"getting lyrics...",
		"artist", song.Artist,
		"title", song.Song,
	)

	data, err := config.Database.SelectLyrics(config.ExitCtx, song)
	if err != nil {
		slog.Debug(
			"getting lyrics: error",
			"artist", song.Artist,
			"title", song.Song,
		)

		return nil, err
	}

	paginatedLyrics := paginateLyrics(data)

	slog.Debug(
		"getting lyrics: success",
		"artist", song.Artist,
		"title", song.Song,
	)

	return paginatedLyrics, nil
}

// Изменение данных песни.
func (l musicLibrary) ChangeSong(song models.NewSongData, params models.FullSongData) error {
	slog.Debug(
		"changing song...",
		"artist", song.Artist,
		"title", song.Song,
	)

	err := config.Database.Update(config.ExitCtx, song, params)
	if err != nil {
		slog.Debug(
			"changing song: error",
			"artist", song.Artist,
			"title", song.Song,
		)

		return err
	}

	slog.Debug(
		"changing song: success",
		"artist", song.Artist,
		"title", song.Song,
	)

	return nil
}

// Удаление песни.
func (l musicLibrary) DeleteSong(song models.NewSongData) error {
	slog.Debug(
		"deleting song...",
		"artist", song.Artist,
		"title", song.Song,
	)

	err := config.Database.Delete(config.ExitCtx, song)
	if err != nil {
		slog.Debug(
			"deleting song: error",
			"artist", song.Artist,
			"title", song.Song,
		)

		return err
	}

	slog.Debug(
		"deleting song: success",
		"artist", song.Artist,
		"title", song.Song,
	)

	return nil
}

// Возвращает новую библиотеку музыки.
func NewLibrary() models.MusicLibrary {
	return &musicLibrary{}
}
