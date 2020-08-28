
# DevBook

## Roadmap

### v0.3

- [x] use interactive mode - use a fancy terminal to show temps (tview)
- [x] switch between pods
- [x] add confirm in tview from channel
- [x] correct color scheme ui
- [x] recipies: calculate other water/size informations
- [x] (I) reduce binary size with ldflags
- [x] (I) ui round TempStart, Temp, TempEnd
- [ ] (I) code quality and simplicity
- [ ] (I) add explicit usage in subcommands of cmds like recipes scale
- [ ] (R) remove brewman.log
- [ ] (I) remove all yaml stuff. Only json is supported
- [ ] (F) out in json or text
- [ ] (I) cmd log output set zstate true in red
- [ ] (I) cmd log ouput add time actual to
- [ ] (I) add successful finish msg to ui
- [ ] (I) remove periph from modules use brian-armstrong/gpio direct
- [ ] (I) handle errors returned to ui. Maybe modal?

#### Bugs

- [x] (F) remove artefacts from modal window after say no to job probe
- [x] (F) jod test eins zu früh, muss beim springen auf abmaischtemperatur angezeigt werden. Im moment wird er bei dem sprung zur letzten Rast angezeigt, in diesem Fall 76C (Simcoe4 Rezept)
- [x] (F) fail count must be more tollerant and reset after time
- [x] (F) control off accept no arguments anymore if you misstype
- [x] (F) after exit does not turn off the heater, it was a logic bug
- [x] (F) nach der letzten Rast beendet das Programm ohne auf Abmaishtemp zu gehen und Schaltet nicht aus!
- [x] (F) maisher run after confirm no jod probe, focus to main window to go back.

#### Bugs - real setup

- [ ] (F) gpio write failed, for first time starting brewman after boot, add retry
         "failed to open gpio 17 direction file for writing" only the first time

### v0.2

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

### v0.1

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

### Untracked features

- [ ] ui set the reciept per pod
- [ ] add continue flag on rest
- [ ] add all test files, use testify and testify/mocking
- [ ] Interface: Sensors (sensors.temp, sensors.flow, sensors.*), Control: 433GHz, Relais
- [ ] check recipe, all nessecary values set. no negative and creepy values? (tobi)
- [ ] grab Status from api
- [ ] Metric Exporter: Prometheus, NodeExporter and internal
- [ ] grab Metrics from Prometheus

### Untracked

