package server

import (
	"context"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

type mutationResolver struct{ *Resolver }

func (r *mutationResolver) CreateAthlete(ctx context.Context, input NewAthlete) (*Athlete, error) {
	panic("not implemented")
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Athletes(ctx context.Context) ([]*Athlete, error) {
	// panic("not implemented")
	ls := make([]*Athlete, 1)
	athlete := Athlete {
		Name: "Rob Rohan",
		Alterego: "Professor Zoom",
	}
	ls[0] = &athlete

	return ls, nil
}
