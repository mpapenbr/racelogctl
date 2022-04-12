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
	"github.com/spf13/viper"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete EVENT_ID",
	Short: "This will delete an event from the database. Needs admin permissions.",
	Long:  `Note: Be careful! This command will be executed without further confirmation! `,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		eventId, err := strconv.Atoi(args[0])
		if err == nil {
			deleteEvent(eventId)
		} else {
			fmt.Fprintf(os.Stderr, "Not a valid eventId: %v\n", err)
		}
	},
}

func init() {
	// println("deleteCmd.init")
	eventCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deleteCmd.PersistentFlags().StringVarP(&internal.AdminPassword, "admin-password", "p", "", "sets the admin password for this action")
	// deleteCmd.Flags().IntVarP(&eventId, "eventId", "e", -1, "sets the admin password for this action")

	//
	viper.BindPFlag("admin.password", deleteCmd.PersistentFlags().Lookup("admin-password"))

}

func deleteEvent(eventId int) {

	fmt.Printf("Deleting now event %v\n", eventId)
	wamp.DeleteEvent(eventId)
	fmt.Printf("Deleted  event %v\n", eventId)
}
