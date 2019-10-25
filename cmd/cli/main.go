package main

import (
	"log"
	"time"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robrohan/HungryLegs/internal/importer"
	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
)

func main() {
	log.Printf("Starting HungryLegs...")

	config, err := models.LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Creating root athlete...")
	rootAthlete := models.NewAthlete(&config.RootAthlete, &config.RootAthlete)
	db, err := repository.OpenDatabase(config, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Put the API on top of the connection
	repo := repository.Attach(*rootAthlete.Alterego, db, config)

	log.Printf("Starting import...")
	start := time.Now()
	// Launch the activity importer
	importer.ImportActivites(config.ImportDir, repo)
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Full import took %v", elapsed)
}
