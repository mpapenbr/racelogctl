/*
Copyright © 2022 Markus Papenbrock

*/
package cmd

import (
	"fmt"
	"racelogctl/internal"
	"racelogctl/wamp"
	"time"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all available events",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		listEvents()
	},
}

func init() {
	eventCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func listEvents() {
	fmt.Printf("Using Realm %s at %s\n", internal.Realm, internal.Url)
	wamp.GetEvents(func(e *internal.Event, idx int) bool {
		// fmt.Printf("idx: %v e: %+v\n", idx, e)

		printEventOverview(e)
		return true

	})

}

func printEventOverview(e *internal.Event) {
	recDate, err := time.Parse("2006-01-02T15:04:05", e.RecordDate)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Id: %3d Date: %s Name: %v \n", e.Id, recDate.Format("2006-01-02"), e.Name)
}
