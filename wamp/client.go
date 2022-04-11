package wamp

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"racelogctl/internal"
	"racelogctl/util"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
)

// consumer function get an Event and the current index within a list.
// consumer may return false if no further Events should be passed.
type EventConsumer func(*internal.Event, int) bool

var logger = log.New(os.Stdout, "client", 0)

func GetEvents(consumer EventConsumer) {
	client := getClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.get_events", nil, nil, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	for i := range result.Arguments {
		ret, _ := wamp.AsList(result.Arguments[i])
		for j := range ret {

			var e internal.Event
			jsonData, _ := json.Marshal(ret[j])
			// logger.Printf("jsonData: %v", string(jsonData))
			json.Unmarshal(jsonData, &e)
			if !consumer(&e, j) {
				return
			}

		}
	}
}

func GetEvent(id int) *internal.Event {
	client := getClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.get_event_info", nil, wamp.List{id}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil
	}
	var e internal.Event
	jsonData, _ := json.Marshal(ret)
	// logger.Printf("jsonData: %v", string(jsonData))
	json.Unmarshal(jsonData, &e)
	return &e
}

type Delta struct {
	idx   int
	value interface{}
}

func GetStates(id int, event *internal.Event, start int, num int) []internal.State {
	client := getClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.archive.state.delta", nil, wamp.List{id, start, num}, nil, nil)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	ret, _ := wamp.AsList(result.Arguments[0])
	lastState := internal.State{}
	resultStates := make([]internal.State, 0)
	for j := range ret {
		s := internal.State{}
		jsonData, _ := json.Marshal(ret[j])
		// logger.Printf("jsonData: %v", string(jsonData))
		json.Unmarshal(jsonData, &s)
		lastState = util.ProcessDeltaStates(lastState, s)
		resultStates = append(resultStates, lastState)
	}
	return resultStates

}

func GetStatesWithCallback(id int, event *internal.Event, start int, num int) []internal.State {
	client := getClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.archive.state.delta", nil, wamp.List{id, start, num}, nil, nil)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	ret, _ := wamp.AsList(result.Arguments[0])
	// fmt.Printf("ret: %v\n", ret)
	s := internal.State{}
	jsonData, _ := json.Marshal(ret)
	logger.Printf("jsonData: %v", string(jsonData))
	json.Unmarshal(jsonData, &s)
	// logger.Printf("%+v\n", s)
	return nil
}

func getClient() *client.Client {
	logger := log.New(os.Stdout, "", 0)
	cfg := client.Config{Realm: internal.Realm, Logger: logger}
	// Connect wampClient session.
	wampClient, err := client.ConnectNet(context.Background(), internal.Url, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return wampClient
}
