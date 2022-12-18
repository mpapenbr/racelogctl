package wamp

import (
	"context"
	"fmt"
	"log"
	"os"
	"racelogctl/internal"

	"github.com/gammazero/nexus/v3/client"
	"github.com/gammazero/nexus/v3/wamp"
	"github.com/mitchellh/mapstructure"
)

type Admin interface {
	GetEvent(eventId int) (*internal.Event, error)
}

type AdminClient struct {
	client *client.Client
}

func NewAdminClient(url string, realm string, ticket string) *AdminClient {
	logger := log.New(os.Stdout, "", 0)

	cfg := client.Config{
		Realm:        realm,
		Logger:       logger,
		HelloDetails: wamp.Dict{"authid": "admin"}, // TODO
		AuthHandlers: map[string]client.AuthFunc{
			"ticket": func(*wamp.Challenge) (string, wamp.Dict) { return ticket, wamp.Dict{} },
		}}

	ret := &AdminClient{client: GetClientWithConfigNew(url, &cfg)}
	return ret
}

func (ac *AdminClient) Close() {
	ac.client.Close()
}

func (ac *AdminClient) DeleteEvent(eventId int) error {
	ctx := context.Background()
	result, err := ac.client.Call(ctx, "racelog.admin.event.delete", nil, wamp.List{eventId}, nil, nil)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Printf("result %v\n", result)
	return nil
}

func (ac *AdminClient) ProcessEvent(id int) internal.ResultMessage {

	ctx := context.Background()
	result, err := ac.client.Call(ctx, "racelog.admin.event.process", nil, nil, wamp.Dict{"eventId": id}, nil)
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
