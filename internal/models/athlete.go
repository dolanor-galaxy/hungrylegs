package models

import (
	"encoding/base64"
)

// Athlete is the root level account
type Athlete struct {
	Alterego     *string
	Name         *string
	FileSafeName *string
	// from the db
	Age         int8
	Weight      int8
	Vo2Max      int8
	RestHR      int8
	RunMaxHR    int8
	RunFTPace   float64
	RunZone5    int8
	RunZone4    int8
	RunZone3    int8
	RunZone2    int8
	RunZone1    int8
	BikeMaxHR   int8
	BikeFTPower float64
	BikeZone5   int8
	BikeZone4   int8
	BikeZone3   int8
	BikeZone2   int8
	BikeZone1   int8
	SwimMaxHR   int8
	SwimFTPace  float64
	SwimZone5   int8
	SwimZone4   int8
	SwimZone3   int8
	SwimZone2   int8
	SwimZone1   int8
	FtpWatts    float64
	WattsZone7  int8
	WattsZone6  int8
	WattsZone5  int8
	WattsZone4  int8
	WattsZone3  int8
	WattsZone2  int8
	WattsZone1  int8
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
