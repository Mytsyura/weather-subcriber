package postgresql

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/viper"

	infrastructure "weather_subscription/internal/db/database_repository/infrastracture"
)

const (
	dbEndpointKey = "dbAddr"
	dbNameKey     = "dbName"
	dbUserKey     = "dbUser"
	dbPasswordKey = "dbPass"
)

func Init(ctx context.Context, schemaName string) (infrastructure.WeatherServiceRepository, error) {
	return initSchemaConnection(ctx, schemaName)
}

func initSchemaConnection(ctx context.Context, schemaName string) (infrastructure.WeatherServiceRepository, error) {
	const funcName = "postgresql.initSchemaConnection"

	// Check environment variables for non-local environment
	dbEndpoint := os.Getenv(dbEndpointKey)
	dbUser := os.Getenv(dbUserKey)
	dbPassword := os.Getenv(dbPasswordKey)
	dbName := os.Getenv(dbNameKey)

	if dbEndpoint == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		setDefaults()
	}

	viper.Set("dbAddr", dbEndpoint)
	viper.Set("dbName", dbName)
	viper.Set("dbUser", dbUser)
	viper.Set("dbPass", dbPassword)

	repo, err := NewPostgresRepo(ctx, schemaName)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to init postgres db: %w", funcName, err)
	}

	return &postgresqlWeatherServiceRepository{repo: repo}, nil
}

func setDefaults() {
	viper.Set("dbAddr", "localhost:5432")
	viper.Set("dbName", "weather_subscription")
	viper.Set("dbUser", "postgres")
	viper.Set("dbPass", "postgres")
	viper.Set("serverPort", "8080")
}
