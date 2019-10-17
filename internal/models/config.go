package models

import (
	"encoding/json"
	"log"
	"os"
)

// DatabaseConfig the connection configuration
type DatabaseConfig struct {
	Driver     string   `json:"driver"`
	Connection string   `json:"connection"`
	Post       []string `json:"post"`
}

// StaticConfig are global configs from the start up file
type StaticConfig struct {
	ImportDir   string         `json:"import"`
	RootAthlete string         `json:"root"`
	Database    DatabaseConfig `json:"database"`
}

// LoadConfig loads the start up configuration
func LoadConfig(file string) (*StaticConfig, error) {
	var config StaticConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Printf("Opening config file failed %v\n", err.Error())
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return &config, nil
}
