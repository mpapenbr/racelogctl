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
	"fmt"
	"log"
	"math/rand"
	"racelogctl/internal"
	"racelogctl/wamp"
	"sync"
	"time"

	nexusWamp "github.com/gammazero/nexus/v3/wamp"
	"github.com/spf13/cobra"
)

// liveCmd represents the live command
var liveCmd = &cobra.Command{
	Use:   "live",
	Short: "Simulates clients listenings to live events",
	Long: `Simulates clients listenings to live events

Each worker will do the following:
- pick one live event
- get the current analysis data
- listen to event topic for <workerListenDuration> 
- wait for <workerPauseDuration> before next run

	`,

	Run: func(cmd *cobra.Command, args []string) {
		setupScenario()
	},
	// TODO: validate args (testDuration)
}

var workerListenDurationArg string = ""
var workerPauseDurationArg string = "5s"

var jobNum int = 0
var workerPause time.Duration
var workerListen time.Duration
var workerListenRandom bool = false

func init() {
	stressCmd.AddCommand(liveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	liveCmd.Flags().StringVarP(&testDurationArg, "duration", "d", testDurationArg, "Overall test duration")
	liveCmd.Flags().StringVar(&workerListenDurationArg, "workerListenDuration", workerListenDurationArg, "Duration for the worker to listen")
	liveCmd.Flags().StringVar(&workerPauseDurationArg, "workerPauseDuration", workerPauseDurationArg, "Pause between 2 worker runs")
	liveCmd.Flags().BoolVar(&workerListenRandom, "workerRandomDuration", workerListenRandom, "if true worker duration will be a random [0..workerListenDuration]")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// liveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setupScenario() {
	ctx, cancel := context.WithCancel(context.Background())

	workerPause, _ = time.ParseDuration(workerPauseDurationArg)
	workerListen, _ = time.ParseDuration(workerListenDurationArg)
	// setup worker for producer
	wg := sync.WaitGroup{}
	queue := make(chan int)
	// go timedResultCollector(queue, results, ctx)

	for i := 0; i < internal.Worker; i++ {
		wg.Add(1)
		fmt.Printf("Starting worker %d\n", i)
		go simBrowserClient(i, queue, &wg, ctx)
		jobNum++
		queue <- jobNum
	}
	// start the test duration timer
	go func() {
		testDuration, _ := time.ParseDuration(testDurationArg)
		log.Printf("Waiting %v to terminate worker\n", testDuration)
		time.Sleep(testDuration)
		log.Printf("signalling cancel\n")
		cancel()
		log.Printf("signalled cancel\n")
	}()

	log.Printf("Waiting for terminating jobs\n")
	wg.Wait()
	log.Printf("All workers finished\n")
}

func simBrowserClient(idx int, queue chan int, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()

	for {
		select {
		case <-ctx.Done():
			log.Printf("test duration reached (outer) Terminating worker %d", idx)
			return
		case dummy := <-queue:
			fmt.Printf("Dummy: %v\n", dummy)

			fmt.Println("get available live events")
			providers, _ := pc.ProviderList()
			if (len(providers)) == 0 {
				log.Println("no event avail. pausing")
				time.Sleep(workerPause)
				go func() {
					log.Println("try issue again")
					queue <- jobNum
					log.Println("put in queue")
				}()

			} else {
				pick := rand.Intn(len(providers))

				go simulateLiveListener(dummy, providers[pick].EventKey, queue)

			}

		}
	}
}

func simulateLiveListener(idx int, eventKey string, queue chan int) {
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)

	defer pc.Close()

	pc.GetLiveAnalysisData(eventKey) // don't need, just to issue the request

	topic := fmt.Sprintf("racelog.public.live.state.%s", eventKey)
	msgNum := 0
	handler := func(event *nexusWamp.Event) {
		msgNum++
		// eventData, b := nexusWamp.AsList(event.Arguments[0])

		log.Printf("Listener %d - Event data for topic %s: msgNum %d \n", idx, eventKey, msgNum)
		// log.Printf("Event: %+v\n", event)

	}
	err := pc.Client().Subscribe(topic, handler, nil)
	if err != nil {
		log.Fatal("subscribe error: ", err)
	}

	go func() {

		unsubTimer := workerListen
		if workerListenRandom {
			unsubTimer, _ = time.ParseDuration(fmt.Sprintf("%ds", rand.Intn(int(workerListen.Seconds()))))
		}
		log.Printf("i: %v Unsub in %v\n", idx, unsubTimer)
		time.Sleep(unsubTimer)
		pc.Client().Unsubscribe(topic)
		log.Printf("i: %v: unsubscribed\n", idx)
		pc.Client().Close()

	}()
	log.Printf("i: %v vor done \n", idx)
	<-pc.Client().Done()
	log.Printf("subsriber %d finished\n", idx)

	time.Sleep(workerPause)
	jobNum++
	queue <- jobNum
}
