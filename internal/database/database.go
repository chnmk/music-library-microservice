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

// Запрос для чтения данных.
const q_select_data = `
	SELECT a.name, s.title, s.lyrics, s.release_date, s.link FROM artists a
	LEFT JOIN songs s
	ON a.id = s.artist_id
	WHERE a.name LIKE @name AND s.title LIKE @title AND s.lyrics LIKE @lyrics 
	AND s.release_date LIKE @release_date AND s.link LIKE @link
`

// Запрос для чтения одной песни.
const q_select_song = `
	SELECT a.name, s.title, s.lyrics, s.release_date, s.link FROM artists a
	LEFT JOIN songs s
	ON a.id = s.artist_id
	WHERE a.name = @name AND s.title = @title
`

// Запрос для получения ID исполнителя по имени.
const q_select_artist_id = `
	SELECT id FROM artists 
	WHERE name = @name
`

// Запрос для проверки количества песен у исполнителя.
const q_select_artist_song_count = `
	SELECT COUNT(*) FROM songs 
	WHERE artist_id = @artist_id
`

// Запрос для добавления исполнителя.
const q_insert_artist = `
	INSERT INTO artists(name)
	VALUES (@name)
	RETURNING id
`

// Запрос для добавления песни.
const q_insert_song = `
	INSERT INTO songs(artist_id, title, lyrics, release_date, link)
	VALUES (@artist_id, @title, @lyrics, @release_date, @link)
	RETURNING id
`

// Запрос для редактирования исполнителя.
const q_update_artist = `
	UPDATE artists
	SET name = @name
	WHERE name = @old_name
	RETURNING id
`

// Запрос для редактирования песни.
const q_update_song = `
	UPDATE songs
	SET title = @title , lyrics = @lyrics, release_date = @release_date, link = @link
	WHERE title = @old_title
	AND artist_id = @artist_id
	RETURNING id
`

// Запрос для удаления исполнителя.
const q_delete_artist = `
	DELETE FROM artists
	WHERE name = @name
	RETURNING id
`

// Запрос для удаления песни.
const q_delete_song = `
	DELETE FROM songs
	WHERE artist_id = @artist_id AND title = @title
	RETURNING id
`

// Имплементация интерфейса models.Database.
type postgresDB struct {
	Conn *pgxpool.Pool
}

// Добавление песни в БД.
func (db postgresDB) Insert(ctx context.Context, song models.FullSongData) error {
	// Проверка, что исполнителя нет.
	args := pgx.NamedArgs{"name": song.Artist}

	var artist_id int
	err := db.Conn.QueryRow(ctx, q_select_artist_id, args).Scan(&artist_id)

	if err != nil {
		// Добавление исполнителя.
		err = db.Conn.QueryRow(ctx, q_insert_artist, args).Scan(&artist_id)
		if err != nil {
			return err
		}
	}

	// Добавление песни.
	args = pgx.NamedArgs{
		"artist_id":    artist_id,
		"title":        song.Song,
		"lyrics":       song.Lyrics,
		"release_date": song.ReleaseDate,
		"link":         song.Link,
	}

	var song_id int
	err = db.Conn.QueryRow(ctx, q_insert_song, args).Scan(&song_id)
	if err != nil {
		return err
	}

	return nil
}

// Получение песен.
func (db postgresDB) SelectAll(ctx context.Context, params models.FullSongData) ([]models.PaginatedSongData, error) {
	var result []models.PaginatedSongData

	// Запрос данных из БД.
	args := pgx.NamedArgs{
		"name":         "%" + params.Artist + "%",
		"title":        "%" + params.Song + "%",
		"lyrics":       "%" + params.Lyrics + "%",
		"release_date": "%" + params.ReleaseDate + "%",
		"link":         "%" + params.Link + "%",
	}

	rows, err := db.Conn.Query(ctx, q_select_data, args)
	if err != nil {
		return result, err
	}

	// Чтение данных с пагинацией.
	var currentPage int

	noData := true
	result = append(result, models.PaginatedSongData{Page: currentPage + 1}) // Нумерация страниц начинается с единицы

	defer rows.Close()
	for rows.Next() {
		noData = false
		var song models.FullSongData

		err = rows.Scan(&song.Artist, &song.Song, &song.Lyrics, &song.ReleaseDate, &song.Link)
		if err != nil {
			return result, err
		}

		if len(result[currentPage].Entries) >= config.PageSize {
			currentPage++
			result = append(result, models.PaginatedSongData{Page: currentPage + 1})
		}

		result[currentPage].Entries = append(result[currentPage].Entries, song)
	}

	if rows.Err() != nil {
		return result, rows.Err()
	}

	if noData {
		return result, errors.New("data not found")
	}

	return result, nil
}

