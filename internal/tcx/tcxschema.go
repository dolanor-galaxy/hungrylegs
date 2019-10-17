package tcx

import (
	"encoding/xml"
	"time"
)

type TCXDB struct {
	XMLName xml.Name    `xml:"TrainingCenterDatabase"`
	Acts    *Activities `xml:"Activities"`
	Auth    Author      `xml:"Author"`
}

type Trackpoint struct {
	Time  time.Time
	Lat   float64 `xml:"Position>LatitudeDegrees"`
	Long  float64 `xml:"Position>LongitudeDegrees"`
	Alt   float64 `xml:"AltitudeMeters,omitempty"`
	Dist  float64 `xml:"DistanceMeters,omitempty"`
	HR    float64 `xml:"HeartRateBpm>Value,omitempty"`
	Cad   float64 `xml:"Cadence,omitempty"`
	Speed float64 `xml:"Extensions>TPX>Speed,omitempty"`
	Power float64 `xml:"Extensions>TPX>Watts,omitempty"`
}

type Track struct {
	Pt []Trackpoint `xml:"Trackpoint"`
}

type Lap struct {
	Start         string  `xml:"StartTime,attr"`
	TotalTime     float64 `xml:"TotalTimeSeconds,omitempty"`
	Dist          float64 `xml:"DistanceMeters,omitempty"`
	Calories      float64 `xml:",omitempty"`
	MaxSpeed      float64 `xml:"MaximumSpeed,omitempty"`
	AvgHr         float64 `xml:"AverageHeartRateBpm>Value,omitempty"`
	MaxHr         float64 `xml:"MaximumHeartRateBpm>Value,omitempty"`
	Intensity     string  `xml:",omitempty"`
	TriggerMethod string  `xml:",omitempty"`
	Trk           *Track  `xml:"Track"`
}

type Activity struct {
	Sport   string `xml:"Sport,attr,omitempty"`
	Id      time.Time
	Laps    []Lap  `xml:"Lap,omitempty"`
	Creator Device `xml:"Creator,omitempty"`
}

type Activities struct {
	Act []Activity `xml:"Activity"`
}

type Author struct {
	Name       string `xml:",omitempty"`
	Build      Build  `xml:",omitempty"`
	LangID     string `xml:",omitempty"`
	PartNumber string `xml:",omitempty"`
}

// Device is the tech used to capture this data - e.g. Fenix 3
type Device struct {
	Name      string       `xml:",omitempty"`
	UnitId    uint         `xml:",omitempty"`
	ProductID string       `xml:",omitempty"`
	Version   BuildVersion `xml:",omitempty"`
}

type Build struct {
	Version BuildVersion `xml:"Version,omitempty"`
	Type    string       `xml:",omitempty"`
	Time    string       `xml:",omitempty"`
	Builder string       `xml:",omitempty"`
}

type BuildVersion struct {
	VersionMajor int `xml:",omitempty"`
	VersionMinor int `xml:",omitempty"`
	BuildMajor   int `xml:",omitempty"`
	BuildMinor   int `xml:",omitempty"`
}
