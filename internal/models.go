package internal

type EventInfo struct {
	RaceloggerVersion string `json:"raceloggerVersion"`
	Name              string `json:"name"`
	EventTime         string `json:"eventTime"`
	// trackInfo of register is merged into EventInfo upon registration
	TrackId               int            `json:"trackId"`
	TrackDisplayName      string         `json:"trackDisplayName"`
	TrackDisplayShortName string         `json:"trackDisplayShortName"`
	TrackConfigName       string         `json:"trackConfigName"`
	TrackLength           float64        `json:"trackLength"`
	Sectors               []Sector       `json:"sectors"`
	Sessions              []EventSession `json:"sessions"`
}

type EventSession struct {
	Num  int    `json:"num"`
	Laps int    `json:"laps"`
	Name string `json:"name"`
	Time int    `json:"time"`
	Type string `json:"type"`
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

// no lowercase json attributes here by design (no conversion from iRacing data)
type Sector struct {
	SectorNum      int
	SectorStartPct float64
}

type ReplayInfo struct {
	MinTimestamp   float64 `json:"minTimestamp"`
	MinSessionTime float64 `json:"minSessionTime"`
	MaxSessionTime float64 `json:"maxSessionTime"`
}

type Manifests struct {
	Car     []string `json:"car"`
	Session []string `json:"session"`
	Message []string `json:"message"`
	Pit     []string `json:"pit"`
}
type Data struct {
	Info       EventInfo  `json:"info"`
	ReplayInfo ReplayInfo `json:"replayInfo"`
	Manifests  Manifests  `json:"manifests"`
}

type Event struct {
	EventKey   string `json:"eventKey"`
	Name       string `json:"name"`
	RecordDate string `json:"recordDate"`
	Data       Data   `json:"data"`

	Id int32 `json:"id"`
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

// common return type for rpc calls.
// if no error occurs, Error is not set/empty.
// In these cases depending on the call Message and/or Data may be set
type ResultMessage struct {
	Error   string        `json:"error"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}
