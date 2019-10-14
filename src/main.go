package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/therohans/HungryLegs/src/importer"
	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
)

func loadConfig(file string) (*models.StaticConfig, error) {
	var config models.StaticConfig
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

/////////////////////////////////////////////////////

func main() {
	log.Printf("Starting HungryLegs...")

	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", config)

	rootAthlete := models.NewAthlete("Professor Zoom")
	db, err := openDatabase(config, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// defer athlete.Close()
	// Put the API on top of the connection
	repo := repository.Attach(db)

	// Launch the activity importer
	importer.ImportActivites(config.ImportDir, repo)

	//////////////////////
	// athlete2 := models.OpenAthlete("Punkin Pie")
	// defer athlete2.Close()
	// // Put the API on top of the connection
	// repo2 := repository.Attach(athlete2)
	// // Launch the activity importer
	// importer.ImportActivites(config.ImportDir, repo2)
}

func openDatabase(config *models.StaticConfig, a *models.Athlete) (*sql.DB, error) {
	conn := strings.ReplaceAll(config.Database.Connection, "{athlete}", a.FileSafeName)
	db, err := sql.Open(config.Database.Driver, conn)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}

	for _, ex := range config.Database.Post {
		cmd := strings.ReplaceAll(ex, "{athlete}", a.Name)
		_, err = db.Exec(cmd)
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	}

	err = updateAthleteStore(config, db)
	if err != nil {
		log.Printf("Failed to upgrade the athlete store")
		return nil, err
	}

	return db, nil
}

func updateAthleteStore(config *models.StaticConfig, db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, config.Database.Driver, migrations, migrate.Up)
	if err != nil {
		log.Printf("Filed migrations\n")
		return err
	}
	log.Printf("Applied %d migrations\n", n)
	return nil
}
