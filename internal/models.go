package internal

type EventInfo struct {
	RaceloggerVersion string `json:"raceloggerVersion"`
	Name              string `json:"name"`
	EventTime         string `json:"eventTime"`
	// trackInfo of register is merged into EventInfo upon registration
	TrackId               int      `json:"trackId"`
	TrackDisplayName      string   `json:"trackDisplayName"`
	TrackDisplayShortName string   `json:"trackDisplayShortName"`
	TrackConfigName       string   `json:"trackConfigName"`
	TrackLength           float64  `json:"trackLength"`
	Sectors               []Sector `json:"sectors"`
}

// when registering an event we need this separate struct
type TrackInfo struct {
	TrackId               int      `json:"trackId"`
	TrackDisplayName      string   `json:"trackDisplayName"`
	TrackDisplayShortName string   `json:"trackDisplayShortName"`
	TrackConfigName       string   `json:"trackConfigName"`
	TrackLength           float64  `json:"trackLength"`
	Sectors               []Sector `json:"sectors"`
}

type Sector struct {
	SectorNum      int
	SectorStartPct float64
}

type ReplayInfo struct {
	MinTimestamp   float64
	MinSessionTime float64
	MaxSessionTime float64
}
type Manifests struct {
	Car     []string
	Session []string
	Message []string
}
type Data struct {
	Info       EventInfo
	ReplayInfo ReplayInfo
	Manifests  Manifests
}

type Event struct {
	EventKey   string
	Name       string
	RecordDate string
	Data       Data

	Id int32
}

type Payload struct {
	Cars     [][]interface{} `json:"cars"`
	Session  []interface{}   `json:"session"`
	Messages [][]interface{} `json:"messages"`
}
type State struct {
	Type      int     `json:"type"` // 1: full data, 8: delta data
	Payload   Payload `json:"payload"`
	Timestamp float64 `json:"timestamp"`
}

type RegisterMessage struct {
	EventKey  string    `json:"eventKey"`
	Manifests Manifests `json:"manifests"`
	Info      EventInfo `json:"info"`
	TrackInfo TrackInfo `json:"trackInfo"`
}
