package server

import (
	"context"
	"log"

	"github.com/therohans/HungryLegs/internal/models"
	"github.com/therohans/HungryLegs/internal/repository"
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

	// s := "1900-01-01"
	// e := "3000-01-01"
	// acts, _ := r.Activities(ctx, *athlete.Alterego, &s, &e)
	// a.Activities = acts

	return &a, nil
}

func (r *queryResolver) Activities(ctx context.Context, athleteID string, startTime *string, endTime *string) ([]*Activity, error) {
	db := DBFromContext(ctx)
	config := ConfigFromContext(ctx)
	repo := repository.Attach(&athleteID, db, config)

	log.Printf("out: %v %v\n", startTime, endTime)

	macts, err := repo.GetActivities(*startTime, *endTime)
	if err != nil {
		return nil, err
	}

	var acts []*Activity
	for _, ac := range macts {
		a := Activity{}
		a.ID = ac.FullUUID
		a.Sid = ac.UUID
		a.Sport = ac.Sport
		a.Time = ac.Time
		acts = append(acts, &a)
	}

	return acts, nil
}
func (r *queryResolver) Laps(ctx context.Context, athleteID string, activityID *string, startTime *string, endTime *string) ([]*Lap, error) {
	log.Printf("%v", ctx)
	log.Printf("%v", r)
	return nil, nil
}
func (r *queryResolver) Trackpoints(ctx context.Context, athleteID string, activityID string, startTime *string, endTime *string) ([]*TrackPoint, error) {
	log.Printf("%v", ctx)
	log.Printf("%v", r)
	return nil, nil
}
