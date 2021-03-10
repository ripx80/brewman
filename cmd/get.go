package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ripx80/brave/exit"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get basic output",
	Long:  `get gives you multiple informations in your prefered output format `,
}

// no json support format
var getConfig = &cobra.Command{
	Use:   "config",
	Short: "get config",
	Long:  `outputs the current config`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(cfg.conf)
	},
}

// no json support format
var getRecipe = &cobra.Command{
	Use:   "recipe",
	Short: "get recipe",
	Long:  `output the current recipe`,
	Run: func(cmd *cobra.Command, args []string) {
		initRecipe()
		fmt.Print(cfg.recipe)
	},
}

// no json support format
var getMetrics = &cobra.Command{
	Use:   "metrics",
	Short: "get metrics",
	Long:  `output the current metrics`,
	Run: func(cmd *cobra.Command, args []string) {
		initRecipe()
		//initChan() // chans before pods, need confirm channel
		initPods()
		m := make(map[string]pod.Metric)
		m["hotwater"] = cfg.pods.hotwater.Metric()
		m["masher"] = cfg.pods.hotwater.Metric()
		m["cooker"] = cfg.pods.hotwater.Metric()
		out, err := json.Marshal(m)
		if err != nil {
			fmt.Println(err)
			exit.Exit(1)
		}
		fmt.Println(string(out))
	},
}

func init() {
	getCmd.AddCommand(getConfig)
	getCmd.AddCommand(getRecipe)
	getCmd.AddCommand(getMetrics)
}
