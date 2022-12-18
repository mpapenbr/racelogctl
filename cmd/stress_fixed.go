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
	"io"
	"log"
	"racelogctl/internal"
	"racelogctl/wamp"
	"time"

	nexusWamp "github.com/gammazero/nexus/v3/wamp"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// fixedCmd represents the fixed command
var fixedCmd = &cobra.Command{
	Use:   "fixed",
	Short: "Performs a stress test by recording an event",
	Long: `Performs a stress test by recording an event.
	
NOTE: This command performs the recording of an event while a number of 
clients will be connected to the live server.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		simulateLiveRecording()
	},
}

func init() {
	stressCmd.AddCommand(fixedCmd)

	fixedCmd.Flags().IntVar(&recordingSpeed, "speed", 1, "Recording speed (<=0 means: go as fast as possible)")
	fixedCmd.Flags().IntVar(&numListener, "num-listener", numListener, "How many states should be fetched in one request")

	fixedCmd.Flags().IntVar(&sourceEventId, "eventId", sourceEventId, "the id of the source event")
	fixedCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")

	fixedCmd.Flags().StringVar(&eventKey, "eventKey", "", "sets the event key")

}

func simulateLiveRecording() {
	if sourceEventId == -1 {
		fmt.Printf("we could pick a random event here. For now we do nothing\n")
		return

	}
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()
	event, err := pc.GetEvent(sourceEventId)
	if err != nil {
		log.Fatalf("Error getting event: %v\n", err)
	}

	if len(event.Data.Info.RaceloggerVersion) == 0 {
		log.Fatalf("Event %v %v not suitable for this function. Need at least raceLoggerVersion 0.4.0", sourceEventId, event.Name)

	}

	registerMsg := internal.RegisterMessage{}

	registerMsg.Manifests = event.Data.Manifests
	registerMsg.Info = event.Data.Info
	registerMsg.Info.Name = fmt.Sprintf("stresstest-%s", time.Now().Format("20060102-150405"))
	h := md5.New()
	io.WriteString(h, registerMsg.Info.Name)
	io.WriteString(h, uuid.New().String())
	registerMsg.EventKey = fmt.Sprintf("%x", h.Sum(nil))

	if eventKey != "" {
		registerMsg.EventKey = eventKey

	}
	dpc := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	dpc.RegisterProvider(registerMsg)
	defer dpc.Close()
	producerDone := make(chan bool)
	// create producer
	go simulateRacelogger(event, registerMsg.EventKey, producerDone)

	// create live consumer
	// wg := sync.WaitGroup{}

	for i := 0; i < numListener; i++ {

		fmt.Printf("Starting listener %d\n", i)
		go simulateBrowserListener(i, registerMsg.EventKey)
	}

	// wg.Wait()

	<-producerDone

	log.Printf("Producer done\n")

	dpc.UnregisterProvider(registerMsg.EventKey)

	log.Printf("Unregistered event\n")

	time.Sleep(time.Duration(2) * time.Second)
	log.Printf("Wait done\n")
}

func simulateBrowserListener(idx int, eventKey string) {

	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()

	topic := fmt.Sprintf("racelog.public.live.state.%s", eventKey)
	msgNum := 0
	handler := func(event *nexusWamp.Event) {
		msgNum++
		// eventData, b := nexusWamp.AsList(event.Arguments[0])

		if (msgNum % 100) == 0 {

			log.Printf("Listener %d - Event data for topic %s: msgNum %d \n", idx, eventKey, msgNum)
			// log.Printf("Event: %+v\n", event)
		}
	}
	err := pc.Client().Subscribe(topic, handler, nil)
	if err != nil {
		log.Fatal("subscribe error: ", err)
	}
	<-pc.Client().Done()
	log.Printf("subsriber %d finished\n", idx)
}

func simulateRacelogger(event *internal.Event, recordingEventKey string, done chan bool) {
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()

	fetches := 0
	numPackets := 0

	sender := make(chan internal.State)
	dataprovider := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	defer dataprovider.Close()
	dataprovider.PublishStateFromChannel(recordingEventKey, sender)

	from := event.Data.ReplayInfo.MinTimestamp
	for goon := true; goon; {
		// fmt.Printf("Fetching %d states beginning at %d\n", numStates, int64(from))
		states := pc.GetStates(int(event.Id), from, numStates)
		fetches += 1
		numPackets += len(states)
		// fmt.Printf("Got %v states\n", len(states))
		goon = len(states) > 0 // == numStates
		if goon {
			for _, state := range states {

				sender <- state
				if recordingSpeed > 0 {
					sleep := int64(1000 / float64(recordingSpeed))
					fmt.Printf("Sleeping for %+v ms\n", sleep)
					time.Sleep(time.Duration(sleep) * time.Millisecond)
				}
			}

			from = states[len(states)-1].Timestamp + 0.0001
		}
	}
	done <- true
}
