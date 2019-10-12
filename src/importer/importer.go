package importer

import (
	"bytes"
	"fmt"
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

func (f *FitFile) Import(file string, repo repository.AthleteRepository) error {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Decode the FIT file data
	_, err = fit.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	// fmt.Println(fit.Type())
	return nil
}

////////////////////////////////////

type TcxFile struct{}

func (f *TcxFile) Import(file string, repo repository.AthleteRepository) error {
	fmt.Println(file)
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
