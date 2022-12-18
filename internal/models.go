package internal

type EventInfo struct {
	RaceloggerVersion string `json:"raceloggerVersion"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	EventTime         string `json:"eventTime"`
	// trackInfo of register is merged into EventInfo upon registration
	TrackId               int     `json:"trackId"`
	TrackDisplayName      string  `json:"trackDisplayName"`
	TrackDisplayShortName string  `json:"trackDisplayShortName"`
	TrackConfigName       string  `json:"trackConfigName"`
	TrackLength           float64 `json:"trackLength"`
	TrackPitSpeed         float64 `json:"trackPitSpeed"`
	MultiClass            bool    `json:"multiClass"`
	TeamRacing            int     `json:"teamRacing"`
	IrSessionId           int     `json:"irSessionId"`
	NumCarTypes           int     `json:"numCarTypes"`
	NumCarClasses         int     `json:"numCarClasses"`
	SpeedmapInterval      int     `json:"speedmapInterval"`

	Sectors  []Sector       `json:"sectors"`
	Sessions []EventSession `json:"sessions"`
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
	Pit                   struct {
		Entry      float64 `json:"entry"`
		Exit       float64 `json:"exit"`
		LaneLength float64 `json:"laneLength"`
	} `json:"pit"`
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
	EventKey    string `json:"eventKey"`
	Name        string `json:"name"`
	Description string `json:"description"`
	RecordDate  string `json:"recordDate"`
	Data        Data   `json:"data"`

	Id int32 `json:"id"`
}

type Payload struct {
	Cars     [][]interface{} `json:"cars"`
	Session  []interface{}   `json:"session"`
	Messages [][]interface{} `json:"messages"`
}

// TODO: rename to WampMessage?
type State struct {
	Type      int     `json:"type"` // 1: full data, 2: delta data
	Payload   Payload `json:"payload"`
	Timestamp float64 `json:"timestamp"`
}

type RegisterMessage struct {
	EventKey   string    `json:"eventKey"`
	Manifests  Manifests `json:"manifests"`
	Info       EventInfo `json:"info"`
	TrackInfo  TrackInfo `json:"trackInfo"`
	RecordDate float64   `json:"recordDate"` // is a timestamp
}

type ProviderData struct {
	EventKey   string     `json:"eventKey"`
	Manifests  Manifests  `json:"manifests"`
	Info       EventInfo  `json:"info"`
	ReplayInfo ReplayInfo `json:"replayInfo"`
	DbId       int        `json:"dbId"`
}

// common return type for rpc calls.
// if no error occurs, Error is not set/empty.
// In these cases depending on the call Message and/or Data may be set
type ResultMessage struct {
	Error   string        `json:"error"`
	Message string        `json:"message"`
	Data    []interface{} `json:"data"`
}

type AverageLapTime struct {
	Timestamp   float64         `json:"timestamp"`
	SessionTime float64         `json:"sessionTime"`
	TimeOfDay   float64         `json:"timeOfDay"`
	TrackTemp   float64         `json:"trackTemp"`
	Laptimes    map[int]float64 `json:"laptimes"`
}
type EventCarMessage struct {
	Type      int       `json:"type"`
	Timestamp float64   `json:"timestamp"`
	Payload   EventCars `json:"payload"`
}
type EventCars struct {
	Cars       []Car        `json:"cars"`
	CarClasses []CarClass   `json:"carClasses"`
	Entries    []EventEntry `json:"entries"`
}
type CarClass struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Car struct {
	Name          string  `json:"name"`
	NameShort     string  `json:"nameShort"`
	CarClassName  string  `json:"carClassName"`
	CarId         int     `json:"carId"`
	CarClassId    int     `json:"carClassId"`
	FuelPct       float64 `json:"fuelPct"`
	PowerAdjust   float64 `json:"powerAdjust"`
	WeightPenalty float64 `json:"weightPenalty"`
	DryTireSets   int     `json:"dryTireSets"`
}

type EventEntry struct {
	Car struct {
		Name         string `json:"name"`
		CarId        int    `json:"carId"`
		CarIdx       int    `json:"carIdx"`
		CarClassId   int    `json:"carClassId"`
		CarNumber    string `json:"carNumber"`
		CarNumberRaw int    `json:"carNumberRaw"`
	} `json:"car"`
	Team struct {
		Id     int    `json:"id"`
		Name   string `json:"name"`
		CarIdx int    `json:"carIdx"`
	} `json:"team"`
	Drivers []struct {
		Id          int    `json:"id"`
		Name        string `json:"name"`
		CarIdx      int    `json:"carIdx"`
		IRating     int    `json:"iRating"`
		Initials    string `json:"initials"`
		LicLevel    int    `json:"licLevel"`
		LicString   string `json:"licString"`
		LicSubLevel int    `json:"licSubLevel"`
		AbbrevName  string `json:"abbrevName"`
	} `json:"drivers"`
}

type SpeedmapMessage struct {
	Type      int             `json:"type"`
	Timestamp float64         `json:"timestamp"`
	Payload   SpeedmapPayload `json:"payload"`
}
type SpeedmapPayload struct {
	Data map[string]struct {
		ChunkSpeeds []float64 `json:"chunkSpeeds"`
		Laptime     float64   `json:"laptime"`
	} `json:"data"`
	ChunkSize   int     `json:"chunkSize"`
	CurrentPos  float64 `json:"currentPos"`
	TrackLength float64 `json:"trackLength"`
	TrackTemp   float64 `json:"trackTemp"`
	SessionTime float64 `json:"sessionTime"`
	TimeOfDay   float64 `json:"timeOfDay"`
}
