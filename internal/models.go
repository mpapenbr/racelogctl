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

type Data struct {
	Info       EventInfo
	ReplayInfo ReplayInfo
}

type Event struct {
	EventKey   string
	Name       string
	RecordDate string
	Data       Data

	Id int32
}

type Payload struct {
	Cars    [][]interface{}
	Session []interface{}
}
type State struct {
	Type    int // 1: full data, 8: delta data
	Payload Payload
}
