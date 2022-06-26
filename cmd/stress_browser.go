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
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"
	"sort"
	"sync"
	"time"

	"github.com/gammazero/nexus/v3/client"
	"github.com/spf13/cobra"
)

var speed = 1         // use this replay speed
var numRuns = 1       // how many repetitions
var numStates = 30    // how many states should be fetched in go
var raceLimitMin = -1 // if > 0, pick only races shorter than this amount

// browserCmd represents the browser command
var browserCmd = &cobra.Command{
	Use:   "browser",
	Short: "Simulates the browser requests to perfom stress tests",

	Run: func(cmd *cobra.Command, args []string) {
		simulateBrowser()
	},
}

type durationStats struct {
	Min time.Duration `json:"min"`
	Max time.Duration `json:"max"`
	Avg time.Duration `json:"avg"`
	Sum time.Duration `json:"sum"`
}

type jobData struct {
	id    int
	event *internal.Event
}
type jobResult struct {
	jobId      int
	workerId   int
	event      *internal.Event
	duration   time.Duration
	numFetches int
	numStates  int
}

type summary struct {
	Id            int             `json:"id"`
	Num           int             `json:"num"`
	NumStates     int             `json:"numStates"`
	NumFetches    int             `json:"numFetches"`
	Durations     []time.Duration `json:"durations"`
	DurationStats durationStats   `json:"durationStats"`
	Name          string          `json:"name"`
}

type statistics struct {
	eventSummary  map[int]summary
	workerSummary map[int]summary
}

// support for sorting summary arrays
type byId []summary

func (s byId) Len() int           { return len(s) }
func (s byId) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byId) Less(i, j int) bool { return s[i].Id < s[j].Id }

func init() {
	stressCmd.AddCommand(browserCmd)
	browserCmd.Flags().IntVar(&speed, "speed", 1, "Replay speed (<=0 means: go as fast as possible)")
	browserCmd.Flags().IntVar(&numStates, "num-states", numStates, "How many states should be fetched in one request")
	browserCmd.Flags().IntVar(&numRuns, "num-runs", numRuns, "Number of test runs")
	browserCmd.Flags().IntVar(&raceLimitMin, "race-limit", raceLimitMin, "max race length (in minutes) to consider")
}

func simulateBrowser() {
	events := wamp.GetEventList()

	queue := make(chan *jobData)
	results := make(chan *jobResult)
	statistics := make(chan *statistics)
	go createJobs(queue, events, numRuns)

	wg := sync.WaitGroup{}

	go createResultCollector(events, results, statistics, &wg)
	for i := 0; i < internal.Worker; i++ {
		wg.Add(1)
		fmt.Printf("Starting worker %d\n", i)
		go worker(i, queue, results, &wg)
	}

	wg.Wait()

	close(results)

	printStatsSummary(<-statistics, events)

}

func printStatsSummary(stats *statistics, events []internal.Event) {
	lookup := make(map[int]internal.Event, len(events))
	for _, e := range events {
		lookup[int(e.Id)] = e
	}

	eventsSummary := processSummary(stats.eventSummary, func(s summary) string {
		return fmt.Sprintf("Event: %v-%v", lookup[s.Id].Id, lookup[s.Id].Name)
	})
	workerSummary := processSummary(stats.workerSummary, func(s summary) string {
		return fmt.Sprintf("Worker: %v", s.Id)
	})
	fmt.Printf("\nSummary by events\n")
	printSummary(eventsSummary)
	fmt.Printf("\nSummary by workers\n")
	printSummary(workerSummary)

	type x struct {
		Events  []summary `json:"events"`
		Workers []summary `json:"workers"`
	}
	j := x{Events: eventsSummary, Workers: workerSummary}
	jsonData, _ := json.Marshal(j)
	// fmt.Printf("\n%v\n", string(jsonData))
	os.WriteFile(fmt.Sprintf("stress-browser-%s.json", time.Now().Format("20060102-150405")), jsonData, 0644)
}

func processSummary(s map[int]summary, title func(summary) string) []summary {
	values := make([]summary, 0, len(s))
	for _, v := range s {
		values = append(values, v)
	}

	sort.Sort(byId(values))
	for idx, item := range values {
		item.DurationStats = minMaxAvg(item.Durations)
		item.Name = title(item)
		values[idx] = item
		// fmt.Printf("%v\n", item)
	}
	// fmt.Printf("%v\n", values)
	return values

}

func printSummary(values []summary) {

	for _, item := range values {
		fmt.Printf("%s\n", item.Name)
		fmt.Printf("Num: %d Total: %s Min: %s Max: %s Avg: %s Fetches: %d States: %d\n", item.Num, item.DurationStats.Sum, item.DurationStats.Min, item.DurationStats.Max, item.DurationStats.Avg, item.NumFetches, item.NumStates)
	}

}

