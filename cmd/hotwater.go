package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	"github.com/spf13/cobra"
)

var hotwaterCmd = &cobra.Command{
	Use:   "hotwater",
	Short: "start heating hotwater pod",
	Long:  `start heating hotwater pod and hold the temperatur set in config for hotwater`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initRecipe()
		initPods()
		initChan()
	},
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			defer cfg.wg.Done()
			cfg.wg.Add(1)
			if err := cfg.pods.hotwater.Run(); err != nil {
				log.WithFields(log.Fields{
					"kettle": "hotwater",
					"error":  err,
				}).Error("kettle func failed")
			}
			cfg.done <- struct{}{}
		}()
		handle()
	},
}

// only json supported at the moment
var hotwaterMetric = &cobra.Command{
	Use:   "metric",
	Short: "get hotwater pod metric",
	Run: func(cmd *cobra.Command, args []string) {
		out, err := json.Marshal(cfg.pods.hotwater.Metric())
		if err != nil {
			fmt.Println(err)
			exit.Exit(1)
		}
		fmt.Println(string(out))
	},
}

func init() {
	hotwaterCmd.AddCommand(hotwaterMetric)
}
