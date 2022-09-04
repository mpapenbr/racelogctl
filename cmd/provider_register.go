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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"

	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register an event (Note: use only for development)",
	Long: `This command performs the register procedure of an event. 
The registration is usually performed by the racelogger.
For debugging purpose this command may be used to initialize the backend in a similar manner.`,

	Run: func(cmd *cobra.Command, args []string) {
		register()
	},
}

func init() {
	providerCmd.AddCommand(registerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")
	registerCmd.Flags().StringVarP(&internal.SampleFile, "sample", "s", "", "Sample event file")
	registerCmd.Flags().StringVarP(&internal.EventName, "name", "n", "", "Event name for registration (default: Sample-YYYY-MM-DD-HH-MM)")
	registerCmd.Flags().StringVarP(&internal.EventKey, "key", "k", "", "Event key for registration")
	registerCmd.Flags().StringVarP(&internal.EventDescription, "description", "d", "", "Event description")
	registerCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func register() {
	registerMsg := internal.RegisterMessage{}
	if len(internal.SampleFile) > 0 {
		event := &internal.Event{}
		json.Unmarshal(readSampleFile(internal.SampleFile), &event)
		registerMsg.EventKey = event.EventKey
		registerMsg.Manifests = event.Data.Manifests
		registerMsg.Info = event.Data.Info

	} else {
		// find some dummy values here
	}
	if len(internal.EventKey) > 0 {
		registerMsg.EventKey = internal.EventKey
	}
	if len(internal.EventName) > 0 {
		registerMsg.Info.Name = internal.EventName
	}

	wamp.RegisterProvider(registerMsg)
}

func readSampleFile(filename string) []byte {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer jsonFile.Close()
	ret, _ := ioutil.ReadAll(jsonFile)
	return ret
}
