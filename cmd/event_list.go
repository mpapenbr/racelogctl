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
	Short: "Lists all available events",
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
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	allEvents, _ := pc.GetEventList()
	for _, e := range allEvents {
		printEventOverview(e)

	}

}

func printEventOverview(e *internal.Event) {
	fmt.Printf("%s\n", composeEventOverview(e))
}

func composeEventOverview(e *internal.Event) string {
	recDate, err := time.Parse("2006-01-02T15:04:05", e.RecordDate)
	if err != nil {
		fmt.Println(err)
	}
	return fmt.Sprintf("Id: %3d Date: %s Name: %v", e.Id, recDate.Format("2006-01-02"), e.Name)
}
