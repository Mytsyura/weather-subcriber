package databasehandler

import (
	"context"
	"fmt"

	"weather_subscription/internal/db/database_repository"
)

var dbHandler *databaseHandler

func Init(ctx context.Context) error {
	weatherServiceRepository, err := database_repository.New(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialized database repository: '%w'", err)
	}

	dbHandler = &databaseHandler{
		weatherServiceRepository: weatherServiceRepository,
	}

	return nil
}

func Close() {
	dbHandler.Close()
}

func Ping(ctx context.Context) error {
	return dbHandler.weatherServiceRepository.Ping(ctx)
}
