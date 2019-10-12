package repository

import (
	"database/sql"

	"github.com/therohans/HungryLegs/src/models"
)

type AthleteRepository struct {
	Db *sql.DB
}

func (r *AthleteRepository) AddActivity(act *models.Activity) (int64, error) {
	statement, err := r.Db.Prepare(`
		INSERT INTO Activity (
			sport, 'time', device
		) VALUES (?, ?, ?)
	`)
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
