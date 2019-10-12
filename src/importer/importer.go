package importer

import (
	"bytes"
	"io/ioutil"

	"github.com/therohans/HungryLegs/src/models"
	"github.com/therohans/HungryLegs/src/repository"
	"github.com/therohans/HungryLegs/src/tcx"
	"github.com/tormoder/fit"
)

type Importer interface {
	Import(file string, repo repository.AthleteRepository) error
}

////////////////////////////////////

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
			// fmt.Printf("act: %v \n", hlAct)
			if err != nil {
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
			// fmt.Printf("lap: %v \n", hlLap)
			if err != nil {
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
				// fmt.Printf("track: %v \n", htTrack)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

////////////////////////////////////

type TcxFile struct{}

func (f *TcxFile) Import(file string, repo *repository.AthleteRepository) error {
	tcxdb, err := tcx.ReadFile(file)
	if err != nil {
		return err
	}

	for i := range tcxdb.Acts.Act {
		act := tcxdb.Acts.Act[i]
		hlAct := models.Activity{
			ID:    act.Id,
			Sport: act.Sport,
		}
		activityID, err := repo.AddActivity(&hlAct)
		if err != nil {
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
					return err
				}
			}
		}
	}

	return nil
}
