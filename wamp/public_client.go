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
)

type PublicAccess interface {
	GetEvent(eventId int) (*internal.Event, error)
}

type PublicClient struct {
	client *client.Client
}

func NewPublicClient(url string, realm string) *PublicClient {
	logger := log.New(os.Stdout, "", 0)
	cfg := client.Config{Realm: realm, Logger: logger}
	// Connect wampClient session.
	wampClient, err := client.ConnectNet(context.Background(), url, cfg)
	if err != nil {
		logger.Fatal(err)
	}

	ret := &PublicClient{client: wampClient}
	return ret
}

func (pc *PublicClient) Close() {
	pc.client.Close()
}

func (pc *PublicClient) Client() *client.Client {
	return pc.client
}

func (pc *PublicClient) ProviderList() ([]*internal.ProviderData, error) {
	ctx := context.Background()

	result, err := pc.client.Call(ctx, "racelog.public.list_providers", nil, wamp.List{}, nil, nil)
	if err != nil {
		return nil, err
	}
	// results in general can contain tuple. Each tuple item can have a different type
	// here we know: this rpc returns one item which is of type list

	// TODO: allgemeine func deserializeEvents(events wamp.List) []*internal.Events

	retEvents := make([]*internal.ProviderData, 0)
	for i := range result.Arguments {
		ret, _ := wamp.AsList(result.Arguments[i])
		for j := range ret {
			var e internal.ProviderData

			jsonData, _ := json.Marshal(ret[j])
			// logger.Printf("jsonData: %v", string(jsonData))

			json.Unmarshal(jsonData, &e)
			retEvents = append(retEvents, &e)
		}
	}
	// fmt.Printf("%v\n", retEvents)
	return retEvents, nil

}

func (pc *PublicClient) GetEvent(eventId int) (*internal.Event, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.get_event_info", nil, wamp.List{eventId}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil, fmt.Errorf("no data for eventId %d", eventId)
	}
	var e internal.Event
	jsonData, _ := json.Marshal(ret)
	// logger.Printf("jsonData: %v", string(jsonData))
	json.Unmarshal(jsonData, &e)
	return &e, nil
}

func (pc *PublicClient) GetEventByKey(eventKey string) (*internal.Event, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.get_event_info_by_key", nil, wamp.List{eventKey}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil, fmt.Errorf("no data for eventKey %s", eventKey)
	}
	var e internal.Event
	jsonData, _ := json.Marshal(ret)
	// logger.Printf("jsonData: %v", string(jsonData))
	json.Unmarshal(jsonData, &e)
	return &e, nil
}

func (pc *PublicClient) GetTrack(id int) (*internal.TrackInfo, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.get_track_info", nil, wamp.List{id}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil, fmt.Errorf("no data for trackId %d", id)
	}
	var t internal.TrackInfo
	jsonData, _ := json.Marshal(ret)
	// logger.Printf("jsonData: %v", string(jsonData))
	json.Unmarshal(jsonData, &t)
	return &t, nil
}

func (pc *PublicClient) GetEventList() ([]*internal.Event, error) {
	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.get_events", nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	retEvents := make([]*internal.Event, 0)
	for i := range result.Arguments {
		ret, _ := wamp.AsList(result.Arguments[i])
		for j := range ret {

			var e internal.Event
			jsonData, _ := json.Marshal(ret[j])
			// logger.Printf("jsonData: %v", string(jsonData))
			json.Unmarshal(jsonData, &e)
			retEvents = append(retEvents, &e)

		}
	}
	return retEvents, nil
}

func (pc *PublicClient) GetStates(id int, start float64, num int) []internal.State {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.archive.state.delta", nil, wamp.List{id, start, num}, nil, nil)
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

func (pc *PublicClient) GetLiveAnalysisData(eventKey string) (map[string]interface{}, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.live.get_event_analysis", nil, wamp.List{eventKey}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil, fmt.Errorf("no data for event %s", eventKey)
	}

	jsonData, _ := json.Marshal(ret)
	var retStruct map[string]interface{}
	json.Unmarshal(jsonData, &retStruct)
	return retStruct, nil

}

func (pc *PublicClient) GetCarData(eventId int) (*internal.EventCarMessage, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.get_event_cars", nil, wamp.List{eventId}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsDict(result.Arguments[0])
	if len(ret) == 0 {
		return nil, fmt.Errorf("no data for event %d", eventId)
	}

	jsonData, _ := json.Marshal(ret)
	var retStruct internal.EventCarMessage
	json.Unmarshal(jsonData, &retStruct)
	return &retStruct, nil

}

func (pc *PublicClient) GetSpeedmaps(id int, start float64, num int) ([]*internal.SpeedmapMessage, error) {
	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.archive.speedmap", nil, wamp.List{id, start, num}, nil, nil)
	if err != nil {
		return nil, err
	}

	ret, _ := wamp.AsList(result.Arguments[0])
	speedmaps := make([]*internal.SpeedmapMessage, 0)
	for j := range ret {
		var s internal.SpeedmapMessage
		jsonData, _ := json.Marshal(ret[j])
		// logger.Printf("jsonData: %v", string(jsonData))
		json.Unmarshal(jsonData, &s)

		speedmaps = append(speedmaps, &s)
	}
	return speedmaps, nil

}

func (pc *PublicClient) GetEventAvgLaps(id int, interval int) ([]*internal.AverageLapTime, error) {

	ctx := context.Background()
	result, err := pc.client.Call(ctx, "racelog.public.archive.avglap_over_time", nil, wamp.List{id, interval}, nil, nil)
	if err != nil {
		return nil, err
	}

	work, _ := wamp.AsList(result.Arguments[0])
	ret := make([]*internal.AverageLapTime, 0)
	for _, item := range work {
		alt := internal.AverageLapTime{}
		jsonData, _ := json.Marshal(item)
		json.Unmarshal(jsonData, &alt)
		ret = append(ret, &alt)
	}
	return ret, nil

}
