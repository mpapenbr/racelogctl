/*
Copyright © 2022 Markus Papenbrock
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// eventCmd represents the event command
var eventCmd = &cobra.Command{
	Use:   "event",
	Short: "Commands regarding events. ",
	// Run: func(cmd *cobra.Command, args []string) {
	// 	fmt.Println("event called")
	// },
}

func init() {
	rootCmd.AddCommand(eventCmd)

}
