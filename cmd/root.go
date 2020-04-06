package cmd

/*
## Improve Code Quality

- [ ] Set Recipes per pod in each cmd or ui
- [ ] Check all namings. simplicity?
- [ ] Check dependencies. All necessary?
- [ ] Remove log in libs. use the metrics and a timer (kettle)
- [ ] When not log in libs. need the logger interface?
- [ ] Use Static Errors?
- [ ] Check all routines if ended with stop (ui) use wg group wait to check
- [ ] check the output format flag: text,json,yaml
*/

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ripx80/brave/exit"
	log "github.com/ripx80/brave/log/logger"
	logrusBrave "github.com/ripx80/brave/log/logrus"
	"github.com/ripx80/brave/work"
	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/pod"
	"github.com/ripx80/recipe"
	"github.com/spf13/cobra"
)

const (
	version    string = "0.3"
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
		f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("can not open log file")
		}
		logr := getLogrus()
		logr.SetOutput(f)
		logr.SetLevel(2) // ErrorLevel
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
