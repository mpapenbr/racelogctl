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
	"racelogctl/internal"
	"racelogctl/wamp"

	"github.com/spf13/cobra"
)

// stressCmd represents the stress command
var stressCmd = &cobra.Command{
	Use:   "stress",
	Short: "This command is used for stress testing the application",
	Long:  "",
}

func init() {
	rootCmd.AddCommand(stressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// stressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	stressCmd.PersistentFlags().IntVarP(&internal.Worker, "worker", "w", 1, "Number of workers to use")
}

// helper functions here

func raceLoggerVersion(e *internal.Event) bool {
	return len(e.Data.Info.RaceloggerVersion) > 0
}
func isMinSessionLength(e *internal.Event, minSessionLengthMinutes int) bool {
	return (e.Data.ReplayInfo.MaxSessionTime - e.Data.ReplayInfo.MinSessionTime) > float64(minSessionLengthMinutes*60)
}

func computeAvailableEvents(minSessionLengthMinutes int) []internal.Event {
	availableEvents := []internal.Event{}
	allEvents := wamp.GetEventList()
	for _, event := range allEvents {
		validSource := raceLoggerVersion(&event) && isMinSessionLength(&event, minSessionLengthMinutes)
		if validSource {
			availableEvents = append(availableEvents, event)
			printEventOverview(&event)
		}
	}
	return availableEvents
	// fmt.Printf("%v\n", availableEvents)
}
