package database

import (
	"context"
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
		config.DBConnectionString+"?sslmode=disable")
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

// TODO: delete me
func (db postgresDB) DeleteMe(ctx context.Context) {
	// Предварительная проверка
	args := pgx.NamedArgs{
		"id":     1,
		"artist": "kiss",
		"song":   "sample text",
		"lyrics": "hello world",
	}
	row := db.Conn.QueryRow(ctx, q_insert, args)

	var song_id int
	err := row.Scan(&song_id)
	if err != nil {
		slog.Error(
			"failed to insert data",
			"err", err.Error,
			"id", 1,
		)
		return
	}
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
