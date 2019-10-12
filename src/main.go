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

func main() {
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v", config)

	db, err := openAthlete("professor_zoom")
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	err = updateAthleteStore(db)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(config.ImportDir)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.AthleteRepository{
		Db: db,
	}

	for _, f := range files {
		name := f.Name()
		name = strings.ToLower(name)

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
	}

	// statement, _ := db.Prepare("create table if not exists people (id INTEGER PRIMARY KEY, firstname TEXT)")
	// statement.Exec()
	// statement, _ := db.Prepare("INSERT INTO PEOPLE (firstname) VALUES (?)")
	// statement.Exec("Rob")
	// rows, _ := db.Query("SELECT id, firstname FROM people")

	// p := P{}
	// for rows.Next() {
	// 	rows.Scan(&p.ID, &p.FirstName)
	// 	// fmt.Println(strconv.Itoa(id) + ": " + firstname)
	// 	fmt.Println(p)
	// }
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
