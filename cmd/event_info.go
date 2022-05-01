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
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info <eventId>",
	Short: "Get information about an event",
	Long:  ``,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		// println("in event_info preRunE")
		bindFlags(cmd, viper.GetViper())
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			eventId, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			} else {
				eventInfo(eventId)
			}
		} else {
			fmt.Println("requires an event id")
		}

	},
}

func init() {
	eventCmd.AddCommand(infoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	// infoCmd.PersistentFlags().IntVar(&internal.EventId, "id", -1, "the event id to fetch")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// infoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	infoCmd.Flags().StringVarP(&internal.OutputFormat, "format", "f", "text", "Output format: text|json.")
	infoCmd.Flags().BoolVarP(&internal.JsonPretty, "pretty", "p", false, "use pretty json format. (Default: false)")
}

func eventInfo(id int) {

	event := wamp.GetEvent(id)
	if event == nil {
		fmt.Printf("No event found for %v\n", id)
		return
	}
	// fmt.Printf("%+v\n", event)
	switch internal.OutputFormat {
	case "json":
		printJson(event)
	default:
		printEvent(event)
	}
}

func printEvent(e *internal.Event) {
	recDate, _ := time.Parse("2006-01-02T15:04:05", e.RecordDate)
	time.Unix(int64(e.Data.ReplayInfo.MinTimestamp), 0)
	minSession, _ := time.ParseDuration(fmt.Sprintf("%.0fs", e.Data.ReplayInfo.MinSessionTime))
	maxSession, _ := time.ParseDuration(fmt.Sprintf("%.0fs", e.Data.ReplayInfo.MaxSessionTime))
	fmt.Printf(`Id: %v (Key: %v)
Name: %v
Recorded: %s (racelogger: %s)
Track: %v
Session begin: %s %.0f
Session end: %s %.0f
Race begin (UTC): %s (%d)
`,
		e.Id, e.EventKey,
		e.Name,
		recDate.Format("2006-01-02 15:04"), e.Data.Info.RaceloggerVersion,
		e.Data.Info.TrackDisplayName,
		minSession.String(), e.Data.ReplayInfo.MinSessionTime,
		maxSession.String(), e.Data.ReplayInfo.MaxSessionTime,
		time.Unix(int64(e.Data.ReplayInfo.MinTimestamp), 0).Format("2006-01-02 15:04"), int64(e.Data.ReplayInfo.MinTimestamp))
}

func printJson(e *internal.Event) {
	var s string = ""
	if internal.JsonPretty {
		jsonData, _ := json.MarshalIndent(e, "", "  ")
		s = string(jsonData)
	} else {
		jsonData, _ := json.Marshal(e)
		s = string(jsonData)
	}
	fmt.Printf("%s\n", s)
}
