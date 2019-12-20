package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/ripx80/brewman/config"
	"github.com/ripx80/brewman/pkgs/brew"
	"github.com/ripx80/brewman/pkgs/recipe"
	log "github.com/sirupsen/logrus"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
	validator "gopkg.in/validator.v2"
)

type ConfigCmd struct {
	configFile   string
	outputFormat string
	debug        *bool
	recipe       *os.File
	kettle       string
}

func absolutePath(fp *os.File) (string, error) {
	return filepath.Abs(fp.Name())
}

func confirm(msg string) bool {
	opt := "use y/n"
	for {
		var response string
		log.Info(msg)
		l, err := fmt.Scan(&response)
		if err != nil {
			log.Error(err)
		}
		if l > 3 {
			log.Warn(opt)
			continue
		}
		response = strings.ToLower(response)
		switch response {
		case "y":
			fallthrough
		case "yes":
			return true
		case "n":
			fallthrough
		case "no":
			return false
		default:
			log.Warn(opt)
		}

	}

}

func goExit(signals chan os.Signal) {
	signals <- syscall.SIGINT // stops all threats and do a cleanup
	select {}
}

func main() {

	// config for cmd flags
	cfg := ConfigCmd{}

	a := kingpin.New("brewman", "A command-line brew application")
	a.Version("1.0")
	a.HelpFlag.Short('h')
	a.Author("https://github/ripx80")

	a.Flag("config.file", "brewman configuration file path.").
		StringVar(&cfg.configFile)

	a.Flag("output.format", "output format").
		HintOptions("text", "json").
		Default("text").StringVar(&cfg.outputFormat)

	cfg.debug = a.Flag("output.debug", "Enable debug mode.").Short('v').Bool()

	sc := a.Command("get", "get basic output")
	sc.Command("config", "output current config")
	sc.Command("recipe", "output control information")

	// save in config file
	sc = a.Command("set", "set values")
	sc.Command("config", "save current config to file")
	sr := sc.Command("recipe", "set recipe to brew")
	sr.Arg("filename", "file of the recipe").Required().FileVar(&cfg.recipe)

	sc = a.Command("mash", "mash brew steps")
	sc.Command("start", "start the mash precedure")

	sc = a.Command("hotwater", "make hotwater in kettle")
	sc.Command("start", "start the hotwater precedure")

	sc = a.Command("cook", "cooking your stuff")
	sc.Command("start", "start the cooking precedure")

	sc = a.Command("control", "control hardware")
	sr = sc.Command("off", "stop all actions")
	sr.Arg("kettle", "stop only actions on kettle (hotwater, masher, cooker)").HintOptions("hotwater", "masher", "cooker").StringVar(&cfg.kettle)

	sr = sc.Command("on", "turn on all")
	sr.Arg("kettle", "turn on kettle (hotwater, masher, cooker)").HintOptions("hotwater", "masher", "cooker").StringVar(&cfg.kettle)

	_, err := a.Parse(os.Args[1:])
	if err != nil {
		log.Error("Error parsing commandline arguments: ", err)
		a.Usage(os.Args[1:])
	}

	// default config
	configFile, _ := config.Load("")

	if cfg.configFile != "" {
		configFile, err = config.LoadFile(cfg.configFile)
		if err != nil {
			log.Error("canot load configuration file: ", err)
			os.Exit(1)
		}

	} else {
		cfg.configFile = "brewman.yaml"
		if _, err := os.Stat(cfg.configFile); err == nil {
			configFile, err = config.LoadFile(cfg.configFile)
			if err != nil {
				log.Error("canot load config file: ", err)
				os.Exit(1)
			}
		}
	}

	if *cfg.debug {
		log.SetLevel(log.DebugLevel)
	}

	if cfg.outputFormat == "json" {
		jf := log.JSONFormatter{}
		//jf.PrettyPrint = true
		log.SetFormatter(&jf)
	}

	if err := validator.Validate(configFile); err != nil {
		log.Error("Config file validation failed: ", err)
	}

	// threads, add data chan, error chan
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	stop := make(chan struct{})
	kettle := &brew.Kettle{} //need for cleanup

	// signal handler
	go func(kettle *brew.Kettle) {
		defer os.Exit(0)
		select {
		case <-signals:
		case <-stop:
		}
		close(stop)
		log.Info("cleanup in controller threat")
		kettle.Cleanup()
		log.Info("go exit")
	}(kettle)

	switch kingpin.MustParse(a.Parse(os.Args[1:])) {

	case "set config":
		configFile.Save(cfg.configFile)
		fallthrough

	case "get config":
		log.Info(fmt.Sprintf("\n%s\n%s", cfg.configFile, configFile))

	case "set recipe":
		configFile.Recipe.File, err = absolutePath(cfg.recipe)
		if err != nil {
			log.Error("set recipe error: ", err)
		}
		configFile.Save(cfg.configFile)
		fallthrough

	case "get recipe":
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
		log.Info(fmt.Sprintf("\n%s\n%s", configFile.Recipe, recipe))

	case "hotwater start":

		if err = kettle.Init(configFile.Hotwater); err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}
		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Infof("main water: %f -->  grouting: %f", recipe.Water.MainCast, recipe.Water.Grouting)

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := kettle.TempIncreaseTo(stop, configFile.Global.HotwaterTemperatur); err != nil {
				log.Error(err)
				goExit(signals)
			}
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := kettle.TempHolder(stop, configFile.Global.HotwaterTemperatur, 0); err != nil {
				log.Error(err)
				goExit(signals)
			}
		}()
		wg.Wait()
		goExit(signals)

	case "mash start":
		if err = kettle.Init(configFile.Masher); err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Info("mash information: ", recipe.Mash)

		if !confirm("start mashing? <y/n>") {
			goExit(signals)
		}

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}

		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := kettle.TempIncreaseTo(stop, recipe.Mash.InTemperatur); err != nil {
				log.Error(err)
				goExit(signals)
			}
		}()
		wg.Wait()

		if !confirm("malt added? continue? <y/n>") {
			goExit(signals)
		}

		for num, rast := range recipe.Mash.Rests {
			log.Infof("Rast %d: Time: %d Temperatur:%f\n", num, rast.Time, rast.Temperatur)

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := kettle.TempIncreaseTo(stop, rast.Temperatur); err != nil {
					log.Error(err)
					goExit(signals)
				}
			}()
			wg.Wait()

			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := kettle.TempHolder(stop, rast.Temperatur, time.Duration(rast.Time*60)*time.Second); err != nil {
					log.Error(err)
					goExit(signals)
				}
			}()
			wg.Wait()
		}

		log.Info("Mashing finished successful")
		goExit(signals)

	case "cook start":
		if err = kettle.Init(configFile.Cooker); err != nil {
			log.Fatal("Failed to init Kettle:", err)
		}

		recipe, err := recipe.LoadFile(configFile.Recipe.File, &recipe.Recipe{})
		if err != nil {
			log.Fatal(err)
		}

		log.Info("using recipe: ", recipe.Global.Name)
		log.Info("cook information: ", recipe.Cook)

		if !confirm("start cooking? <y/n>") {
			goExit(signals)
		}

		if kettle.Agitator != nil && !kettle.Agitator.State() {
			kettle.Agitator.On()
		}
		wg := new(sync.WaitGroup)
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := kettle.TempIncreaseTo(stop, configFile.Global.CookingTemperatur); err != nil {
				log.Error(err)
				goExit(signals)
			}
		}()
		wg.Wait()

		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := kettle.TempHolder(stop, configFile.Global.CookingTemperatur, time.Duration(recipe.Cook.Time*60)*time.Second); err != nil {
				log.Error(err)
				goExit(signals)
			}
		}()
		wg.Wait()

		log.Info("Cooking finished successful")
		goExit(signals)

	case "control off":
		switch cfg.kettle {
		case "hotwater":
			log.Info("stop hotwater")
			if err = kettle.Init(configFile.Hotwater); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		case "masher":
			log.Info("stop masher")
			if err = kettle.Init(configFile.Masher); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		case "cooker":
			log.Info("stop cooker")
			if err = kettle.Init(configFile.Cooker); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		default:
			log.Info("stop all actions and cleanup")
			if err = kettle.Init(configFile.Hotwater); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
			kettle.Cleanup()

			if err = kettle.Init(configFile.Masher); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
			kettle.Cleanup()

			if err = kettle.Init(configFile.Cooker); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		}
		goExit(signals)
	case "control on":
		switch cfg.kettle {
		case "hotwater":
			log.Info("turn on hotwater")
			if err = kettle.Init(configFile.Hotwater); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		case "masher":
			log.Info("turn on masher")
			if err = kettle.Init(configFile.Masher); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		case "cooker":
			log.Info("turn on cooker")
			if err = kettle.Init(configFile.Cooker); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		default:
			log.Info("turn all kettle on")
			if err = kettle.Init(configFile.Hotwater); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
			kettle.On()

			if err = kettle.Init(configFile.Masher); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
			kettle.On()

			if err = kettle.Init(configFile.Cooker); err != nil {
				log.Fatal("Failed to init Kettle:", err)
			}
		}
		kettle.On()
	}

}
