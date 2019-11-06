package repository

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/robrohan/HungryLegs/internal/models"
	migrate "github.com/rubenv/sql-migrate"
)

// OpenDatabase Open up the database connection
func OpenDatabase(driver string, connection string, post string, a *models.Athlete) (*sql.DB, error) {
	conn := strings.ReplaceAll(connection, "{athlete}", *a.FileSafeName)
	db, err := sql.Open(driver, conn)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}

	posts := strings.Split(post, ";")

	for _, ex := range posts {
		cmd := strings.ReplaceAll(ex, "{athlete}", *a.Name)
		_, err = db.Exec(cmd)
		if err != nil {
			log.Printf("%v\n", err.Error())
		}
	}

	err = UpdateAthleteStore(driver, db)
	if err != nil {
		log.Printf("Failed to upgrade the athlete store")
		return nil, err
	}

	return db, nil
}

// UpdateAthleteStore Run any migrations that need to run
func UpdateAthleteStore(driver string, db *sql.DB) error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db, driver, migrations, migrate.Up)
	if err != nil {
		log.Printf("Filed migrations\n")
		return err
	}
	log.Printf("Applied %d migrations\n", n)
	return nil
}
