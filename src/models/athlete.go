package models

import (
	"encoding/base64"
)

// Athlete is the root level account
type Athlete struct {
	Name         string
	FileSafeName string
}

func NewAthlete(name string) *Athlete {
	a := Athlete{
		Name:         name,
		FileSafeName: base64.URLEncoding.EncodeToString([]byte(name)),
	}
	return &a
}
