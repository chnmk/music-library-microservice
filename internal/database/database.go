package database

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

/*
TODO: Обогащенную информацию положить в БД postgres (структура БД должна
быть создана путем миграций при старте сервиса)
*/

const q_insert = `
	INSERT INTO muslib(id, artist, song, lyrics)
	VALUES (@id, @artist, @song, @lyrics)
	RETURNING id
`

const q_restore = `
	SELECT * FROM muslib
	LIMIT @lim
`

type postgresDB struct {
	Conn *pgxpool.Pool
}

func NewDatabase(ctx context.Context) models.Database {
	slog.Debug(
		"database setup",
		"string", config.DBConnectionString,
	)

	// Миграции
	m, err := migrate.New(
		"file://migrations",
		// TODO: ssl
		config.DBConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		slog.Info(err.Error())
	}

	// Подключение к БД
	conn, err := pgxpool.New(ctx, config.DBConnectionString)
	if err != nil {
		slog.Error(
			"unable to connect to database",
			"err", err.Error(),
		)
	}

	return &postgresDB{Conn: conn}
}

func (db postgresDB) AddSong(ctx context.Context, id int, song models.SongData) error {
	args := pgx.NamedArgs{
		"id":     id,
		"artist": song.Group,
		"song":   song.Song,
		"lyrics": song.Lyrics,
	}

	row := db.Conn.QueryRow(ctx, q_insert, args)

	var song_id int

	err := row.Scan(&song_id)
	if err != nil {
		return err
	}

	return nil
}

func (db postgresDB) RestoreData(ctx context.Context) (int, map[int]models.SongData, error) {
	var maxID int
	result := make(map[int]models.SongData)

	args := pgx.NamedArgs{
		"lim": config.MaxEntries,
	}

	rows, err := db.Conn.Query(ctx, q_restore, args)
	if err != nil {
		return 0, nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var song models.SongData

		err = rows.Scan(&id, &song.Group, &song.Song, &song.Lyrics)
		if err != nil {
			return 0, nil, err
		}

		result[id] = song
		if id > maxID {

			maxID = id
		}
	}

	if rows.Err() != nil {
		return 0, nil, err
	}

	if len(result) == 0 {
		return 0, nil, errors.New("data not found")
	}

	return maxID, result, nil
}
