package repository

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/robrohan/HungryLegs/internal/models"
)

// OpenDatabase Open up the database connection
func OpenDatabase(config *models.StaticConfig, a *models.Athlete) (*sql.DB, error) {
	conn := strings.ReplaceAll(config.Database.Connection, "{athlete}", *a.FileSafeName)
	db, err := sql.Open(config.Database.Driver, conn)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}

	for _, ex := range config.Database.Post {
		cmd := strings.ReplaceAll(ex, "{athlete}", *a.Name)
		_, err = db.Exec(cmd)
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	}

	err = UpdateAthleteStore(config, db)
	if err != nil {
		log.Printf("Failed to upgrade the athlete store")
		return nil, err
	}

	return db, nil
}

// UpdateAthleteStore Run any migrations that need to run
func UpdateAthleteStore(config *models.StaticConfig, db *sql.DB) error {
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
