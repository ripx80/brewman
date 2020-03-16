package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mashCmd represents the mash command
var mashCmd = &cobra.Command{
	Use:   "mash",
	Short: "start the mash procedure",
	Long:  `start the mash procedure given in recipe`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mash called")
	},
}

var mashMetric = &cobra.Command{
	Use:   "metric",
	Short: "get mash pod metric",
	Long:  `get metrics of mash mod`,
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mash metric called")
	},
}

var mashRest = &cobra.Command{
	Use:   "rest",
	Short: "mash the given rest",
	Long:  `mash the given rest. after finishing stop mashing`,
	Args:  cobra.MinimumNArgs(1),
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mash metric called %s\n", args[0])
	},
}

// add continue flag on rest
func init() {
	mashCmd.AddCommand(mashMetric)
	mashCmd.AddCommand(mashRest)
}
