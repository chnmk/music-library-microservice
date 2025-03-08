package database

import (
	"context"
	"errors"
	"log/slog"

	"github.com/chnmk/music-library-microservice/internal/config"
	"github.com/chnmk/music-library-microservice/internal/models"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Строка добавления песни в БД.
const q_insert = `
	INSERT INTO muslib(id, artist, song, lyrics)
	VALUES (@id, @artist, @song, @lyrics)
	RETURNING id
`

// Строка восстановления данных из БД.
const q_restore = `
	SELECT * FROM muslib
	LIMIT @lim
`

// Имплементация интерфейса models.Database.
type postgresDB struct {
	Conn *pgxpool.Pool
}

// Добавление песни в БД.
func (db postgresDB) AddSong(ctx context.Context, id int, song models.SongData) error {
	slog.Info(
		"adding song to DB",
		"id", id,
	)

	args := pgx.NamedArgs{
		"id":     id,
		"artist": song.Group,
		"song":   song.Song,
		"lyrics": song.Lyrics,
	}

	row := db.Conn.QueryRow(ctx, q_insert, args)

	// Проверка результата.
	var song_id int
	err := row.Scan(&song_id)
	if err != nil {
		return err
	}

	slog.Info(
		"song successfully added to DB",
		"id", id,
	)

	return nil
}

// Восстановление данных из БД.
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

// Запускает миграции, устанавливает подключение, возвращает новую БД.
func NewDatabase(ctx context.Context) models.Database {
	slog.Debug(
		"database setup",
		"string", config.DBConnectionString,
	)

	// Миграции.
	m, err := migrate.New(
		"file://migrations",
		config.DBConnectionString)
	if err != nil {
		slog.Error(
			"fatal error",
			"err", err,
		)
		config.Exit()
	}
	if err := m.Up(); err != nil {
		slog.Info(err.Error())
	}

	// Подключение к БД.
	conn, err := pgxpool.New(ctx, config.DBConnectionString)
	if err != nil {
		slog.Error(
			"unable to connect to database",
			"err", err.Error(),
		)
		config.Exit()
	}

	return &postgresDB{Conn: conn}
}
