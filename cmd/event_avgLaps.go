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
	"fmt"
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// avgLapsCmd represents the avgLaps command
var avgLapsCmd = &cobra.Command{
	Use:   "avgLaps <eventId>",
	Short: "Get average laptimes by car classes over time for an event",
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
				eventAvgLaps(eventId)
			}
		} else {
			fmt.Println("requires an event id")
		}

	},
}

func init() {
	eventCmd.AddCommand(avgLapsCmd)

	avgLapsCmd.Flags().IntVar(&internal.Interval, "interval", 300, "Interval in seconds")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// avgLapsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// avgLapsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func eventAvgLaps(id int) {

	avgLaps := wamp.GetEventAvgLaps(id, internal.Interval)
	for _, item := range avgLaps {
		fmt.Printf("%.0f (%s): track: %.0f %+v\n", item.Timestamp, time.Unix(int64(item.Timestamp), 0), item.TrackTemp, item.Laptimes)
	}
}
