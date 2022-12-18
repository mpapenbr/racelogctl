package wamp

import (
	"context"
	"log"
	"os"

	"github.com/gammazero/nexus/v3/client"
)

var logger = log.New(os.Stdout, "client", 0)

func GetClientWithConfigNew(url string, cfg *client.Config) *client.Client {

	// Connect wampClient session.
	wampClient, err := client.ConnectNet(context.Background(), url, *cfg)
	if err != nil {
		logger.Fatal(err)
	}

	return wampClient
}
