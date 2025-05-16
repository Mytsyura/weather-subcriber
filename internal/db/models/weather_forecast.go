package models

import "time"

type WeatherForecast struct {
	City        string    `json:"city"`
	Temperature float64   `json:"temperature"`
	Humidity    int       `json:"humidity"`
	WindSpeed   float64   `json:"wind_speed"`
	Description string    `json:"description"`
	ForecastFor time.Time `json:"forecast_for"`
	CreatedAt   time.Time `json:"created_at"`
}
