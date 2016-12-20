package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// Config contains the configuration for the service
type Config struct {
	Port       int    `json:"port"`
	DBUsername string `json:"db_username"`
	DBPassword string `json:"db_password"`
	DBHost     string `json:"db_host"`
	DBPort     int    `json:"db_port"`
	Database   string `json:"db"`
	SecretKey  string `json:"secret_key"`
}

// LoadConfig returns the config from the config.json file
func LoadConfig() Config {

	file, err := ioutil.ReadFile("./config/config.json")
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	return config
}
