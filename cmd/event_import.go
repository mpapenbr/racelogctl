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
	"bufio"
	"encoding/json"
	"log"
	"os"
	"racelogctl/internal"
	"racelogctl/wamp"

	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import <eventId> <input>",
	Short: "Reads data from a file and sends it the racelogger backend.",
	Long:  `TODO: requirements when to use....`,
	Run: func(cmd *cobra.Command, args []string) {
		importData()
	},
}

func init() {
	eventCmd.AddCommand(importCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// importCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	importCmd.Flags().StringVarP(&internal.Input, "input", "i", "", "Input file containing to states to be imported")
	importCmd.MarkFlagRequired("input")
	importCmd.Flags().StringVarP(&internal.EventKey, "eventKey", "k", "", "Key of the event recieving the data")
	importCmd.MarkFlagRequired("eventKey")
	importCmd.Flags().StringVarP(&internal.DataproviderPassword, "dataprovider-password", "p", "", "sets the Dataprovider password for this action")
}

func importData() {
	file, err := os.Open(internal.Input)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	sender := make(chan internal.State)

	dataprovider := wamp.NewDataProviderClient(internal.Url, internal.Realm, internal.DataproviderPassword)
	defer dataprovider.Close()
	dataprovider.PublishStateFromChannel(internal.EventKey, sender)

	idx := 0
	for scanner.Scan() {
		line := scanner.Text()
		s := internal.State{}
		json.Unmarshal([]byte(line), &s)
		idx++
		sender <- s
		// fmt.Printf("%v\n", s.Payload.Session)
	}
	close(sender)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
