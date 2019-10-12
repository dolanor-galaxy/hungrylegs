package models

import (
	"time"
)

// Trackpoint is a single capture event
type Trackpoint struct {
	Time  time.Time
	Lat   float64
	Long  float64
	Alt   float64
	Dist  float64
	HR    float64
	Cad   float64
	Speed float64
	Power float64
}

// Track is just a group of Trackpoints - a joining table
type Track struct {
	Pt []Trackpoint
}

// Lap is a single cycle, say 1km or when you click the lap button
type Lap struct {
	Start         string
	TotalTime     float64
	Dist          float64
	Calories      float64
	MaxSpeed      float64
	AvgHr         float64
	MaxHr         float64
	Intensity     string
	TriggerMethod string
	Trk           *Track
}

// Activity is a single session of a particular activity Bike, Run, Swim, etc
type Activity struct {
	Sport   string
	ID      time.Time
	Laps    []Lap
	Creator Device
}

// Activities is the list of activities for this session
type Activities struct {
	Act []Activity
}

//////////////////////////////////////////////////////////////////////

type Author struct {
	Name       string
	Build      Build
	LangID     string
	PartNumber string
}

//////////////////////////////////////////////////////////////////////

// Device is the tech used to capture this data - e.g. Fenix 3
type Device struct {
	Name      string
	UnitId    int
	ProductID string
	Version   BuildVersion
}

type Build struct {
	Version BuildVersion
	Type    string
	Time    string
	Builder string
}

type BuildVersion struct {
	VersionMajor int
	VersionMinor int
	BuildMajor   int
	BuildMinor   int
}
