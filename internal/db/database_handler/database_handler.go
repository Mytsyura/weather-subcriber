package databasehandler

import (
	infrastructure "weather_subscription/internal/db/database_repository/infrastracture"
)

type databaseHandler struct {
	weatherServiceRepository infrastructure.WeatherServiceRepository
}

func (d *databaseHandler) Close() {
	d.weatherServiceRepository.Close()
}
