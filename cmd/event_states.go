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
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"

	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:   "states",
	Short: "Retrieves state data from server",
	Long:  ``,
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
				fetchStates(eventId, outFile)
			}
		} else {
			fmt.Println("requires an event id")
		}
	},
}

func init() {
	eventCmd.AddCommand(stateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	stateCmd.PersistentFlags().IntVar(&internal.From, "from", 0, "Fetch states beginning from timestamp (Default: 0=first available entry)")
	stateCmd.PersistentFlags().IntVar(&internal.Num, "num", 10, "How many states should be fetches in one request")
	stateCmd.PersistentFlags().BoolVar(&internal.FullStateData, "full", false, "retrieves all data for this event")
	stateCmd.PersistentFlags().StringVar(&internal.Output, "output", "-", "Output filename. (Default: stdout)")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchStates(eventId int, outFile *os.File) {
	event := wamp.GetEvent(eventId)
	fmt.Printf("event: %v\n", event)
	if internal.FullStateData {
		fetchFullData(event, outFile)
		return
	}
	fmt.Printf("Fetching %d states beginning at %d\n", internal.Num, internal.From)
	states := wamp.GetStates(eventId, event, internal.From, internal.Num)
	fmt.Printf("\n---\nresulting states\n")
	for _, entry := range states {
		jsonData, _ := json.Marshal(entry)
		outFile.WriteString(fmt.Sprintln(string(jsonData)))

	}
}

func fetchFullData(event *internal.Event, outFile *os.File) {
	// var lastTimestamp float64 = 0
	from := int(event.Data.ReplayInfo.MinTimestamp)
	if internal.From != 0 {
		from = internal.From
	}
	for goon := true; goon; {
		fmt.Printf("Fetching %d states beginning at %d\n", internal.Num, from)
		states := wamp.GetStates(int(event.Id), event, from, internal.Num)
		goon = len(states) == internal.Num
		for _, entry := range states {
			jsonData, _ := json.Marshal(entry)
			outFile.WriteString(fmt.Sprintln(string(jsonData)))
		}
		from = int(states[len(states)-1].Timestamp)
	}

}
