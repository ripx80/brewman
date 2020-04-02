package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	"github.com/spf13/cobra"
)

// hotwaterCmd represents the hotwater command
var cookCmd = &cobra.Command{
	Use:   "cook",
	Short: "start cooking",
	Long:  `start cooking set in recipe and hold temperatur set in config`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initRecipe()
		initPods()
		initChan()
	},
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			defer cfg.wg.Done()
			cfg.wg.Add(1)
			if err := cfg.pods.cooker.Run(); err != nil {
				log.WithFields(log.Fields{
					"kettle": "cooker",
					"error":  err,
				}).Error("kettle func failed")
			}
			cfg.done <- struct{}{}
		}()
		handle()
	},
}

var cookMetric = &cobra.Command{
	Use:   "metric",
	Short: "get metric from cooker",
	Long:  `get metric of hotwater pod`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := json.Marshal(cfg.pods.cooker.Metric())
		if err != nil {
			fmt.Println(err)
			exit.Exit(1)
		}
		fmt.Println(string(out))
	},
}

func init() {
	cookCmd.AddCommand(cookMetric)
}
