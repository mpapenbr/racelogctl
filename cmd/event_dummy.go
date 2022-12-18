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
	"log"
	"racelogctl/internal"
	"racelogctl/wamp"
	"strconv"

	"github.com/blang/semver/v4"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// dummyCmd represents the dummy command
var dummyCmd = &cobra.Command{
	Use:   "dummy",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		// println("in event_info preRunE")
		bindFlags(cmd, viper.GetViper())
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		eventId, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			return
		} else {
			dummy(eventId)
		}

	},
	Args: cobra.ExactArgs(1),
}

func init() {
	eventCmd.AddCommand(dummyCmd)

}

func dummy(eventId int) {
	pc := wamp.NewPublicClient(internal.Url, internal.Realm)
	event, _ := pc.GetEvent(eventId)

	sourceVersion := event.Data.Info.RaceloggerVersion
	if len(sourceVersion) == 0 {
		sourceVersion = "0.0.0"
	}
	fmt.Printf("recieving raceloggerVersion: %s\n", sourceVersion)
	v, err := semver.Parse(sourceVersion)
	if err != nil {
		log.Fatalf("Error parsing version: %v\n", err)
	}
	r, err := semver.ParseRange(">=0.4.4")
	if err != nil {
		log.Fatalf("Error parsing range: %v\n", err)
	}
	fmt.Printf("hasCarData: %v\n", r(v))

}
