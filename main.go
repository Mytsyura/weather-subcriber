package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"weather_subscription/config"
	databasehandler "weather_subscription/internal/db/database_handler"
	models "weather_subscription/internal/db/models"
	"weather_subscription/internal/services/email"

	"weather_subscription/internal/services/scheduler"
	weatherClient "weather_subscription/internal/weatherClient"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	serverPortKey = "serverPort"
	receiverKey   = "from"
	passwordKey   = "password"
	smtpHostKey   = "smtphost"
	smtpPortKey   = "smtpport"
)

var emailService *email.EmailService

func main() {
	config.LoadConfig()

	// Initialize email service
	env := os.Getenv("ENV")
	if env == "" {
		env = "local" // Default to local environment
	}

	var emailConfig map[string]string
	if env == "local" || env == "development" {
		emailConfig = viper.GetStringMapString("email.local")
	} else {
		emailConfig = viper.GetStringMapString("email.production")
	}

	emailService = email.NewEmailService(
		emailConfig[receiverKey],
		emailConfig[passwordKey],
		emailConfig[smtpHostKey],
		emailConfig[smtpPortKey],
	)

	// Initialize database
	err := databasehandler.Init(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer databasehandler.Close()

	// Initialize WeatherAPI client
	apiKey := viper.GetString("weather_api.key")
	if apiKey == "" {
		log.Fatal("Weather API key not found in configuration")
		return
	}
	configuration := weatherClient.NewConfiguration()
	configuration.AddDefaultHeader("key", apiKey)

	weatherClient := weatherClient.NewAPIClient(configuration)
	if weatherClient == nil {
		log.Fatalf("Failed to create weather client")
		return
	}

	// Initialize and start scheduler
	scheduler := scheduler.NewWeatherScheduler(weatherClient, emailService)
	ctx := context.Background()
	go scheduler.Start(ctx)

	// Setup routes and start server
	startAPIServer()
}

func startAPIServer() {
	router := gin.Default()

	// Serve static files
	router.Static("/frontend", "./frontend")
	router.LoadHTMLGlob("frontend/*.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// Initialize WeatherAPI client
	apiKey := viper.GetString("weather_api.key")
	if apiKey == "" {
		log.Fatal("Weather API key not found in configuration")
		return
	}
	configuration := weatherClient.NewConfiguration()
	configuration.AddDefaultHeader("key", apiKey)

	weatherClient := weatherClient.NewAPIClient(configuration)

	if weatherClient == nil {
		log.Fatalf("Failed to create weather client")
		return
	}

	//emailService := &EmailService{}
	//scheduler := scheduler.NewWeatherScheduler(weatherClient, emailService)

	registerRoutes(router, weatherClient)
	// Start scheduler in background
	//ctx := context.Background()
	//go scheduler.Start(ctx)

	port := viper.GetString(serverPortKey)
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

func registerRoutes(router *gin.Engine, weatherClient *weatherClient.APIClient) {
	router.GET("/health", healthCheck())

	router.GET("/api/weather/:city", getWeather(weatherClient))
	router.POST("/api/subscribe", subscribe())
	router.GET("/api/unsubscribe/:email", unsubscribe())
	router.GET("/api/confirm/:token", confirm())
}

func healthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	}
}

func subscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Email     string                       `json:"email" binding:"required,email"`
			City      string                       `json:"city" binding:"required"`
			Frequency models.SubscriptionFrequency `json:"frequency" binding:"required,oneof=daily hourly"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := databasehandler.CreateSubscription(c.Request.Context(), req.Email, req.City, req.Frequency)
		if err != nil {
			if err.Error() == "subscription already exists" {
				c.JSON(http.StatusConflict, gin.H{"error": "Subscription already exists"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		// Send confirmation email
		if err := emailService.SendConfirmationEmail(req.Email, *token); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Failed to send confirmation email: %v", err)})
		}

		c.JSON(http.StatusCreated, gin.H{
			"status": "Subscription created successfully. Please check your email to confirm.",
			"token":  *token,
		})
	}
}

func unsubscribe() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Param("email")

		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email is required"})
			return
		}

		if err := databasehandler.DeleteSubscription(c.Request.Context(), email); err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Subscription is cancelled successfully"})
	}
}

func confirm() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Param("token")
		if err := databasehandler.ConfirmSubscription(c.Request.Context(), token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "Subscription is confirmed"})
	}
}

func getWeather(weatherClient *weatherClient.APIClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		city := c.Param("city")
		weather, response, err := weatherClient.APIsApi.RealtimeWeather(c.Request.Context(), city, nil)
		if err != nil && response != nil {
			c.JSON(response.StatusCode, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, weather)
	}
}

// EmailService implements the EmailService interface
type EmailService struct{}

func (s *EmailService) SendWeatherUpdate(ctx context.Context, email string, forecast *models.WeatherForecast) error {
	// TODO: Implement actual email sending
	log.Printf("Sending weather update to %s for %s: %.1f°C, %s",
		email, forecast.City, forecast.Temperature, forecast.Description)
	return nil
}
