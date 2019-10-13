package main

import (
	"encoding/json"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/therohans/HungryLegs/src/importer"
	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
)

func loadConfig(file string) (*models.StaticConfig, error) {
	var config models.StaticConfig
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		log.Printf("Opening log file failed %v\n", err.Error())
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return &config, nil
}

/////////////////////////////////////////////////////

func main() {
	// Load configs
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", config)

	athlete := models.OpenAthlete("Professor Zoom")
	defer athlete.Close()
	// Put the API on top of the connection
	repo := repository.Attach(athlete)
	// Launch the activity importer
	importer.ImportNewActivity(config, repo)

	//////////////////////
	athlete2 := models.OpenAthlete("Punkin Pie")
	defer athlete2.Close()
	// Put the API on top of the connection
	repo2 := repository.Attach(athlete2)
	// Launch the activity importer
	importer.ImportNewActivity(config, repo2)
}
