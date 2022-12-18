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

	"github.com/spf13/cobra"
)

// unregisterCmd represents the unregister command
var unregisterCmd = &cobra.Command{
	Use:   "unregister <eventKey>",
	Short: "Unregisters a race provider by its eventKey",
	Long: `Unregisters a race provider by its eventKey.

This command may be used if the race provider was terminated before the race ended. 
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {

			unregisterEventProvider(args[0])

		} else {
			fmt.Println("requires an event key")
		}
	},
}

func init() {
	providerCmd.AddCommand(unregisterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unregisterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unregisterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	unregisterCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")
}

func unregisterEventProvider(eventKey string) {
	dpc := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	dpc.UnregisterProvider(eventKey)
}
