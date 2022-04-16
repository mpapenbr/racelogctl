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
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"

	"github.com/spf13/cobra"
)

// processCmd represents the process command
var processCmd = &cobra.Command{
	Use:   "process EVENT_ID",
	Short: "Reprocess the event",
	Args:  cobra.ExactArgs(1),
	Long: `This command can be used to reprocess an event. 
It will read all existing states from the database and process them again.
The analysis result will be stored in the database.
This feature may be useful when there were bugfixes in the analysis module`,
	Run: func(cmd *cobra.Command, args []string) {
		eventId, err := strconv.Atoi(args[0])
		if err == nil {
			processEvent(eventId)
		} else {
			fmt.Fprintf(os.Stderr, "Not a valid eventId: %v\n", err)
		}
	},
}

func init() {
	eventCmd.AddCommand(processCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// processCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// processCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	processCmd.Flags().StringVarP(&internal.AdminPassword, "admin-password", "p", "", "sets the admin password for this action")
}

func processEvent(eventId int) {

	fmt.Printf("Processing now event %v\n", eventId)
	result := wamp.ProcessEvent(eventId)
	if len(result.Error) > 0 {
		fmt.Printf("Error: %s\n", result.Error)
	} else {
		fmt.Printf("Finished: %s\n", result.Message)
	}

}
