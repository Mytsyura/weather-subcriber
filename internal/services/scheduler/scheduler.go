package scheduler

import (
	"context"
	"log"
	"time"

	databasehandler "weather_subscription/internal/db/database_handler"
	"weather_subscription/internal/db/models"
	"weather_subscription/internal/services/email"
	weatherClient "weather_subscription/internal/weatherClient"
)

type WeatherScheduler struct {
	weatherClient *weatherClient.APIClient
	emailService  *email.EmailService
}

func NewWeatherScheduler(weatherClient *weatherClient.APIClient, emailService *email.EmailService) *WeatherScheduler {
	return &WeatherScheduler{
		weatherClient: weatherClient,
		emailService:  emailService,
	}
}

func (s *WeatherScheduler) Start(ctx context.Context) {
	// Start daily check at midnight
	go s.scheduleDailyCheck(ctx)

	// Start hourly check
	go s.scheduleHourlyCheck(ctx)
}

func (s *WeatherScheduler) scheduleDailyCheck(ctx context.Context) {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		time.Sleep(time.Until(next))

		// Get all daily subscriptions
		subscriptions, err := databasehandler.ListActiveSubscriptions(ctx)
		if err != nil {
			log.Printf("Error fetching daily subscriptions: %v", err)
			continue
		}

		for _, sub := range subscriptions {
			if sub.Frequency == models.Daily {
				s.sendWeatherUpdate(ctx, sub)
			}
		}
	}
}

func (s *WeatherScheduler) scheduleHourlyCheck(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Get all hourly subscriptions
			subscriptions, err := databasehandler.ListActiveSubscriptions(ctx)
			if err != nil {
				log.Printf("Error fetching hourly subscriptions: %v", err)
				continue
			}

			for _, sub := range subscriptions {
				if sub.Frequency == models.Hourly {
					s.sendWeatherUpdate(ctx, sub)
				}
			}
		}
	}
}

func (s *WeatherScheduler) sendWeatherUpdate(ctx context.Context, subscription *models.Subscription) {
	// Get weather data
	weather, _, err := s.weatherClient.APIsApi.RealtimeWeather(ctx, subscription.City, nil)
	if err != nil {
		log.Printf("Error fetching weather for %s: %v", subscription.City, err)
		return
	}

	// Create weather forecast model
	forecast := &models.WeatherForecast{
		City:        subscription.City,
		Temperature: weather.Current.TempC,
		Description: weather.Current.Condition.Text,
		Humidity:    int(weather.Current.Humidity),
		WindSpeed:   weather.Current.WindKph,
	}

	// Send email
	if err := s.emailService.SendWeatherUpdate(ctx, subscription.Email, forecast); err != nil {
		log.Printf("Error sending weather update to %s: %v", subscription.Email, err)
		return
	}
}
