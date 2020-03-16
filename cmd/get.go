package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get basic output",
	Long:  `get gives you multiple informations in your prefered output format `,
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		//state.Id = viper.GetInt("jobid")
		fmt.Println("get called")
	},
}

// no json support format
var getConfig = &cobra.Command{
	Use:   "config",
	Short: "get config",
	Long:  `get gives you multiple informations in your prefered output format `,
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cfg.conf)
	},
}

// no json support format
var getRecipe = &cobra.Command{
	Use:   "recipe",
	Short: "get recipe",
	Long:  `get gives you multiple informations in your prefered output format `,
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cfg.recipe)
	},
}

func init() {
	getCmd.AddCommand(getConfig)
	getCmd.AddCommand(getRecipe)
}
