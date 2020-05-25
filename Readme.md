# Brewman

## Setup

```bash
make; ./bin/brewman
```

You will find your executables in bin

- brewman - Brewery Helper
- m3c - Maische-Malz-und-Mehr recipe converter
- m3d - Maische-Malz-und-Mehr recipe downloader

## Features

describe comming soon

## Not Supported

- At the moment no comments supported in recipe or config file. If you set recipe the config will be overwritten. Parser dont support preserving comments
- Recipe Parser not support Dekoktion and has only kg as unit for malt
- Recipe Parser not support Value Ranges in Recipe (Gaertemperatur: 20-24) and Unit "ml/l".
- Periph.io onewire, I'm not gonna be able to get it to work. Sorry I used the files in sys

## Roadmap

0.1

- [x] read config
- [x] save config
- [x] convert, save and read recipe
- [x] build cmd tool for convert m3 recipes
- [x] change dep to go modules
- [x] add dummy mode
- [x] TempWatcher, Temp.Get(), targetTemp, Hold Time
- [x] one log handler with logrus. in lib also
- [x] os signal handling. Turn off all controls
- [x] Control all on/off switch
- [x] convert all yaml packages to one (recipe, config)
- [x] build Makefile

0.2

- [x] jump to rast
- [x] lib improvements reduce dependencies
- [x] remove in mash config dep
- [x] document exported funcs, check private
- [x] weird control. if something happen crazy (heater on -> temp down) log warnings
- [x] correct os.Exit codes
- [x] add data channel to lib(not abroved)
- [x] validate (try run, with a demo recipe) (check all pins, heat cattle for 2*C)
- [x] add 433mh control unit
- [x] improve cmd parsing, change to cobra (check structure)
- [x] removed confirmations

0.3

- [x] use interactive mode - use a fancy terminal to show temps (tview)
- [x] switch between pods
- [x] add confirm in tview from channel
- [x] correct color scheme ui
- [x] recipies: calculate other water/size informations
- [ ] improve code quality and simplicity
- [ ] add explicit usage in subcommands of cmds
- [ ] bug: out in json, yaml or text
-  bug: maisher run after confirm no focus to main window to go back. get n and y 

0.4

- [ ] ui set the reciept per pod
- [ ] add continue flag on rest
- [ ] add all test files, use testify and testify/mocking
- [ ] Interface: Sensors (sensors.temp, sensors.flow, sensors.*), Control: 433GHz, Relais
- [ ] check recipe, all nessecary values set. no negative and creepy values? (tobi)

0.5

- [ ] grab Status from api
- [ ] Metric Exporter: Prometheus, NodeExporter and internal
- [ ] grab Metrics from Prometheus

## Units

- all weight sizes are gramm
- all time units are minutes
