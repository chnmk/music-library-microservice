package database

import "github.com/chnmk/music-library-microservice/internal/models"

/*
TODO: Обогащенную информацию положить в БД postgres (структура БД должна
быть создана путем миграций при старте сервиса)
*/

type postgresDB struct {
	// TODO
}

func NewDatabase() models.Database {
	return &postgresDB{}
}
