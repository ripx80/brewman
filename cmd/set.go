package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("set called")
	},
}

var setConfig = &cobra.Command{
	Use:   "config",
	Short: "set config",
	Long:  `get gives you multiple informations in your prefered output format `,
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Run: func(cmd *cobra.Command, args []string) {
		cfg.conf.Save(cfg.file)
	},
}

var setRecipe = &cobra.Command{
	Use:   "recipe",
	Short: "set recipe",
	Long:  `get gives you multiple informations in your prefered output format `,
	// PersistentPreRun
	//	PreRun: func(cmd *cobra.Command, args []string) {},
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		cfg.conf.Recipe.File, err = filepath.Abs(args[0])
		if err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"recipe": cfg.conf.Recipe.File,
				"config": cfg.file,
			}).Error("can not get recipe file")
			exit.Exit(1)
		}
		// todo parsing recipe and check content
		if err = cfg.conf.Save(cfg.file); err != nil {
			log.WithFields(log.Fields{
				"error":      err,
				"configFile": cfg.file,
			}).Error("set recipe in configuration")
			exit.Exit(1)
		}
	},
}

var (
	filename string
)

func init() {
	setCmd.AddCommand(setConfig)
	setCmd.AddCommand(setRecipe)
	//setRecipe.Flags().StringVarP(&filename, "filename", "s", "", "Source directory to read from")
	// setRecipe.Flags().StringP("host", "s", "", "export host connect to")
	// viper.BindPFlag("end", blockCmd.Flags().Lookup("end"))

}
