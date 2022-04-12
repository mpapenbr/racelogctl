package internal

// this holds the resolved configuration values
var (
	Url                  string // url of WAMP server
	Realm                string // realm to use
	EventId              int    // used for actions against a single event
	Num                  int    // used to hold the number of items to fetch (for example when retrieving states)
	From                 int    // used to hold the from timestamp when fetching states
	Output               string // used to hold the output filename
	Input                string // used to hold the input filename (when importing data)
	FullStateData        bool   // if true all states for an event should be fetched
	OutputFormat         string // output format to be used (text,json)
	JsonPretty           bool   // prettify json output
	SkipPersistence      bool   // if true the backend will not persist any data (useful for replay)
	SampleFile           string // file name of sample for specific action
	EventName            string // event name for registration
	EventKey             string // event key for registration
	EventDescription     string // event description for registration
	AdminPassword        string // the password used to perform admin commands
	DataproviderPassword string // the password used to perform dataprovider commands

)