- [ ] (I) write tests
- [ ] (I) recipe change in ui
- [ ] (I) prom metric exporter
- [ ] (I) clear docs
- [ ] (I) gocardreport
- [ ] (I) remove unuses subcommands like get called
- [ ] (I) include dummy and overwrite config if dummy is selected
- [ ] (I) include docker
- [ ] (I) include es docker-compose up
- [ ] (I) serach for a way to sort the yaml to have the same view [doc](https://blog.labix.org/2014/09/22/announcing-yaml-v2-for-go#mapslice)
- [ ] (I) cmd metric flag for humans
- [ ] (I) recipe, read, verify (time, temp values)
- [ ] (I) wire proto direct not over file, hanging
- [ ] (I) add contrib dir with additonal soft, like scripts or es stuff

- [ ] (HF) maisher dreht nicht, motor läuft dreht aber nicht, hin und wieder
- [ ] (HF) schleifen beim maishen
- [ ] (HF) hotwater nicht dicht am Hahn
- [ ] (HF) längere kabel für temp sensoren
- [ ] (H) verlängerungskabel
- [ ] (H) Loch im Maisher für Temp Sensor
- [ ] (H) Schaltkasten für Raspberry
- [ ] (HF) vordere schraube maisher am Dockel motor anziehen. Starke vibration

- [ ] rot: hotwater
- [ ] weiß: masher
- [ ] green: cooker

Rezept geändert:
Whirpool von 27g auf 30g erhöht.

## Mesh

- Water in Masher (Open/Close Valves), Count Water Amount, Checks (Open Valves -> Water flows)
- Heat Water in Masher (Heater On/Off, Agitator On/Off), Checks (Heater On -> Temp increase) Sensors: Plate#1, Temp#1
- Malt Fill into Masher (Stop: Confirm), Program for Rests
- Läutern (Stop: confirm, Start Hot Tube Water: confirm), Valves to Cooker, Font Hops/Ingredients

## HotTube

- (P) Water in Hot Tube (Same as 1. Other Valves) Sensors Plate#2, Temp#2
- (P) Heat Water in Hot Tube (Same as 2.)

## Cooking

- Heat Water in Cooker (Temp to cook and hold: 97.5, set with cmd while cooking for precision)
- Reach Cook Temp: Cooking Time, Info about Hops/Ingredients get in! (message with beep and terminal)
- Finish Cook Time: Stop(confirm), Whirlpool Info

## Yaml

Using ghodss/yaml and not the pure go-yaml/yaml.
In short, this library first converts YAML to JSON using go-yaml and then uses json.Marshal and json.Unmarshal to convert to or from the struct. This means that it effectively reuses the JSON struct tags as well as the custom JSON methods MarshalJSON and UnmarshalJSON unlike go-yaml

## Implementation

```go
// Control Elements
struct Valve interface
    func Open()
    func Close()

struct Control interface
    func On() // HIGH, LOW or Programm (433Utils)
    func Off() interface

// Sensors
func (registerSensor(Sensor, func))
func (getHardwareSensors)

struct Sensor interface
    Get()

struct OneWire interface

struct WaterFlow type
struct Thermometer type
struct WaterLevel type

// Outputs
struct Output interface
struct Screen type
struct Terminal type
struct Prometheus type
```

## Ideas

- set the max volume of mesher, cooker, fermenter. So you can auto calc max Outcome with MainCast and Grouting!
- build cmd for convert recipes
- recipe:Download recipe, search recipe

## Config

- Temperatur of HotWaterTube

```yaml
pods:
  # support only one temp at the moment
  hotwater:
    temperatur:
      device: "ds18b20"
      bus: 5
      address: 0x3343839898
    control: 3
  masher:
    control: 4
    agiator: 10
    temperatur:
      device: "ds18b20"
      bus: 4
      address: 0x387839898
  cooker:
    temperatur:
      device: "ds18b20"
      bus: 14
      address: 0x3354839898
    control: 30


sensors:
  hotwater: 4 # add address
  masher: 11
  cooker: 12
  flowin: 13
control:
  heater-water: 29
  heater-mash: 31
  heater-cooker: 32

```

## cmd

```bash

brewman
    set: set things, configs
    get: basic output
    describe: verbose output
    validate: validate things
    logs: dump logs

Output format:
    -o=json
    -o=yaml

Verbosity:

    --v=0 quiet
    --v=1 show extendet infromation
    --v=2 debug information
    --v=3 display all sensor information
    --v=4 display http requests


# recipe
brewman set recipe file #save file in ENV VAR
brewman get recipe # recipe name and additional information
brewman get hops
brewman get rast <number>
brewman get cooking
brewman get fermentation
brewman describe recipe # print full recipe
brewman validate recipe

# config
brewman set config foo=bar
brewman get config
brewman set control.flow=gpio_pin?

# run demo programm
brewman validate # run demo programm with recipe
brewman validate demo # run demo with demo recipe

# brewing
brewman start # full steps
brewman start rast <number> # only the given rast
brewman start cooking # only start at cooking


# sensors
brewman get sensors # print connected Sensor information
brewman set sensor # add sensor to config

# api

brewman start server # run only the api server and wait for instructions
brewman stop # api call to stop
brewman set remote=https://remoteserver:8000

```

## Dependencies

```text
- https://github.com/martinohmann/rfoutlet
- [wiringpi](https://tutorials-raspberrypi.de/wiringpi-installieren-pinbelegung/)
- https://github.com/stianeikeland/go-rpio

- https://gobot.io/documentation/platforms/raspi/
- https://godoc.org/periph.io/x/periph/conn/onewire
- https://github.com/kidoman/embd
```

## Info

- [go-log](https://github.com/go-log/log) logging with interfaces from [this](https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern) article
- [datadog-guide](https://www.datadoghq.com/blog/go-logging/) std logging interface with const values
- [clog](https://github.com/go-clog/clog) use go chan to log back
