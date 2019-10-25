package repository

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/robrohan/HungryLegs/internal/models"
)

type AthleteRepository struct {
	Athlete            *models.Athlete
	Db                 *sql.DB
	hasImportedQuery   *sql.Stmt
	recordImportQuery  *sql.Stmt
	addActivityQuery   *sql.Stmt
	addLapQuery        *sql.Stmt
	addTrackPointQuery *sql.Stmt
	getActivities      *sql.Stmt
	getLaps            *sql.Stmt
	getTrackpoints     *sql.Stmt
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
func Attach(schema string, db *sql.DB, config *models.StaticConfig) *AthleteRepository {
	a := AthleteRepository{
		Db: db,
	}

	a.hasImportedQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		SELECT import_time FROM "%v".fileimport WHERE file_name = $1
	`, config), schema), db)

	a.recordImportQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".fileimport (
			import_time, "file_name"
		) VALUES ($1, $2)
	`, config), schema), db)

	a.addActivityQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".activity (
			uuid, suuid, sport, "time", device
		) VALUES ($1, $2, $3, $4, $5)
	`, config), schema), db)

	a.addLapQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".lap (
			"time", "start", total_time, dist, calories, max_speed,
			avg_hr, max_hr, intensity, trigger, activity_uuid
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`, config), schema), db)

	a.addTrackPointQuery = prepareQuery(fmt.Sprintf(sqlForDriver(`
		INSERT INTO "%v".trackpoint (
			"time", lat, long, alt, dist, hr, cad, speed, "power", activity_uuid
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, config), schema), db)

	a.getActivities = prepareQuery(fmt.Sprintf(sqlForDriver(`
		SELECT uuid, suuid, sport, "time", device
		FROM "%v".activity
		WHERE ("time" >= $1 AND "time" <= $2)
		ORDER BY "time" desc
		LIMIT 100
`, config), schema), db)

	a.getLaps = prepareQuery(fmt.Sprintf(sqlForDriver(`
		SELECT
			"time", start, total_time, dist, calories,
			max_speed, avg_hr, max_hr, intensity
		FROM "%v".lap
		WHERE activity_uuid = $1
		ORDER BY "time" asc
	`, config), schema), db)

	a.getTrackpoints = prepareQuery(fmt.Sprintf(sqlForDriver(`
		SELECT 
			"time", lat, long, alt, dist, hr,
			cad, speed, power
		FROM "%v".trackpoint
		WHERE activity_uuid = $1
		ORDER BY "time" asc
	`, config), schema), db)

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
		act.UUID, act.SUUID, act.Sport, act.Time, act.Creator.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddLap(act *models.Activity, lap *models.Lap) error {
	_, err := r.addLapQuery.Exec(
		lap.Time, lap.Start, lap.TotalTime, lap.Dist, lap.Calories, lap.MaxSpeed,
		lap.AvgHr, lap.MaxHr, lap.Intensity, lap.TriggerMethod, act.UUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddTrackPoint(act *models.Activity, tp *models.Trackpoint) error {
	_, err := r.addTrackPointQuery.Exec(
		tp.Time, tp.Lat, tp.Long, tp.Alt, tp.Dist,
		tp.HR, tp.Cad, tp.Speed, tp.Power, act.UUID,
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) GetActivities(start string, end string) ([]*models.Activity, error) {
	rows, err := r.getActivities.Query(start, end)
	if err != nil {
		return nil, err
	}
	var acts []*models.Activity

	if rows != nil {
		for rows.Next() {
			a := models.Activity{}
			rows.Scan(
				&a.UUID,
				&a.SUUID,
				&a.Sport,
				&a.Time,
				&a.Creator)
			acts = append(acts, &a)
		}
	}
	return acts, nil
}

func (r *AthleteRepository) GetLaps(activity_uuid string) ([]*models.Lap, error) {
	// log.Printf("%v\n", r.getLaps)
	rows, err := r.getLaps.Query(activity_uuid)
	if err != nil {
		return nil, err
	}
	var laps []*models.Lap

	if rows != nil {
		for rows.Next() {
			l := models.Lap{}
			rows.Scan(
				&l.Time,
				&l.Start,
				&l.TotalTime,
				&l.Dist,
				&l.Calories,
				&l.MaxSpeed,
				&l.AvgHr,
				&l.MaxHr,
				&l.Intensity)
			laps = append(laps, &l)
		}
	}
	return laps, nil
}

func (r *AthleteRepository) GetTrackpoints(activity_uuid string) ([]*models.Trackpoint, error) {
	rows, err := r.getTrackpoints.Query(activity_uuid)
	if err != nil {
		return nil, err
	}
	var tp []*models.Trackpoint

	if rows != nil {
		for rows.Next() {
			l := models.Trackpoint{}
			rows.Scan(
				&l.Time,
				&l.Lat,
				&l.Long,
				&l.Alt,
				&l.Dist,
				&l.HR,
				&l.Cad,
				&l.Speed,
				&l.Power)
			tp = append(tp, &l)
		}
	}
	return tp, nil
}
