package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tormoder/fit"

	migrate "github.com/rubenv/sql-migrate"
)

type P struct {
	ID        int
	FirstName string
}

func main() {
	Example()

	db, _ := sql.Open("sqlite3", "./bogo.db")

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	fmt.Printf("%v\n", migrations)

	n, _ := migrate.Exec(db, "sqlite3", migrations, migrate.Up)
	fmt.Printf("Applied %d migrations!\n", n)

	// statement, _ := db.Prepare("create table if not exists people (id INTEGER PRIMARY KEY, firstname TEXT)")
	// statement.Exec()
	// statement, _ := db.Prepare("INSERT INTO PEOPLE (firstname) VALUES (?)")
	// statement.Exec("Rob")
	rows, _ := db.Query("SELECT id, firstname FROM people")

	p := P{}
	for rows.Next() {
		rows.Scan(&p.ID, &p.FirstName)
		// fmt.Println(strconv.Itoa(id) + ": " + firstname)
		fmt.Println(p)
	}
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
