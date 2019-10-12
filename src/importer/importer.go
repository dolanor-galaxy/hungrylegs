package importer

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
	"github.com/therohans/HungryLegs/src/tcx"
	"github.com/tormoder/fit"
)

type Importer interface {
	Import(file string, repo repository.AthleteRepository) error
}

func ImportNewActivity(config *models.StaticConfig, repo *repository.AthleteRepository) {
	log.Println("Beginning import of new files...")
	files, err := ioutil.ReadDir(config.ImportDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()
		name = strings.ToLower(name)

		have, err := repo.HasImported(name)
		if err != nil {
			log.Fatal(err)
		}

		if have == false {
			start := time.Now()
			if strings.HasSuffix(name, ".tcx") {
				tcxFile := TcxFile{}
				err := tcxFile.Import(filepath.Join(config.ImportDir, name), repo)
				if err != nil {
					log.Fatal(err)
				}
			} else if strings.HasSuffix(name, ".fit") {
				fitFile := FitFile{}
				err := fitFile.Import(filepath.Join(config.ImportDir, name), repo)
				if err != nil {
					log.Fatal(err)
				}
			}
			repo.RecordImport(name)

			t := time.Now()
			elapsed := t.Sub(start)
			log.Printf("%v took %v", name, elapsed)
		} else {
			log.Printf("Already imported %v\n", name)
		}
	}
	log.Println("Done import")
}

////////////////////////////////////

// FitFile represents a .fit file (standard garmin)
type FitFile struct{}

func (f *FitFile) Import(file string, repo *repository.AthleteRepository) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// Decode the FIT file data
	fitFile, err := fit.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	if fitFile.Type() == fit.FileTypeActivity {
		tx, err := repo.Begin()
		if err != nil {
			return err
		}

		activity, err := fitFile.Activity()
		if err != nil {
			return err
		}

		var activityID int64
		for _, session := range activity.Sessions {
			hlAct := models.Activity{
				ID:    session.Timestamp,
				Sport: session.Sport.String(),
			}
			activityID, err = repo.AddActivity(&hlAct)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, lap := range activity.Laps {
			hlLap := models.Lap{
				Start:         lap.StartTime.String(),
				TotalTime:     float64(lap.TotalElapsedTime / 1000),
				Dist:          float64(lap.TotalDistance) / 100000,
				Calories:      float64(lap.TotalCalories),
				MaxSpeed:      float64(lap.MaxSpeed),
				AvgHr:         float64(lap.AvgHeartRate),
				MaxHr:         float64(lap.MaxHeartRate),
				Intensity:     lap.Intensity.String(),
				TriggerMethod: lap.LapTrigger.String(),
			}
			lapID, err := repo.AddLap(activityID, &hlLap)
			if err != nil {
				tx.Rollback()
				return err
			}

			for _, track := range activity.Records {
				htTrack := models.Trackpoint{
					Time:  track.Timestamp,
					Lat:   track.PositionLat.Degrees(),
					Long:  track.PositionLong.Degrees(),
					Alt:   float64(track.Altitude),
					Dist:  float64(track.Distance),
					HR:    float64(track.HeartRate),
					Cad:   float64(track.Cadence),
					Speed: float64(track.Speed),
					Power: float64(track.Power),
				}
				_, err := repo.AddTrackPoint(lapID, &htTrack)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		tx.Commit()
	}

	return nil
}

////////////////////////////////////

// TcxFile represents an older .tcx file (old garmin)
type TcxFile struct{}

func (f *TcxFile) Import(file string, repo *repository.AthleteRepository) error {
	tcxdb, err := tcx.ReadFile(file)
	if err != nil {
		return err
	}

	for i := range tcxdb.Acts.Act {
		tx, err := repo.Begin()
		if err != nil {
			return err
		}

		act := tcxdb.Acts.Act[i]
		hlAct := models.Activity{
			ID:    act.Id,
			Sport: act.Sport,
		}
		activityID, err := repo.AddActivity(&hlAct)
		if err != nil {
			tx.Rollback()
			return err
		}

		for l := range act.Laps {
			lap := act.Laps[l]
			hlLap := models.Lap{
				Start:         lap.Start,
				TotalTime:     lap.TotalTime,
				Dist:          lap.Dist,
				Calories:      lap.Calories,
				MaxSpeed:      lap.MaxSpeed,
				AvgHr:         lap.AvgHr,
				MaxHr:         lap.MaxHr,
				Intensity:     lap.Intensity,
				TriggerMethod: lap.TriggerMethod,
			}
			lapID, err := repo.AddLap(activityID, &hlLap)
			if err != nil {
				tx.Rollback()
				return err
			}

			for t := range lap.Trk.Pt {
				track := lap.Trk.Pt[t]
				htTrack := models.Trackpoint{
					Time:  track.Time,
					Lat:   track.Lat,
					Long:  track.Long,
					Alt:   track.Alt,
					Dist:  track.Dist,
					HR:    track.HR,
					Cad:   track.Cad,
					Speed: track.Speed,
					Power: track.Power,
				}
				_, err := repo.AddTrackPoint(lapID, &htTrack)
				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		tx.Commit()
	}

	return nil
}
