package models

import (
	"encoding/base64"
)

// Athlete is the root level account
type Athlete struct {
	Alterego     *string
	Name         *string
	FileSafeName *string
}

func NewAthlete(name *string, alterego *string) *Athlete {
	safe := base64.URLEncoding.EncodeToString([]byte(*alterego))
	a := Athlete{
		Alterego:     alterego,
		Name:         name,
		FileSafeName: &safe,
	}
	return &a
}
