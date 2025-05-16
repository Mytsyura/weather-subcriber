package database_repository

import (
	"context"

	infrastructure "weather_subscription/internal/db/database_repository/infrastracture"

	postgresql "weather_subscription/internal/db/database_repository/postgresql"
)

func New(ctx context.Context) (infrastructure.WeatherServiceRepository, error) {
	return postgresql.Init(ctx, "public")
}
