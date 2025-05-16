package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// LoadConfig tries to load config.yaml, otherwise falls back to defaults.
func LoadConfig() {
	viper.SetConfigName("config") // config.yaml
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config") // root

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("⚠️  Config file not found or unreadable, using defaults: %v\n", err)
		setDefaults()
	} else {
		fmt.Println("✅ config loaded:", viper.ConfigFileUsed())
	}
}

func setDefaults() {
	viper.Set("dbAddr", "localhost:5432")
	viper.Set("dbName", "weather_subscription")
	viper.Set("dbUser", "postgres")
	viper.Set("dbPass", "postgres")
	viper.Set("serverPort", "8080")
}
