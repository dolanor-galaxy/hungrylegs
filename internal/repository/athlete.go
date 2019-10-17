package repository

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/therohans/HungryLegs/internal/models"
)

type AthleteRepository struct {
	Athlete            *models.Athlete
	Db                 *sql.DB
	hasImportedQuery   *sql.Stmt
	recordImportQuery  *sql.Stmt
	addActivityQuery   *sql.Stmt
	addLapQuery        *sql.Stmt
	addTrackPointQuery *sql.Stmt
}

func prepareQuery(query string, db *sql.DB) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}

var re = regexp.MustCompile(`(\$[0-9]+)`)

func sqlForDriver(query string, config *models.StaticConfig) string {
	if config.Database.Driver == "sqlite3" {
		noSchema := strings.ReplaceAll(query, "\"%v\".", "")
		placeholders := re.ReplaceAllString(noSchema, "?")
		// We need to add the schema back somehow to the sprintfs
		// don't fail
		return placeholders + "\n -- %v\n"
	}
	return query
}

// Attach creates a new repository and sets up needed bits
func Attach(athlete *models.Athlete, db *sql.DB, config *models.StaticConfig) *AthleteRepository {
	a := AthleteRepository{
		Athlete: athlete,
		Db:      db,
	}

	a.hasImportedQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		SELECT import_time FROM "%v".fileimport WHERE file_name = $1
	`, config), athlete.Name), db)

	a.recordImportQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".fileimport (
			import_time, "file_name"
		) VALUES ($1, $2)
	`, config), athlete.Name), db)

	a.addActivityQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".activity (
			uuid, suuid, sport, "time", device
		) VALUES ($1, $2, $3, $4, $5)
	`, config), athlete.Name), db)

	a.addLapQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".lap (
			"time", "start", total_time, dist, calories, max_speed,
			avg_hr, max_hr, intensity, trigger, activity_uuid
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, config), athlete.Name), db)

	a.addTrackPointQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".trackpoint (
			"time", lat, long, alt, dist, hr, cad, speed, "power", activity_uuid
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, config), athlete.Name), db)

	return &a
}

func (r *AthleteRepository) Begin() (*sql.Tx, error) {
	return r.Db.Begin()
}

func (r *AthleteRepository) HasImported(file string) (bool, error) {
	res, err := r.hasImportedQuery.Query(file)
	defer res.Close()

	if err != nil {
		return false, err
	}
	exists := res.Next()
	return exists, nil
}

func (r *AthleteRepository) RecordImport(file string) error {
	_, err := r.recordImportQuery.Exec(time.Now(), file)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddActivity(act *models.Activity) error {
	_, err := r.addActivityQuery.Exec(
		act.FullUUID, act.UUID, act.Sport, act.ID, act.Creator.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddLap(act *models.Activity, lap *models.Lap) error {
	_, err := r.addLapQuery.Exec(
		lap.Time, lap.Start, lap.TotalTime, lap.Dist, lap.Calories, lap.MaxSpeed,
		lap.AvgHr, lap.MaxHr, lap.Intensity, lap.TriggerMethod, act.FullUUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddTrackPoint(act *models.Activity, tp *models.Trackpoint) error {
	_, err := r.addTrackPointQuery.Exec(
		tp.Time, tp.Lat, tp.Long, tp.Alt, tp.Dist,
		tp.HR, tp.Cad, tp.Speed, tp.Power, act.FullUUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) GetActivities(start time.Time, end time.Time) ([]*models.Activity, error) {
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
	return nil, nil
}
