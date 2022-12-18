package wamp

import (
	"context"
	"fmt"
	"log"
	"os"
	"racelogctl/internal"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
)

type Dataprovider interface {
	GetEvent(eventId int) (*internal.Event, error)
}

type DataProviderClient struct {
	client *client.Client
}

func NewDataProviderClient(url string, realm string, ticket string) *DataProviderClient {
	logger := log.New(os.Stdout, "", 0)

	cfg := client.Config{
		Realm:        realm,
		Logger:       logger,
		HelloDetails: wamp.Dict{"authid": "dataprovider"}, // TODO
		AuthHandlers: map[string]client.AuthFunc{
			"ticket": func(*wamp.Challenge) (string, wamp.Dict) { return ticket, wamp.Dict{} },
		}}

	ret := &DataProviderClient{client: GetClientWithConfigNew(url, &cfg)}
	return ret
}

func (dpc *DataProviderClient) Close() {
	dpc.client.Close()
}

// registers a new provider
func (dpc *DataProviderClient) RegisterProvider(registerMsg internal.RegisterMessage) error {
	ctx := context.Background()
	_, err := dpc.client.Call(ctx, "racelog.dataprovider.register_provider", nil, wamp.List{registerMsg}, nil, nil)
	return err

}

// unregisters a provider
func (dpc *DataProviderClient) UnregisterProvider(eventKey string) error {
	ctx := context.Background()
	_, err := dpc.client.Call(ctx, "racelog.dataprovider.remove_provider", nil, wamp.List{eventKey}, nil, nil)
	return err
}

// recieves data via channel and publishes it on the racelog.public.live.state.<eventKey> topic
func (dpc *DataProviderClient) PublishStateFromChannel(eventKey string, rcv chan internal.State) {

	go func() {
		for {
			s, more := <-rcv
			err := dpc.client.Publish(fmt.Sprintf("racelog.public.live.state.%s", eventKey), nil, wamp.List{s}, nil)
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

func (dpc *DataProviderClient) PublishCarData(eventKey string, carData *internal.EventCarMessage) {

	err := dpc.client.Publish(fmt.Sprintf("racelog.public.live.cardata.%s", eventKey), nil, wamp.List{carData}, nil)
	if err != nil {
		log.Fatal(err)
	}

}
func (dpc *DataProviderClient) PublishSpeedmapDataFromChannel(eventKey string, rcv chan internal.SpeedmapMessage) {

	go func() {
		for {
			s, more := <-rcv
			err := dpc.client.Publish(fmt.Sprintf("racelog.public.live.speedmap.%s", eventKey), nil, wamp.List{s}, nil)
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
