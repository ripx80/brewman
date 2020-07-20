package cmd

import (
	log "github.com/ripx80/brave/log/logger"
	"github.com/spf13/cobra"
)

// controlCmd represents the control command
var controlCmd = &cobra.Command{
	Use:   "control",
	Short: "hardware control",
	Long:  `turn on or off pods and validate hardware`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		initRecipe()
		initPods()
	},
}

var offCmd = &cobra.Command{
	Use:   "off",
	Short: "turn off pods",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "hotwater",
		}).Debug("control all off")
		cfg.pods.hotwater.Kettle.Off()
		cfg.pods.cooker.Kettle.Off()
		cfg.pods.masher.Kettle.Off()
	},
}

var offHotwater = &cobra.Command{
	Use:   "hotwater",
	Short: "turn off hotwater pod",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "hotwater",
		}).Debug("control off")
		cfg.pods.hotwater.Kettle.Off()
	},
}

var offCooker = &cobra.Command{
	Use:   "cooker",
	Short: "turn off cooker pod",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "cooker",
		}).Debug("control off")
		cfg.pods.cooker.Kettle.Off()
	},
}

var offMasher = &cobra.Command{
	Use:   "masher",
	Short: "turn off masher pod",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "masher",
		}).Debug("control off")
		cfg.pods.masher.Kettle.Off()
	},
}

var onCmd = &cobra.Command{
	Use:   "on",
	Short: "turn on pods",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "hotwater",
		}).Debug("control all on")
		cfg.pods.hotwater.Kettle.On()
		cfg.pods.cooker.Kettle.On()
		cfg.pods.masher.Kettle.On()
	},
}

var onHotwater = &cobra.Command{
	Use:   "hotwater",
	Short: "turn on hotwater pod",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "hotwater",
		}).Debug("control on")
		cfg.pods.hotwater.Kettle.On()
	},
}

var onMasher = &cobra.Command{
	Use:   "masher",
	Short: "turn on masher pod",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "masher",
		}).Debug("control on")
		cfg.pods.masher.Kettle.On()
	},
}

var onCooker = &cobra.Command{
	Use:   "cooker",
	Short: "turn on cooker kettle",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		log.WithFields(log.Fields{
			"kettle": "cooker",
		}).Debug("control on")
		cfg.pods.cooker.Kettle.On()
	},
}

var validateCmd = &cobra.Command{
	Use:       "validate",
	Short:     "validate contorl units with a test programm",
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"hotwater", "cook", "mash"},
	Run: func(cmd *cobra.Command, args []string) {
		cfg.pods.hotwater.Validate(cfg.pods.hotwater.Kettle.Metric().Temp)
		cfg.pods.masher.Validate(cfg.pods.masher.Kettle.Metric().Temp)
		cfg.pods.cooker.Validate(cfg.pods.cooker.Kettle.Metric().Temp)
		// can be run all together in the future
		log.WithFields(log.Fields{
			"kettle": "hotwater",
		}).Debug("validate")
		cfg.pods.hotwater.Run()
		log.WithFields(log.Fields{
			"kettle": "masher",
		}).Debug("validate")
		cfg.pods.masher.Run()
		log.WithFields(log.Fields{
			"kettle": "cooker",
		}).Debug("validate")
		cfg.pods.cooker.Run()
	},
}

func init() {
	controlCmd.AddCommand(offCmd)
	offCmd.AddCommand(offHotwater)
	offCmd.AddCommand(offMasher)
	offCmd.AddCommand(offCooker)
	controlCmd.AddCommand(onCmd)
	onCmd.AddCommand(onHotwater)
	onCmd.AddCommand(onMasher)
	onCmd.AddCommand(onCooker)
	controlCmd.AddCommand(validateCmd)
}
