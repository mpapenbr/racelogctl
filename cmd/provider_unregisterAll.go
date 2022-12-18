/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
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
	dpc := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	providers, err := pc.ProviderList()
	if err != nil {
		log.Fatalf("Error reading provider list: %v\n", err)
	}
	for _, e := range providers {
		dpc.UnregisterProvider(e.EventKey)
	}

}
