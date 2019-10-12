package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/therohans/HungryLegs/src/importer"
	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
	"github.com/tormoder/fit"

	migrate "github.com/rubenv/sql-migrate"
)

func openAthlete(name string) (*sql.DB, error) {
	athletePath := filepath.Join("store", "athletes", name+".db")
	db, err := sql.Open("sqlite3", athletePath)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}
	return db, nil
}

func updateAthleteStore(db *sql.DB) error {
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

func importNewActivity(config *models.StaticConfig, repo *repository.AthleteRepository) {
	log.Println("Beginning import of new files...")
	files, err := ioutil.ReadDir(config.ImportDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()
		name = strings.ToLower(name)

		have, err := repo.HasImported(name)
		if err != nil {
			log.Fatal(err)
		}

		if have == false {
			start := time.Now()
			if strings.HasSuffix(name, ".tcx") {
				tcxFile := importer.TcxFile{}
				err := tcxFile.Import(filepath.Join(config.ImportDir, name), repo)
				if err != nil {
					log.Fatal(err)
				}
			} else if strings.HasSuffix(name, ".fit") {
				fitFile := importer.FitFile{}
				err := fitFile.Import(filepath.Join(config.ImportDir, name), repo)
				if err != nil {
					log.Fatal(err)
				}
			}
			repo.RecordImport(name)

			t := time.Now()
			elapsed := t.Sub(start)
			log.Printf("%v took %v", name, elapsed)
		} else {
			log.Printf("Already imported %v\n", name)
		}
	}
	log.Println("Done import")
}

func main() {
	// Load configs
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", config)

	// Open an athlete (an sqlite database)
	db, err := openAthlete("professor_zoom")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	// Run any migrations
	err = updateAthleteStore(db)
	if err != nil {
		log.Fatal(err)
	}

	// Put the API on top of the connection
	repo := repository.AthleteRepository{
		Db: db,
	}

	importNewActivity(config, &repo)
}

func Example() {
	// Read our FIT test file data
	testFile := filepath.Join("import", "run_no_heart.FIT")
	testData, err := ioutil.ReadFile(testFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Decode the FIT file data
	fit, err := fit.Decode(bytes.NewReader(testData))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Inspect the TimeCreated field in the FileId message
	fmt.Println(fit.FileId.TimeCreated)

	fmt.Println(fit.Type())

	// Inspect the dynamic Product field in the FileId message
	fmt.Println(fit.FileId.GetProduct())

	// Inspect the FIT file type
	// fmt.Println(fit.FileType())

	// Get the actual activity
	activity, err := fit.Activity()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the latitude and longitude of the first Record message
	for _, record := range activity.Records {
		fmt.Println(record.PositionLat)
		fmt.Println(record.PositionLong)
		break
	}

	// Print the sport of the first Session message
	for _, session := range activity.Sessions {
		fmt.Println(session.Sport)
		break
	}
}
