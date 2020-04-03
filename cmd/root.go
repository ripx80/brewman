package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	logrusBrave "github.com/ripx80/brave/log/logrus"
	"github.com/ripx80/brave/work"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/ripx80/recipe"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/validator.v2"
)

const (
	version    string = "0.2"
	name       string = "brewman"
	configfile string = "brewman.yaml"
	logfile    string = "brewman.log"
)

var rootCmd = &cobra.Command{
	Use:     name,
	Version: version,
	Short:   "A command-line brew application with a beer in your hand",
	Long: `When you brew your own beer the time is comming to do it with some more cyberpunk stuff.
controls multiple pods with different types of recipes and tasks.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		initRecipe()
		initChan() // chan before pods!
		initPods()
		cfg.ui = true
		//set logfile for ui or use a in memory logger
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("canot open log file")
		}
		logr := getLogrus()
		logr.SetOutput(f)

		log.Set(logrusBrave.Configured(logr))

		go func() {
			defer cfg.wg.Done()
			cfg.wg.Add(1)
			if err := view(); err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("view run func failed")
			}
			cfg.done <- struct{}{}
		}()
		handle()
	},
}

var cfg cfgCmd

/*Execute the rootCmd*/
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

	pods podsup
	conf *config.Config

	wg *sync.WaitGroup
	stop,
	done chan struct{}
	signals     chan os.Signal
	confirm     chan pod.Quest
	confirmFunc func() error

	ui bool
}

type podsup struct {
	hotwater *pod.Pod
	cooker   *pod.Pod
	masher   *pod.Pod
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)
	rootCmd.PersistentFlags().StringVar(&cfg.file, "config", "", fmt.Sprintf("config file (default is $HOME/.%s)", configfile))
	rootCmd.PersistentFlags().StringVar(&cfg.format, "format", "text", "output format: json,text")
	rootCmd.PersistentFlags().BoolVarP(&cfg.debug, "debug", "d", false, "debug messages")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(hotwaterCmd)
	rootCmd.AddCommand(mashCmd)
	rootCmd.AddCommand(cookCmd)
	rootCmd.AddCommand(controlCmd)
}

func initChan() {
	// threads, add data chan, error chan
	cfg.signals = make(chan os.Signal, 1)
	cfg.stop = make(chan struct{})
	cfg.confirm = make(chan pod.Quest)
	cfg.done = make(chan struct{}, 2) // buff because after closing stop nobody will recive
	signal.Notify(cfg.signals, syscall.SIGINT, syscall.SIGTERM)
	cfg.wg = new(sync.WaitGroup)
}

func confirmConsole() error {
	var quest pod.Quest
	opt := "use y/n"
	for {
		select {
		case <-cfg.stop:
			return nil
		case quest = <-cfg.confirm:
			var response string
			fmt.Printf("%s (Y/n)", quest.Msg)
			l, err := fmt.Scan(&response)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
				}).Error("canot read response")
			}
			if l > 3 {
				log.Info(opt) // was a warning
				continue
			}
			response = strings.ToLower(response)
			switch response {
			case "n":
				cfg.confirm <- pod.Quest{Msg: "n", Asw: false}
			default:
				cfg.confirm <- pod.Quest{Msg: "y", Asw: true}
			}

		}
	}
}

func getLogrus() *logrus.Logger {
	logr := logrus.New()
	if cfg.debug {
		logr.SetLevel(logrus.DebugLevel)
	}

	if cfg.format == "json" {
		jf := logrus.JSONFormatter{}
		//jf.PrettyPrint = true
		logr.SetFormatter(&jf)
	}
	return logr
}

func initLogger() {
	// use a wrapper for the WithFields func
	log.Set(logrusBrave.Configured(getLogrus()))
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
	cfg.pods.hotwater = pod.New(hotwater, cfg.recipe)
	cfg.pods.hotwater.Hotwater(cfg.conf.Global.HotwaterTemperatur)

	masher, err := getKettle(cfg.conf.Masher)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "masher",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.masher = pod.New(masher, cfg.recipe)
	cfg.pods.masher.Mash(cfg.conf.Global.HoldTemperatur, cfg.confirm)

	cooker, err := getKettle(cfg.conf.Cooker)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"kettle": "cooker",
		}).Error("init failed")
		exit.Exit(1)
	}
	cfg.pods.cooker = pod.New(cooker, cfg.recipe)
	cfg.pods.cooker.Cook(cfg.conf.Global.CookingTemperatur)
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
		cfg.file = configfile
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

func handle() {
	for {
		select {
		case <-cfg.signals:
		case <-cfg.stop:
		case <-cfg.done:
		}
		close(cfg.stop)
		cfg.pods.hotwater.Stop()
		cfg.pods.masher.Stop()
		cfg.pods.cooker.Stop()
		//cfg.wg.Wait() check this when finish, confirm hangs
		work.WaitTimeout(cfg.wg, 1*time.Second) // wait for all workers with timeout
		exit.Exit(0)                            // check the return
	}
}