// Получение текста одной песни.
func (db postgresDB) SelectLyrics(ctx context.Context, song models.NewSongData) (string, error) {
	var result string

	// Запрос данных из БД.
	args := pgx.NamedArgs{
		"name":  song.Artist,
		"title": song.Song,
	}

	var new models.FullSongData
	err := db.Conn.QueryRow(ctx, q_select_song, args).Scan(&new.Artist, &new.Song, &new.Lyrics, &new.ReleaseDate, &new.Link)
	if err != nil {
		return result, err
	}

	return new.Lyrics, nil
}

// Редактирование песни.
func (db postgresDB) Update(ctx context.Context, song models.NewSongData, params models.FullSongData) error {
	// Ищёт нужную песню.
	args := pgx.NamedArgs{"name": song.Artist, "title": song.Song}

	var new models.FullSongData
	err := db.Conn.QueryRow(ctx, q_select_song, args).Scan(&new.Artist, &new.Song, &new.Lyrics, &new.ReleaseDate, &new.Link)
	if err != nil {
		return err
	}

	// Проверяет какие поля стоит изменить.
	if params.Artist != "" {
		new.Artist = params.Artist
	}
	if params.Song != "" {
		new.Song = params.Song
	}
	if params.Lyrics != "" {
		new.Lyrics = params.Lyrics
	}
	if params.ReleaseDate != "" {
		new.ReleaseDate = params.ReleaseDate
	}
	if params.Link != "" {
		new.Link = params.Link
	}

	// Редактирует исполнителя.
	args = pgx.NamedArgs{
		"old_name": song.Artist,
		"name":     new.Artist,
	}

	var artist_id int
	err = db.Conn.QueryRow(ctx, q_update_artist, args).Scan(&artist_id)
	if err != nil {
		return err
	}

	// Редактирует песню.
	args = pgx.NamedArgs{
		"artist_id":    artist_id,
		"old_title":    song.Song,
		"title":        new.Song,
		"lyrics":       new.Lyrics,
		"release_date": new.ReleaseDate,
		"link":         new.Link,
	}

	var song_id int
	err = db.Conn.QueryRow(ctx, q_update_song, args).Scan(&song_id)
	if err != nil {
		return err
	}

	return nil
}

// Удаление песни. Если у исполнителя не осталось песен, исполнитель тоже удаляется.
func (db postgresDB) Delete(ctx context.Context, song models.NewSongData) error {
	// Ищет нужного исполнителя.
	args := pgx.NamedArgs{"name": song.Artist}

	var artist_id int
	err := db.Conn.QueryRow(ctx, q_select_artist_id, args).Scan(&artist_id)
	if err != nil {
		return err
	}

	// Удаляет песню.
	args = pgx.NamedArgs{"artist_id": artist_id, "title": song.Song}

	var song_id int
	err = db.Conn.QueryRow(ctx, q_delete_song, args).Scan(&song_id)
	if err != nil {
		return err
	}

	// Проверяет остались ли песни у исполнителя.
	args = pgx.NamedArgs{"artist_id": artist_id}

	var song_count int
	err = db.Conn.QueryRow(ctx, q_select_artist_song_count, args).Scan(&song_count)
	if err != nil {
		return err
	}

	// Удаляет исполнителя.
	args = pgx.NamedArgs{"name": song.Artist}

	err = db.Conn.QueryRow(ctx, q_delete_artist, args).Scan(&artist_id)
	if err != nil {
		return err
	}
	return nil
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
