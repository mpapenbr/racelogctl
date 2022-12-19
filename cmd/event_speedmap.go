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
	"log"
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"

	"github.com/spf13/cobra"
)

// speedmapCmd represents the speedmap command
var speedmapCmd = &cobra.Command{
	Use:   "speedmap",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			eventId, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			} else {
				var outFile *os.File = nil
				var err error
				if internal.Output != "-" {
					outFile, err = os.Create(internal.Output)
					if err != nil {
						fmt.Printf("Error creating output file %v: %v", internal.Output, err)
						return
					}
					defer outFile.Close()
				} else {
					outFile = os.Stdout
				}
				fetchSpeedmapRangeEntries(eventId, outFile)
			}
		} else {
			fmt.Println("requires an event id")
		}
	},
}

func init() {
	eventCmd.AddCommand(speedmapCmd)

	speedmapCmd.Flags().IntVar(&internal.From, "from", 0, "Fetch speedmap beginning from timestamp (Default: 0=first available entry)")
	speedmapCmd.Flags().IntVar(&internal.Num, "num", 10, "How many entries should be fetches in one request")
	speedmapCmd.Flags().BoolVar(&internal.FullStateData, "full", false, "retrieves all data for this event")
	speedmapCmd.Flags().StringVar(&internal.Output, "output", "-", "Output filename. (Default: stdout)")

}

func fetchSpeedmapRangeEntries(eventId int, outFile *os.File) {
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()
	event, err := pc.GetEvent(eventId)
	if err != nil {
		log.Fatalf("Error getting event: %v\n", err)
	}
	fmt.Printf("event: %v\n", event)
	if internal.FullStateData {
		fetchSpeedmapFull(event, outFile)
		return
	}
	fmt.Printf("Fetching %d entries beginning at %d\n", internal.Num, internal.From)
	entries, err := pc.GetSpeedmaps(eventId, float64(internal.From), internal.Num)
	fmt.Printf("\n---\nresulting speedmap entries\n")
	for _, entry := range entries {
		jsonData, _ := json.Marshal(entry.Payload)
		outFile.WriteString(fmt.Sprintln(string(jsonData)))

	}
}

func fetchSpeedmapFull(event *internal.Event, outFile *os.File) {
	// var lastTimestamp float64 = 0
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	defer pc.Close()

	from := event.Data.ReplayInfo.MinTimestamp
	if internal.From != 0 {
		from = float64(internal.From)
	}
	for goon := true; goon; {
		fmt.Printf("Fetching %d speedmaps beginning at %v\n", internal.Num, from)
		speedmaps, _ := pc.GetSpeedmaps(int(event.Id), from, internal.Num)
		goon = len(speedmaps) == internal.Num
		for _, entry := range speedmaps {
			jsonData, _ := json.Marshal(entry)
			outFile.WriteString(fmt.Sprintln(string(jsonData)))
		}
		from = speedmaps[len(speedmaps)-1].Timestamp + 0.0001
	}

}