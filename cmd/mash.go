package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	"github.com/spf13/cobra"
)

// mashCmd represents the mash command
var mashCmd = &cobra.Command{
	Use:   "mash",
	Short: "start the mash procedure",
	Long:  `start the mash procedure given in recipe`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initRecipe()
		initPods()
		initChan()
		cfg.pods.masher.Mash(cfg.conf.Global.HoldTemperatur)
	},
	Run: func(cmd *cobra.Command, args []string) {
		go func() {
			defer cfg.wg.Done()
			cfg.wg.Add(1)
			if err := cfg.pods.masher.Run(); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
					"error":  err,
				}).Error("kettle func failed")
			}
			cfg.done <- struct{}{}
		}()
		handle()
	},
}

var mashMetric = &cobra.Command{
	Use:   "metric",
	Short: "get mash pod metric",
	Long:  `get metrics of mash mod`,
	Run: func(cmd *cobra.Command, args []string) {
		out, err := json.Marshal(cfg.pods.masher.Metric())
		if err != nil {
			fmt.Println(err)
			exit.Exit(1)
		}
		fmt.Println(string(out))
	},
}

var mashRest = &cobra.Command{
	Use:   "rest",
	Short: "mash the given rest",
	Long:  `mash the given rest. after finishing stop mashing`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		rastNum, err := strconv.Atoi(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"kettle": "masher",
				"error":  err,
			}).Error("wrong argument")
			exit.Exit(1)
		}
		if rastNum > 8 || rastNum <= 0 {
			log.WithFields(log.Fields{
				"kettle": "masher",
				"error":  err,
			}).Error("rast number out of range [1-8]")
			exit.Exit(1)
		}

		if len(cfg.recipe.Mash.Rests) < rastNum {
			log.WithFields(log.Fields{
				"kettle": "masher",
				"error":  err,
			}).Error("rast number not in recipe")

		}
		cfg.pods.masher.MashRast(rastNum - 1) // set defined task with steps

		go func() {
			defer cfg.wg.Done()
			cfg.wg.Add(1)
			if err := cfg.pods.masher.Run(); err != nil {
				log.WithFields(log.Fields{
					"kettle": "masher",
					"error":  err,
				}).Error("kettle func failed")
			}
			cfg.done <- struct{}{}
		}()
		handle()
	},
}

func init() {
	mashCmd.AddCommand(mashMetric)
	mashCmd.AddCommand(mashRest)
}
