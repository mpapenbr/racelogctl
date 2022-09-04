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
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"math/rand"
	"racelogctl/internal"
	"racelogctl/wamp"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

// timedCmd represents the timed command
var timedCmd = &cobra.Command{
	Use:   "timed",
	Short: "Simulate running a number of raceloggers for a limited time",
	Long: `Simulate running a number of raceloggers for a limited time.
	
NOTE: This command performs the recording of events.
Example:
racelogctl stress timed --worker 5 --speed 2 --duration 60m --minSessionDuration 30m

This will start a test for 60 minutes with 5 workers. 
They will produce copies of existing races which last at least 30 minutes.
The recording speed is 2 which means, instead of sending a packet each second, 
they will send a packet every 500 milliseconds.`,
	Run: func(cmd *cobra.Command, args []string) {
		setupTimedProducer()
	},
}

var minSessionDuration string = "5m" // the minimum session duration of the source
var nextJobId = 1
var availableEvents []internal.Event

type TimedJobRequest struct {
	id          int
	eventSource *internal.Event
}

type TimedJobResult struct {
	jobId    int
	workerId int
}

func (j TimedJobRequest) output() string {
	return fmt.Sprintf("JobId: %d Event: %s", j.id, composeEventOverview(j.eventSource))
}

func init() {
	stressCmd.AddCommand(timedCmd)
	timedCmd.Flags().IntVar(&speed, "speed", 1, "Recording speed (<=0 means: go as fast as possible)")
	timedCmd.Flags().StringVar(&testDurationArg, "duration", testDurationArg, "How long should the test run")
	timedCmd.Flags().StringVar(&minSessionDuration, "min-session-duration", minSessionDuration, "the minimum session duration of the source")
	timedCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// timedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// timedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupTimedProducer() {
	minDuration, _ := time.ParseDuration(minSessionDuration)
	availableEvents = computeAvailableEvents(int(minDuration.Minutes()))

	ctx, cancel := context.WithCancel(context.Background())
	// setup worker for producer
	wg := sync.WaitGroup{}
	queue := make(chan *TimedJobRequest)
	results := make(chan *TimedJobResult)

	go timedResultCollector(queue, results, ctx)
	for i := 0; i < internal.Worker; i++ {
		wg.Add(1)
		fmt.Printf("Starting worker %d\n", i)
		go raceloggerWorker(i, queue, results, &wg, ctx)
	}

	for jobId := 1; jobId <= internal.Worker; jobId++ {
		queue <- &TimedJobRequest{id: jobId, eventSource: &availableEvents[rand.Intn(len(availableEvents))]}
	}
	nextJobId = internal.Worker + 1

	// start the test duration timer
	go func() {
		testDuration, _ := time.ParseDuration(testDurationArg)
		log.Printf("Waiting %s to terminate worker\n", testDuration)
		time.Sleep(testDuration)
		log.Printf("signalling cancel\n")
		cancel()
		log.Printf("signalled cancel\n")
	}()

	log.Printf("Waiting for terminating jobs\n")
	wg.Wait()
	log.Printf("All workers finished\n")

	// handle producer finish
	// done
	log.Printf("All done\n")
}

func raceloggerWorker(idx int, requestChan chan *TimedJobRequest, resultChan chan *TimedJobResult, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	client := wamp.GetClient()
	defer client.Close()
	// get a job from the queue
	// do something

	for {
		select {
		case <-ctx.Done():
			log.Printf("test duration reached (outer) Terminating worker %d", idx)
			return
		case job := <-requestChan:
			currentRun := 0
			currentIdx := 0
			log.Printf("Worker %d got job %v\n", idx, job.output())
			recordingEventKey := registerNewEvent(job.eventSource)
			sender := make(chan internal.State)
			finalizeRecorder := func() {
				wamp.UnregisterProvider(recordingEventKey)
				resultChan <- &TimedJobResult{jobId: job.id, workerId: idx}
			}

			wamp.WithDataProviderClient(recordingEventKey, sender)

			from := job.eventSource.Data.ReplayInfo.MinTimestamp
			states := wamp.GetStatesWithClient(client, int(job.eventSource.Id), job.eventSource, from, numStates)

			for goon := true; goon; {
				select {
				case <-ctx.Done():
					log.Printf("test duration reached (inner). Terminating worker %d", idx)
					finalizeRecorder()
					return
				default:
					// log.Printf("Worker %d Iter %d\n", idx, currentRun)
					sender <- states[currentIdx]
					if speed > 0 {
						sleep := 1000 / speed
						// fmt.Printf("Worker %d Sleeping for %+v ms\n", idx, sleep)
						time.Sleep(time.Duration(sleep) * time.Millisecond)
					}
					currentRun++
					currentIdx++
					goon = currentIdx < len(states)
					// check if more states are available
					if !goon {
						from = states[len(states)-1].Timestamp + 0.0001
						states = wamp.GetStatesWithClient(client, int(job.eventSource.Id), job.eventSource, from, numStates)
						log.Printf("Worker %d get %d new states\n", idx, len(states))

						goon = len(states) > 0
						currentIdx = 0
					}
				}
			}

			log.Printf("Worker %d End of task. \n", idx)
			finalizeRecorder()
		}

	}
}

func registerNewEvent(event *internal.Event) string {
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
	wamp.RegisterProvider(registerMsg)
	return registerMsg.EventKey
}

func timedResultCollector(requests chan *TimedJobRequest, results chan *TimedJobResult, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("test duration reached. Terminating result collector")
			for x := range results {

				log.Printf("Collect result: %v\n", x)
			}
			return
		case result, ok := <-results:
			log.Printf("Got result: %v ok: %v\n", result, ok)
			nextJobId++
			requests <- &TimedJobRequest{id: nextJobId, eventSource: &availableEvents[rand.Intn(len(availableEvents))]}
		}
	}
}
