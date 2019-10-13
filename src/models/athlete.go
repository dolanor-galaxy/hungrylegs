package models

import (
	"database/sql"
	"encoding/base64"
	"log"
	"path/filepath"

	migrate "github.com/rubenv/sql-migrate"
)

// Athlete is the root level account
type Athlete struct {
	Name         string
	FileSafeName string
	Db           *sql.DB
}

func OpenAthlete(name string) *Athlete {
	a := Athlete{
		Name:         name,
		FileSafeName: base64.URLEncoding.EncodeToString([]byte(name)),
	}

	// Open the database and apply migrations is needed
	db, err := openDatabase(&a)
	if err != nil {
		log.Fatal(err)
	}
	a.Db = db
	return &a
}

// Close release the database and what not
func (a *Athlete) Close() {
	a.Db.Close()
}

func openDatabase(a *Athlete) (*sql.DB, error) {
	athletePath := filepath.Join("store", "athletes", a.FileSafeName+".db")
	// db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	db, err := sql.Open("sqlite3", athletePath)
	if err != nil {
		log.Printf("Failed to open athlete store")
		return nil, err
	}
	_, err = db.Exec("PRAGMA synchronous = OFF")
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
	_, err = db.Exec("PRAGMA journal_mode = MEMORY")
	if err != nil {
		log.Printf("%v\n", err.Error())
	}
	_, err = db.Exec("PRAGMA cache_size = -16000")
	if err != nil {
		log.Printf("%v\n", err.Error())
	}

	// Ensure the database is up2date
	err = updateAthleteStore(a, db)
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}

func updateAthleteStore(a *Athlete, db *sql.DB) error {
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
