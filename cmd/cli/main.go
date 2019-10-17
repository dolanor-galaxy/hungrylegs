package main

import (
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/therohans/HungryLegs/internal/importer"
	"github.com/therohans/HungryLegs/internal/models"
	"github.com/therohans/HungryLegs/internal/repository"
)

func main() {
	log.Printf("Starting HungryLegs...")

	config, err := models.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", config)

	rootAthlete := models.NewAthlete("Professor Zoom")
	db, err := repository.OpenDatabase(config, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Put the API on top of the connection
	repo := repository.Attach(rootAthlete, db, config)

	// Launch the activity importer
	importer.ImportActivites(config.ImportDir, repo)
}
