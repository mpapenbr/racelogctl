/*
Copyright © 2022 Markus Papenbrock

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"crypto/md5"
	"fmt"
	"log"
	"racelogctl/internal"
	"racelogctl/util"
	"racelogctl/wamp"
	"strconv"
	"time"

	"github.com/blang/semver/v4"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// copyCmd represents the copy command
var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copies complete event data from one server to another.",
	Long: `Copies complete event data from one server to another.	

Example: This copies the event with id 42 from the server running at crossbar.mydomain.com 
racelogctl event copy 42 \
   --source-url wss://crossbar.mydomain.com/ws \
   --dataprovider-password verySecret
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		// println("in event_info preRunE")
		bindFlags(cmd, viper.GetViper())
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		eventId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			return
		} else {
			eventCopy(eventId)
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	eventCmd.AddCommand(copyCmd)
	copyCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")

	copyCmd.Flags().StringVar(&internal.SourceUrl, "source-url", "", "sets the url of the source server")

	// TODO: reactivate when doing a real copy
	// copyCmd.MarkFlagRequired("target-url")
	// copyCmd.MarkFlagRequired("target-dataprovider-password")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// copyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// copyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type copyParam struct {
	source         *wamp.PublicClient
	target         *wamp.DataProviderClient
	sourceEventId  int
	targetEventKey string
}

func eventCopy(eventId int) {
	// TODO: need to get information from source-url
	var sourcePc *wamp.PublicClient
	var destPc *wamp.PublicClient

	if len(internal.SourceUrl) != 0 {
		sourcePc = wamp.NewPublicClient(internal.SourceUrl, internal.Realm)
		destPc = wamp.NewPublicClient(internal.Url, internal.Realm)
		defer destPc.Close()
	} else {
		sourcePc = wamp.NewPublicClient(internal.Url, internal.Realm)
		destPc = sourcePc
	}
	defer sourcePc.Close()
	event, err := sourcePc.GetEvent(eventId)
	if err != nil {
		log.Fatalf("Source not found")
	}
	// fmt.Printf("%+v\n", event)

	track, err := sourcePc.GetTrack(event.Data.Info.TrackId)
	if err != nil {
		log.Fatalf("Track not found")
	}

	dpc := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	uuid, _ := uuid.NewRandom()
	md5 := md5.New()
	md5.Write([]byte(uuid.String()))

	eventKey := fmt.Sprintf("%x", md5.Sum(nil))

	recDate, _ := time.Parse("2006-01-02T15:04:05", event.RecordDate)
	registerMsg := internal.RegisterMessage{
		Manifests:  event.Data.Manifests,
		EventKey:   eventKey,
		Info:       event.Data.Info,
		TrackInfo:  *track,
		RecordDate: float64(recDate.Unix()),
	}
	err = dpc.RegisterProvider(registerMsg)
	if err != nil {
		log.Fatalf("Error registering event: %v", err)
	}

	targetEvent, _ := destPc.GetEventByKey(eventKey)
	fmt.Println("Created event on target:")
	printEvent(targetEvent)

	speedAndCarDataAvail := semver.MustParseRange(">=0.4.4")
	param := copyParam{source: sourcePc, target: dpc, sourceEventId: eventId, targetEventKey: eventKey}

	copyStandardData(param)
	if speedAndCarDataAvail(semver.MustParse(util.GetEventRaceloggerVersion(event))) {
		copyCarData(param)
		copySpeedData(param)
	}
	err = dpc.UnregisterProvider(eventKey)
	if err != nil {
		log.Fatalf("Error unregistering event: %v", err)
	}

}

func copyStandardData(param copyParam) {
	log.Println("begin copy states")

	fetches := 0
	numPackets := 0

	sender := make(chan internal.State)

	param.target.PublishStateFromChannel(param.targetEventKey, sender)

	from := 0.0
	for goon := true; goon; {
		// fmt.Printf("Fetching %d states beginning at %d\n", numStates, int64(from))
		states := param.source.GetStates(param.sourceEventId, from, 100)
		fetches += 1
		numPackets += len(states)
		// fmt.Printf("Got %v states\n", len(states))
		goon = len(states) > 0 // == numStates
		if goon {
			for _, state := range states {
				sender <- state
			}
			from = states[len(states)-1].Timestamp + 0.0001
		}

	}
	log.Printf("done copy states: fetches %d packets: %d", fetches, numPackets)

}

func copyCarData(param copyParam) {
	log.Println("begin copy car data")

	carData, _ := param.source.GetCarData(param.sourceEventId)
	param.target.PublishCarData(param.targetEventKey, carData)

	log.Println("done copy car data")

}

func copySpeedData(param copyParam) {
	log.Println("begin copy speedmap data")
	sender := make(chan internal.SpeedmapMessage)
	fetches := 0
	numPackets := 0

	param.target.PublishSpeedmapDataFromChannel(param.targetEventKey, sender)

	from := 0.0
	for goon := true; goon; {
		// fmt.Printf("Fetching %d speedmaps beginning at %d\n", numStates, int64(from))
		speedmaps, _ := param.source.GetSpeedmaps(param.sourceEventId, from, 100)

		fetches += 1
		numPackets += len(speedmaps)
		// fmt.Printf("Got %v states\n", len(speedmaps))
		goon = len(speedmaps) > 0 // == numStates
		if goon {
			for _, speedmap := range speedmaps {
				sender <- *speedmap
			}
			from = speedmaps[len(speedmaps)-1].Timestamp + 0.0001
		}

	}
	log.Printf("done copy speedmaps: fetches %d packets: %d", fetches, numPackets)
}
