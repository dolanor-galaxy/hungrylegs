package importer

import (
	"bytes"
	"crypto/sha1"
	"fmt"
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

func ImportActivity(name string, directory string, repo *repository.AthleteRepository) error {
	err := importFile(name, directory, repo)
	if err != nil {
		log.Printf("Error importing file: %v : %v", name, err.Error())
		return err
	}
	return nil
}

func ImportActivites(directory string, repo *repository.AthleteRepository) error {
	log.Println("Beginning import of new files...")
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		name := f.Name()
		err = importFile(name, directory, repo)
		if err != nil {
			log.Printf("Error importing file: %v : %v", f, err.Error())
		}
	}
	log.Println("Done import")
	return nil
}

func importFile(name string, directory string, repo *repository.AthleteRepository) error {
	// Check if this file has already been imported
	have, err := repo.HasImported(name)
	if err != nil {
		log.Fatal(err)
	}

	if have == false {
		start := time.Now()
		lower := strings.ToLower(name)
		// We only support tcx and fit files
		if strings.HasSuffix(lower, ".tcx") {
			tcxFile := TcxFile{}
			err := tcxFile.Import(filepath.Join(directory, name), repo)
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(lower, ".fit") {
			fitFile := FitFile{}
			err := fitFile.Import(filepath.Join(directory, name), repo)
			if err != nil {
				return err
			}
		} else {
			return nil
		}
		repo.RecordImport(name)

		t := time.Now()
		elapsed := t.Sub(start)
		log.Printf("%v took %v", name, elapsed)
	} else {
		log.Printf("Already imported %v\n", name)
	}
	return nil
}

// ActivityHash makes a mostly unique hash
func ActivityHash(sport string, time time.Time, file string) string {
	s := sport + "::" + string(time.Unix()) + "::" + file + "::hungrylegs"
	h := sha1.New()
	h.Write([]byte(s))
	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

////////////////////////////////////

// FitFile represents a .fit file (standard Garmin)
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
			sport := session.Sport.String()
			hash := ActivityHash(sport, session.Timestamp, file)
			hlAct := models.Activity{
				ID:       session.Timestamp,
				UUID:     hash[:8],
				FullUUID: hash,
				Sport:    sport,
			}
			activityID, err = repo.AddActivity(&hlAct)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, lap := range activity.Laps {
			hlLap := models.Lap{
				Time:          lap.Timestamp,
				Start:         lap.StartTime.String(),
				TotalTime:     lap.GetTotalElapsedTimeScaled(),
				Dist:          lap.GetTotalDistanceScaled(),
				Calories:      float64(lap.TotalCalories),
				MaxSpeed:      lap.GetMaxSpeedScaled(),
				AvgHr:         float64(lap.AvgHeartRate),
				MaxHr:         float64(lap.MaxHeartRate),
				Intensity:     lap.Intensity.String(),
				TriggerMethod: lap.LapTrigger.String(),
			}
			_, err := repo.AddLap(activityID, &hlLap)
			if err != nil {
				tx.Rollback()
				return err
			}
		}

		for _, track := range activity.Records {
			htTrack := models.Trackpoint{
				Time:  track.Timestamp,
				Lat:   track.PositionLat.Degrees(),
				Long:  track.PositionLong.Degrees(),
				Alt:   track.GetAltitudeScaled(),
				Dist:  track.GetDistanceScaled(),
				HR:    float64(track.HeartRate),
				Cad:   track.GetCadence256Scaled(),
				Speed: track.GetSpeedScaled(),
				Power: float64(track.Power),
			}
			_, err := repo.AddTrackPoint(activityID, &htTrack)
			if err != nil {
				tx.Rollback()
				return err
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
		hash := ActivityHash(act.Sport, act.Id, file)
		hlAct := models.Activity{
			ID:       act.Id,
			UUID:     hash[:8],
			FullUUID: hash,
			Sport:    act.Sport,
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
			_, err := repo.AddLap(activityID, &hlLap)
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
				_, err := repo.AddTrackPoint(activityID, &htTrack)
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
