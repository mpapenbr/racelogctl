package internal

// this holds the resolved configuration values
var (
	Url     string // url of WAMP server
	Realm   string // realm to use
	EventId int    // used for actions against a single event
	Num     int    // used to hold the number of items to fetch (for example when retrieving states)
	From    int    // used to hold the from timestamp when fetching states
)
