package config

import (
	"os"
	"strconv"
)

// Config contains the configuration for the service
type Config struct {
	Port                  int
	ProfileServiceBaseurl string
	DBUsername            string
	DBPassword            string
	DBHost                string
	DBPort                int
	Database              string
}

// LoadConfig returns the config from the config.json file
func LoadConfig() Config {

	var config Config

	if _, ok := os.LookupEnv("PORT"); ok {
		portString := os.Getenv("PORT")
		port, err := strconv.Atoi(portString)
		if err != nil {
			config.Port = port
		}
	}

	if _, ok := os.LookupEnv("PROFIE_SERVICE_URL"); ok {
		config.ProfileServiceBaseurl = os.Getenv("PROFIE_SERVICE_URL")
	}

	if _, ok := os.LookupEnv("DB_USERNAME"); ok {
		config.DBUsername = os.Getenv("DB_USERNAME")
	}

	if _, ok := os.LookupEnv("DB_PASSWORD"); ok {
		config.DBPassword = os.Getenv("DB_PASSWORD")
	}

	if _, ok := os.LookupEnv("DB_HOST"); ok {
		config.DBHost = os.Getenv("DB_HOST")
	}

	if _, ok := os.LookupEnv("DB_PORT"); ok {
		portString := os.Getenv("DB_PORT")
		port, err := strconv.Atoi(portString)
		if err != nil {
			config.DBPort = port
		}
	}

	if _, ok := os.LookupEnv("DB"); ok {
		config.Database = os.Getenv("DB")
	}

	return config
}
