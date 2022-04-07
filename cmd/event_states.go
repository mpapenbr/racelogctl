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
				fetchStates(eventId)
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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// stateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func fetchStates(eventId int) {
	fmt.Printf("Fetching %d states beginning at %d\n", internal.Num, internal.From)
	states := wamp.GetStates(eventId, internal.From, internal.Num)
	fmt.Printf("\n---\nresulting states\n")
	for i, entry := range states {
		fmt.Printf("%d, %+v\n", i, entry)
	}
}
