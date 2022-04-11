package internal

type EventInfo struct {
	RaceloggerVersion string
	Name              string
	EventTime         string
	TrackDisplayName  string
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
