package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// hotwaterCmd represents the hotwater command
var hotwaterCmd = &cobra.Command{
	Use:   "hotwater",
	Short: "start heating hotwater pod",
	Long:  `start heating hotwater pod and hold the temperatur set in config for hotwater`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hotwater called")
	},
}

var hotwaterMetric = &cobra.Command{
	Use:   "metric",
	Short: "get hotwater pod metric",
	Long:  `get metric of hotwater pod`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hotwater metric called")
	},
}

func init() {
	hotwaterCmd.AddCommand(hotwaterMetric)
}
