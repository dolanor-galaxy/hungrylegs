package server

import (
	"context"

	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type activitiesResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAthlete(ctx context.Context, input NewAthlete) (*Athlete, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Athlete(ctx context.Context, alterego string) (*Athlete, error) {
	// db := DBFromContext(ctx)
	// config := ConfigFromContext(ctx)

	athlete := models.NewAthlete(&alterego, &alterego)
	a := Athlete{
		Alterego: *athlete.Alterego,
		Name:     *athlete.Name,
	}
	// repo := repository.Attach(athlete.Name, db, config)

	s := "1900-01-01"
	e := "3000-01-01"
	acts, _ := r.Activities(ctx, a.Alterego, &s, &e)
	a.Activities = acts

	return &a, nil
}

func (r *queryResolver) Activities(ctx context.Context, athleteID string, startTime *string, endTime *string) ([]*Activity, error) {
	db := DBFromContext(ctx)
	config := ConfigFromContext(ctx)
	repo := repository.Attach(&athleteID, db, config)

	macts, err := repo.GetActivities(*startTime, *endTime)
	if err != nil {
		return nil, err
	}

	var acts []*Activity
	for _, ac := range macts {
		a := Activity{}
		a.ID = ac.UUID
		a.Sid = ac.SUUID
		a.Sport = ac.Sport
		a.Time = ac.Time

		a.Laps, _ = r.Laps(ctx, athleteID, ac.UUID)
		a.Trackpoints, _ = r.Trackpoints(ctx, athleteID, ac.UUID)

		acts = append(acts, &a)
	}

	return acts, nil
}
func (r *queryResolver) Laps(ctx context.Context, athleteID string, activityID string) ([]*Lap, error) {
	db := DBFromContext(ctx)
	config := ConfigFromContext(ctx)
	repo := repository.Attach(&athleteID, db, config)

	mlaps, err := repo.GetLaps(activityID)
	if err != nil {
		return nil, err
	}

	var laps []*Lap
	for _, ac := range mlaps {
		a := Lap{}
		a.Time = ac.Time
		a.Duration = ac.TotalTime
		a.Distance = ac.Dist
		a.Calories = ac.Calories
		a.MaxSpeed = ac.MaxSpeed
		a.AvgHr = ac.AvgHr
		a.MaxHr = ac.MaxHr
		a.Intensity = ac.Intensity
		laps = append(laps, &a)
	}

	return laps, nil
}
func (r *queryResolver) Trackpoints(ctx context.Context, athleteID string, activityID string) ([]*TrackPoint, error) {
	db := DBFromContext(ctx)
	config := ConfigFromContext(ctx)
	repo := repository.Attach(&athleteID, db, config)

	mlaps, err := repo.GetTrackpoints(activityID)
	if err != nil {
		return nil, err
	}

	var tps []*TrackPoint
	for _, ac := range mlaps {
		a := TrackPoint{}
		a.Time = ac.Time
		a.Lat = ac.Lat
		a.Long = ac.Long
		a.Altitude = ac.Alt
		a.Distance = ac.Dist
		a.Hr = ac.HR
		a.Cadence = ac.Cad
		a.Speed = ac.Speed
		a.Power = ac.Power
		tps = append(tps, &a)
	}

	return tps, nil
}
