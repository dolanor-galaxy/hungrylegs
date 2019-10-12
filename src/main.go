package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/therohans/HungryLegs/src/importer"
	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
)

// Athlete is the root level account
type Athlete struct {
	Name string
}

func (a *Athlete) OpenAthlete() (*sql.DB, error) {
	athletePath := filepath.Join("store", "athletes", a.Name+".db")
	db, err := sql.Open("sqlite3", athletePath)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}
	return db, nil
}

func (a *Athlete) UpdateAthleteStore(db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	if err != nil {
		log.Printf("Filed migrations\n")
		return err
	}
	log.Printf("Applied %d migrations\n", n)
	return nil
}

/////////////////////////////////////////////////////

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

	a := Athlete{
		Name: "professor_zoom",
	}

	// Open an athlete (an sqlite database)
	db, err := a.OpenAthlete()
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Run any migrations
	err = a.UpdateAthleteStore(db)
	if err != nil {
		log.Fatal(err)
	}

	// Put the API on top of the connection
	repo := repository.AthleteRepository{
		Db: db,
	}

	importer.ImportNewActivity(config, &repo)
}
