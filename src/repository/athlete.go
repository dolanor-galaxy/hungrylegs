package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/therohans/HungryLegs/src/models"
)

type AthleteRepository struct {
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

// Attach creates a new repository and sets up needed bits
func Attach(athlete *models.Athlete) *AthleteRepository {
	db := athlete.Db
	a := AthleteRepository{
		Db: db,
	}

	a.hasImportedQuery = prepareQuery(`
		SELECT id FROM FileImport WHERE file_name = ?
	`, db)

	a.recordImportQuery = prepareQuery(`
		INSERT INTO FileImport (
			import_time, 'file_name'
		) VALUES (?, ?)
	`, db)

	a.addActivityQuery = prepareQuery(`
		INSERT INTO Activity (
			uuid, full_uuid, sport, 'time', device
		) VALUES (?, ?, ?, ?, ?)
	`, db)

	a.addLapQuery = prepareQuery(`
		INSERT INTO Lap (
			'time', 'start', total_time, dist, calories, max_speed, 
			avg_hr, max_hr, intensity, trigger, activity_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, db)

	a.addTrackPointQuery = prepareQuery(`
		INSERT INTO TrackPoint (
			'time', lat, long, alt, dist, hr, cad, speed, 'power', activity_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, db)

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

func (r *AthleteRepository) AddActivity(act *models.Activity) (int64, error) {
	res, err := r.addActivityQuery.Exec(
		act.UUID, act.FullUUID, act.Sport, act.ID, act.Creator.Name)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *AthleteRepository) AddLap(activityID int64, lap *models.Lap) (int64, error) {
	res, err := r.addLapQuery.Exec(
		lap.Time, lap.Start, lap.TotalTime, lap.Dist, lap.Calories, lap.MaxSpeed,
		lap.AvgHr, lap.MaxHr, lap.Intensity, lap.TriggerMethod, activityID,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r *AthleteRepository) AddTrackPoint(activityID int64, tp *models.Trackpoint) (int64, error) {
	res, err := r.addTrackPointQuery.Exec(
		tp.Time, tp.Lat, tp.Long, tp.Alt, tp.Dist,
		tp.HR, tp.Cad, tp.Speed, tp.Power, activityID,
	)
	if err != nil {
		return -1, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
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
