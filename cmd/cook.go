package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// hotwaterCmd represents the hotwater command
var cookCmd = &cobra.Command{
	Use:   "cook",
	Short: "start cooking",
	Long:  `start cooking set in recipe and hold temperatur set in config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cook called")
	},
}

var cookMetric = &cobra.Command{
	Use:   "metric",
	Short: "get metric from cooker",
	Long:  `get metric of hotwater pod`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cook metric called")
	},
}

func init() {
	cookCmd.AddCommand(cookMetric)
}
