package server

import (
	"context"
	"database/sql"

	"github.com/robrohan/HungryLegs/internal/models"
	"github.com/robrohan/HungryLegs/internal/repository"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct {
	DB     *sql.DB
	Driver string
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Activity() ActivityResolver {
	return &activityResolver{r}
}

func (r *Resolver) Laps() ActivityResolver {
	return &activityResolver{r}
}

func (r *Resolver) Trackpoints() ActivityResolver {
	return &activityResolver{r}
}

func (r *Resolver) Athlete() AthleteResolver {
	return &athleteResolver{r}
}

func (r *Resolver) Activities() AthleteResolver {
	return &athleteResolver{r}
}

type activitiesResolver struct{ *Resolver }

type mutationResolver struct{ *Resolver }

type activityResolver struct{ *Resolver }

type athleteResolver struct{ *Resolver }

type queryResolver struct{ *Resolver }

////////////////////////////////////////////////////////////////
// Common function

func GetLaps(repo *repository.AthleteRepository, athleteID *string, activityID *string) ([]*Lap, error) {
	// log.Printf("Lap: %v %v\n", *athleteID, *activityID)
	mlaps, err := repo.GetLaps(*activityID)
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

func GetTrackpoints(repo *repository.AthleteRepository, athleteID *string, activityID *string) ([]*TrackPoint, error) {
	// log.Printf("Track: %v %v\n", *athleteID, *activityID)
	mlaps, err := repo.GetTrackpoints(*activityID)
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

func GetActivities(repo *repository.AthleteRepository, athleteID string, startTime *string, endTime *string) ([]*Activity, error) {
	// log.Printf("Act: %v %v %v\n", athleteID, *startTime, *endTime)
	macts, err := repo.GetActivities(*startTime, *endTime)
	if err != nil {
		return nil, err
	}

	ath := Athlete{
		ID:       *repo.Athlete.Alterego,
		Name:     *repo.Athlete.Name,
		Alterego: *repo.Athlete.Alterego,
	}

	var acts []*Activity
	for _, ac := range macts {
		a := Activity{}
		a.ID = ac.UUID
		a.Sid = ac.SUUID
		a.Sport = ac.Sport
		a.Time = ac.Time
		a.Athlete = &ath
		acts = append(acts, &a)
	}

	return acts, nil
}

////////////////////////////////////////////////////////////////
// Activity Resolver

func (a *activityResolver) Laps(ctx context.Context, obj *Activity) ([]*Lap, error) {
	repo := repository.Attach(obj.Athlete.Alterego, a.DB, a.Driver)

	return GetLaps(repo, &obj.Athlete.Alterego, &obj.ID)
}

func (a *activityResolver) Trackpoints(ctx context.Context, obj *Activity) ([]*TrackPoint, error) {
	repo := repository.Attach(obj.Athlete.Alterego, a.DB, a.Driver)

	return GetTrackpoints(repo, &obj.Athlete.Alterego, &obj.ID)
}

////////////////////////////////////////////////////////////////
// Athlete Resolver

func (r *athleteResolver) Activities(ctx context.Context, obj *Athlete, startTime *string, endTime *string) ([]*Activity, error) {
	repo := repository.Attach(obj.Alterego, r.DB, r.Driver)

	ath := models.NewAthlete(&obj.Alterego, &obj.Alterego)
	repo.Athlete = ath

	return GetActivities(repo, obj.Alterego, startTime, endTime)
}

////////////////////////////////////////////////////////////////
// Mutation

func (r *mutationResolver) CreateAthlete(ctx context.Context, input NewAthlete) (*Athlete, error) {
	panic("not implemented")
}

////////////////////////////////////////////////////////////////
// Query Resolver

func (r *queryResolver) Athlete(ctx context.Context, alterego string) (*Athlete, error) {
	// db := DBFromContext(ctx)
	// config := ConfigFromContext(ctx)

	athlete := models.NewAthlete(&alterego, &alterego)
	a := Athlete{
		Alterego: *athlete.Alterego,
		Name:     *athlete.Name,
	}
	// repo := repository.Attach(athlete.Name, db, config)

	// s := "1900-01-01"
	// e := "3000-01-01"
	// acts, _ := GetActivities(repo, a.Alterego, &s, &e)
	// a.Activities = acts

	return &a, nil
}

func (r *queryResolver) Activities(ctx context.Context, athleteID string, startTime *string, endTime *string) ([]*Activity, error) {
	repo := repository.Attach(athleteID, r.DB, r.Driver)

	return GetActivities(repo, athleteID, startTime, endTime)
}

func (r *queryResolver) Laps(ctx context.Context, athleteID string, activityID string) ([]*Lap, error) {
	repo := repository.Attach(athleteID, r.DB, r.Driver)

	return GetLaps(repo, &athleteID, &activityID)
}

func (r *queryResolver) Trackpoints(ctx context.Context, athleteID string, activityID string) ([]*TrackPoint, error) {
	repo := repository.Attach(athleteID, r.DB, r.Driver)

	return GetTrackpoints(repo, &athleteID, &activityID)
}
