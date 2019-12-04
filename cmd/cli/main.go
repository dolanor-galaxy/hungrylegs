package main

import (
	"expvar"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	"github.com/robrohan/HungryLegs/internal/importer"
	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
)

// will be replaced with git hash
var build = "develop"

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "HL : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	cfg := models.Config

	if err := conf.Parse(os.Args[1:], "HL", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("HL", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// =========================================================================
	// App Starting
	expvar.NewString("build").Set(build)
	log.Printf("Starting HungryLegs initializing : version %q", build)
	defer log.Println("Completed")

	out, err := conf.String(&cfg)
	if err != nil {
		return errors.Wrap(err, "generating config for output")
	}
	log.Printf("Config :\n%v\n", out)

	// =========================================================================
	// Root athlete
	log.Printf("Creating root athlete...")
	rootAthlete := models.NewAthlete(&cfg.Base.Root, &cfg.Base.Root)

	// =========================================================================
	// Start Database
	log.Println("Initializing database support")
	db, err := repository.OpenDatabase(cfg.DB.Driver, cfg.DB.Connection, cfg.DB.Post, rootAthlete)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Printf("Database Stopping : %s", cfg.DB.Connection)
		db.Close()
	}()

	// Put the API on top of the connection
	repo := repository.Attach(*rootAthlete.Alterego, db, cfg.DB.Driver)

	log.Printf("Starting import...")
	start := time.Now()
	// Launch the activity importer
	importer.ImportActivites(log, cfg.Base.Import, repo)
	t := time.Now()
	elapsed := t.Sub(start)
	log.Printf("Full import took %v", elapsed)

	return nil
}
