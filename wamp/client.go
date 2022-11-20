package wamp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"racelogctl/internal"
	"racelogctl/util"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/mitchellh/mapstructure"
)

// consumer function get an Event and the current index within a list.
// consumer may return false if no further Events should be passed.
type EventConsumer func(*internal.Event, int) bool

var logger = log.New(os.Stdout, "client", 0)

func GetEventList() []internal.Event {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.get_events", nil, nil, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}

	retEvents := make([]internal.Event, 0)
	for i := range result.Arguments {
		ret, _ := wamp.AsList(result.Arguments[i])
		for j := range ret {

			var e internal.Event
			jsonData, _ := json.Marshal(ret[j])
			// logger.Printf("jsonData: %v", string(jsonData))
			json.Unmarshal(jsonData, &e)
			retEvents = append(retEvents, e)

		}
	}
	return retEvents
}

func GetEvents(consumer EventConsumer) {
	client := GetClient()
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
	client := GetClient()
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

func DeleteEvent(id int) {
	client := getAdminClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.admin.event.delete", nil, wamp.List{id}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("result %v\n", result)

}

func ProcessEvent(id int) internal.ResultMessage {
	client := getAdminClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.admin.event.process", nil, nil, wamp.Dict{"eventId": id}, nil)
	if err != nil {
		logger.Fatal(err)
	}
	if len(result.Arguments) > 0 {
		var resultMsg internal.ResultMessage
		// fmt.Printf("%+v", result.Arguments[0])
		mapstructure.Decode(result.Arguments[0], &resultMsg)
		// fmt.Printf("%+v", resultMsg)
		return resultMsg
	}
	return internal.ResultMessage{}

}

func GetStates(id int, event *internal.Event, start float64, num int) []internal.State {
	client := GetClient()
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
		s := internal.State{Payload: internal.Payload{}}
		jsonData, _ := json.Marshal(ret[j])
		// logger.Printf("jsonData: %v", string(jsonData))
		json.Unmarshal(jsonData, &s)
		lastState = util.ProcessDeltaStates(lastState, s)
		resultStates = append(resultStates, lastState)
	}
	return resultStates

}

func GetStatesWithClient(client *client.Client, id int, event *internal.Event, start float64, num int) []internal.State {

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
		s := internal.State{Payload: internal.Payload{}}
		jsonData, _ := json.Marshal(ret[j])
		// logger.Printf("jsonData: %v", string(jsonData))
		json.Unmarshal(jsonData, &s)
		lastState = util.ProcessDeltaStates(lastState, s)
		resultStates = append(resultStates, lastState)
	}
	return resultStates

}

func GetSpeedmaps(id int, event *internal.Event, start float64, num int) []internal.State {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.archive.speedmap", nil, wamp.List{id, start, num}, nil, nil)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	ret, _ := wamp.AsList(result.Arguments[0])
	fmt.Printf("%v", ret)
	return []internal.State{}

}

func GetEventAvgLaps(id int, interval int) []internal.AverageLapTime {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.archive.avglap_over_time", nil, wamp.List{id, interval}, nil, nil)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	work, _ := wamp.AsList(result.Arguments[0])
	ret := make([]internal.AverageLapTime, 0)
	for _, item := range work {
		alt := internal.AverageLapTime{}
		jsonData, _ := json.Marshal(item)
		json.Unmarshal(jsonData, &alt)
		ret = append(ret, alt)
	}
	return ret

}

func RegisterProvider(registerMsg internal.RegisterMessage) {
	client := getDataproviderClient()
	defer client.Close()
	ctx := context.Background()
	_, err := client.Call(ctx, "racelog.dataprovider.register_provider", nil, wamp.List{registerMsg}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	// logger.Printf("%v", result)

}

func UnregisterProvider(eventKey string) {
	client := getDataproviderClient()
	defer client.Close()
	ctx := context.Background()
	_, err := client.Call(ctx, "racelog.dataprovider.remove_provider", nil, wamp.List{eventKey}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	// logger.Printf("%v", result)

}

func ConsumeProviders(consumer EventConsumer) {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.list_providers", nil, wamp.List{}, nil, nil)
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

func ListProviders() []internal.Event {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.list_providers", nil, wamp.List{}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}

	retEvents := make([]internal.Event, 0)
	for i := range result.Arguments {
		ret, _ := wamp.AsList(result.Arguments[i])
		for j := range ret {

			var e internal.Event
			jsonData, _ := json.Marshal(ret[j])
			// logger.Printf("jsonData: %v", string(jsonData))
			json.Unmarshal(jsonData, &e)
			retEvents = append(retEvents, e)

		}
	}
	return retEvents

}

func GetLiveAnalysisData(eventKey string) map[string]interface{} {
	client := GetClient()
	defer client.Close()
	ctx := context.Background()
	result, err := client.Call(ctx, "racelog.public.live.get_event_analysis", nil, wamp.List{eventKey}, nil, nil)
	if err != nil {
		logger.Fatal(err)
		return nil
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil
	}

	jsonData, _ := json.Marshal(ret)
	var retStruct map[string]interface{}
	json.Unmarshal(jsonData, &retStruct)
	return retStruct

}

func WithDataProviderClient(eventKey string, rcv chan internal.State) {

	// myCount := 0
	go func() {
		client := getDataproviderClient()
		defer client.Close()
		// ctx := context.Background()

		for {
			s, more := <-rcv
			err := client.Publish(fmt.Sprintf("racelog.public.live.state.%s", eventKey), nil, wamp.List{s}, nil)
			if err != nil {
				log.Fatal(err)
			}
			// fmt.Printf("chanValue: %v more: %v\n", s.Timestamp, more)
			// time.Sleep(100 * time.Millisecond)
			if !more {
				fmt.Println("closed channel signaled")
				return
			}
		}
	}()

}

func getDataproviderClient() *client.Client {
	logger := log.New(os.Stdout, "", 0)
	cfg := client.Config{
		Realm:        internal.Realm,
		Logger:       logger,
		HelloDetails: wamp.Dict{"authid": "dataprovider"}, // TODO
		AuthHandlers: map[string]client.AuthFunc{
			"ticket": func(*wamp.Challenge) (string, wamp.Dict) { return internal.DataproviderPassword, wamp.Dict{} }, //TODO:
		}}
	return getClientWithConfig(&cfg)
}

func getAdminClient() *client.Client {
	logger := log.New(os.Stdout, "", 0)
	cfg := client.Config{
		Realm:        internal.Realm,
		Logger:       logger,
		HelloDetails: wamp.Dict{"authid": "admin"}, // TODO
		AuthHandlers: map[string]client.AuthFunc{
			"ticket": func(*wamp.Challenge) (string, wamp.Dict) { return internal.AdminPassword, wamp.Dict{} },
		}}
	return getClientWithConfig(&cfg)
}

func GetClient() *client.Client {
	logger := log.New(os.Stdout, "", 0)
	cfg := client.Config{Realm: internal.Realm, Logger: logger}
	// Connect wampClient session.
	wampClient, err := client.ConnectNet(context.Background(), internal.Url, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return wampClient
}

func getClientWithConfig(cfg *client.Config) *client.Client {

	// Connect wampClient session.
	wampClient, err := client.ConnectNet(context.Background(), internal.Url, *cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return wampClient
}
