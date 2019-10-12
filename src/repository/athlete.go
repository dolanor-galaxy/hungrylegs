package repository

import (
	"database/sql"
	"time"

	"github.com/therohans/HungryLegs/src/models"
)

type AthleteRepository struct {
	Db *sql.DB
}

func (r *AthleteRepository) Begin() (*sql.Tx, error) {
	return r.Db.Begin()
}

func (r *AthleteRepository) HasImported(file string) (bool, error) {
	statement, err := r.Db.Prepare(`
		SELECT id FROM FileImport WHERE file_name = ?
	`)
	if err != nil {
		return false, err
	}
	res, err := statement.Query(file)
	defer res.Close()

	if err != nil {
		return false, err
	}
	exists := res.Next()
	return exists, nil
}

func (r *AthleteRepository) RecordImport(file string) error {
	statement, err := r.Db.Prepare(`
		INSERT INTO FileImport (
			import_time, 'file_name'
		) VALUES (?, ?)
	`)
	defer statement.Close()
	if err != nil {
		return err
	}
	_, err = statement.Exec(time.Now(), file)
	if err != nil {
		return err
	}
	return nil
}

func (r *AthleteRepository) AddActivity(act *models.Activity) (int64, error) {
	statement, err := r.Db.Prepare(`
		INSERT INTO Activity (
			sport, 'time', device
		) VALUES (?, ?, ?)
	`)
	defer statement.Close()
	if err != nil {
		return -1, err
	}
	res, err := statement.Exec(act.Sport, act.ID, act.Creator.Name)
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
	statement, err := r.Db.Prepare(`
		INSERT INTO Lap (
			'start', total_time, dist, calories, max_speed, 
			avg_hr, max_hr, intensity, trigger, activity_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	defer statement.Close()
	if err != nil {
		return -1, err
	}
	res, err := statement.Exec(
		lap.Start, lap.TotalTime, lap.Dist, lap.Calories, lap.MaxSpeed,
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

func (r *AthleteRepository) AddTrackPoint(lapID int64, tp *models.Trackpoint) (int64, error) {
	statement, err := r.Db.Prepare(`
		INSERT INTO TrackPoint (
			'time', lat, long, alt, dist, hr, cad, speed, 'power', lap_id
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	defer statement.Close()
	if err != nil {
		return -1, err
	}
	res, err := statement.Exec(
		tp.Time, tp.Lat, tp.Long, tp.Alt, tp.Dist,
		tp.HR, tp.Cad, tp.Speed, tp.Power, lapID,
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
