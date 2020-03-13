package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
)

var cfg cfgCmd
var rootCmd = &cobra.Command{
	Use:     "brewman",
	Version: "0.2",
	Short:   "A command-line brew application with a beer in my hand",
	Long: `When you brew your own beer the time is comming to do it with some more cyberpunk stuff.
brewman controls multiple pods with different types of recipes and tasks.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Root().Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type cfgCmd struct {
	file   string
	format string
	debug  bool
	recipe *os.File
	kettle string
}

func init() {
	var cfg cfgCmd
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfg.file, "config", "", "config file (default is $HOME/.brewman.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.format, "format", "text", "output format: json,text")
	rootCmd.PersistentFlags().BoolVarP(&cfg.debug, "debug", "d", false, "debug messages")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(getCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfg.file != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfg.file)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".brewman" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".brewman")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