func minMaxAvg(items []time.Duration) durationStats {
	if len(items) == 0 {
		return durationStats{}
	}
	min := time.Duration(math.MaxInt64)
	max := time.Duration(0)
	sum := time.Duration(0)
	for _, arg := range items {
		if arg < min {
			min = arg
		}
		if arg > max {
			max = arg
		}
		sum += arg
	}
	avg := int(sum) / len(items)
	return durationStats{min, max, time.Duration(avg), sum}
}

func createResultCollector(events []internal.Event, results chan *jobResult, stats chan *statistics, wg *sync.WaitGroup) {

	byEventSummary := make(map[int]summary)
	byWorkerSummary := make(map[int]summary)

	for {
		job, ok := <-results
		if ok {
			fmt.Printf("Worker %2d finishied Job %d: %v-%v used %d batches for %d packets in %s\n", job.workerId, job.jobId, job.event.Id, job.event.Name, job.numFetches, job.numStates, job.duration)
			val, other := byEventSummary[int(job.event.Id)]
			if !other {
				val = summary{Id: int(job.event.Id), Durations: []time.Duration{}}
			} else {
				// fmt.Printf("Adding mode for %v\n", job.event.Id)
			}
			val.Num += 1
			val.NumFetches += job.numFetches
			val.NumStates += job.numStates
			val.Durations = append(val.Durations, job.duration)

			byEventSummary[int(job.event.Id)] = val

			val, other = byWorkerSummary[job.workerId]
			if !other {
				val = summary{Id: job.workerId, Durations: []time.Duration{}}
			}
			val.Num += 1
			val.NumFetches += job.numFetches
			val.NumStates += job.numStates
			val.Durations = append(val.Durations, job.duration)
			byWorkerSummary[job.workerId] = val

		} else {
			fmt.Printf("ResultCollector: no more results. Terminating\n")
			// fmt.Printf("%+v\n", byEventSummary)
			stats <- &statistics{eventSummary: byEventSummary, workerSummary: byWorkerSummary}

			return
		}
	}
}

func worker(idx int, queue chan *jobData, results chan *jobResult, wg *sync.WaitGroup) {
	client := wamp.GetClient()
	defer wg.Done()
	defer client.Close()
	for {
		job, ok := <-queue
		if ok {
			start := time.Now()
			numFetches, numPackets := simulateFrontendFetching(client, job.event)
			duration := time.Since((start))
			results <- &jobResult{workerId: idx + 1, jobId: job.id, event: job.event, duration: duration, numFetches: numFetches, numStates: numPackets}
			// fmt.Printf("Job %3d %v-%v done in %s\n", job.id, job.event.Id, job.event.Name, duration)
		} else {
			fmt.Printf("Worker %2d: no more jobs available. Terminating\n", idx+1)

			return
		}
	}
}

func createJobs(ch chan<- *jobData, events []internal.Event, numRuns int) {
	pickShortRace := func() int {
		for {
			pick := rand.Intn(len(events))
			event := events[pick]
			if raceLimitMin > 0 {
				if (event.Data.ReplayInfo.MaxSessionTime - event.Data.ReplayInfo.MinSessionTime) < float64(raceLimitMin) {
					return pick
				}
			} else {
				return pick
			}
		}
	}
	defer close(ch)
	for i := 0; i < numRuns; i++ {
		pick := pickShortRace()
		event := events[pick]

		fmt.Printf("Run %03d picked event %d: %s\n", i+1, event.Id, event.Name)
		ch <- &jobData{id: i + 1, event: &event}
	}
}

func simulateFrontendFetching(client *client.Client, event *internal.Event) (int, int) {
	fetches := 0
	numPackets := 0
	from := event.Data.ReplayInfo.MinTimestamp
	for goon := true; goon; {
		// fmt.Printf("Fetching %d states beginning at %d\n", numStates, int64(from))
		states := wamp.GetStatesWithClient(client, int(event.Id), event, from, numStates)
		fetches += 1
		numPackets += len(states)
		// fmt.Printf("Got %v states\n", len(states))
		goon = len(states) > 0 // == numStates
		if goon {
			if speed > 0 {
				sleep := int64(float64(len(states)*1000) / float64(speed))
				// fmt.Printf("Sleeping for %+v ms\n", sleep)
				time.Sleep(time.Duration(sleep) * time.Millisecond)
			}

			from = states[len(states)-1].Timestamp + 0.0001
		}
	}
	return fetches, numPackets
}
