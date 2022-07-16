/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"racelogctl/internal"
	"racelogctl/wamp"

	"github.com/spf13/cobra"
)

// unregisterAllCmd represents the unregisterAll command
var unregisterAllCmd = &cobra.Command{
	Use:   "unregisterAll",
	Short: "Unregisters all current providers",

	Run: func(cmd *cobra.Command, args []string) {
		unregisterAll()
	},
}

func init() {
	providerCmd.AddCommand(unregisterAllCmd)
	unregisterAllCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")
}

func unregisterAll() {
	wamp.ConsumeProviders(func(e *internal.Event, i int) bool {
		wamp.UnregisterProvider(e.EventKey)
		return true
	})
}
