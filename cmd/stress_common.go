package cmd

var sourceEventId int = -1 // the eventId for the source
var recordingSpeed int = 1 // used to simulate a faster recording speed
var numListener int = 5    // the number of simulated live listeners

var eventKey string = ""           // the number of simulated live listeners
var testDurationArg string = "10m" // default testDuration

var speed = 1         // use this replay speed
var numRuns = 1       // how many repetitions
var numStates = 30    // how many states should be fetched in go
var numSpeedMaps = 30 // how many speedmaps should be fetched in go
var raceLimitMin = -1 // if > 0, pick only races shorter than this amount
