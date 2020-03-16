package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	logrusBrave "github.com/ripx80/brave/log/logrus"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/ripx80/recipe"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
)

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

var cfg cfgCmd

func Execute() {
	defer exit.Safe()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		exit.Exit(1)
	}
}

type cfgCmd struct {
	file       string
	format     string
	debug      bool
	recipeFile string //not used
	recipe     *recipe.Recipe

	pods Pods
	conf *config.Config

	wg *sync.WaitGroup
	stop,
	done chan struct{}
	signals chan os.Signal
}

type Pods struct {
	hotwater *pod.Pod
	cooker   *pod.Pod
	masher   *pod.Pod
}

func init() {
	cobra.OnInitialize(initLogger, initConfig, initRecipe, initPods, initChan)
	rootCmd.PersistentFlags().StringVar(&cfg.file, "config", "", "config file (default is $HOME/.brewman.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.format, "format", "text", "output format: json,text")
	rootCmd.PersistentFlags().BoolVarP(&cfg.debug, "debug", "d", false, "debug messages")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(mashCmd)
	rootCmd.AddCommand(hotwaterCmd)
	rootCmd.AddCommand(controlCmd)
}

func initChan() {
	// threads, add data chan, error chan
	cfg.signals = make(chan os.Signal, 1)
	cfg.stop = make(chan struct{})
	cfg.done = make(chan struct{})
	signal.Notify(cfg.signals, syscall.SIGINT, syscall.SIGTERM)
	cfg.wg = new(sync.WaitGroup)
}

func initLogger() {
	logr := logrus.New()
	if cfg.debug {
		logr.SetLevel(logrus.DebugLevel)
	}

	if cfg.format == "json" {
		jf := logrus.JSONFormatter{}
		//jf.PrettyPrint = true
		logr.SetFormatter(&jf)
	}
	// use a wrapper for the WithFields func
	log.Set(logrusBrave.Configured(logr))
}

func initPods() {
	var err error

	hotwater, err := getKettle(cfg.conf.Hotwater)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "hotwater",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.hotwater = pod.New(hotwater, cfg.recipe, cfg.stop)

	masher, err := getKettle(cfg.conf.Masher)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "masher",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.masher = pod.New(masher, cfg.recipe, cfg.stop)

	cooker, err := getKettle(cfg.conf.Cooker)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "cooker",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.cooker = pod.New(cooker, cfg.recipe, cfg.stop)
}

func initRecipe() {
	var err error
	cfg.recipe, err = recipe.LoadFile(cfg.conf.Recipe.File, &recipe.Recipe{})
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"recipe": cfg.conf.Recipe.File,
		}).Error("init recipe")
		exit.Exit(1)
	}
}

func initConfig() {
	cfg.conf, _ = config.Load("")
	var err error

	if cfg.file != "" {
		cfg.conf, err = config.LoadFile(cfg.file)
		if err != nil {
			fmt.Printf("can not load configuration file: %s\n", cfg.file)
			exit.Exit(1)
		}
	} else {
		cfg.file = "brewman.yaml"
		if _, err := os.Stat(cfg.file); err == nil {
			cfg.conf, err = config.LoadFile(cfg.file)
			if err != nil {
				fmt.Printf("can not load configuration file: %s\n", cfg.file)
				exit.Exit(1)
			}
		}
	}
	if err := validator.Validate(cfg.conf); err != nil {
		fmt.Printf("config file validation failed: %s\n", err)
		exit.Exit(1)
	}
}
