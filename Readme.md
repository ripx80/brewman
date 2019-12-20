# Brewman

## Features

- Log: logger to console, file, socket: Yes with MultipleWriter.

## Not Supported

- At the moment no comments supported in recipe or config file. If you set recipe the config will be overwritten. Parser dont support preserving comments
- Recipe Parser not support Dekoktion and has only kg as unit for malt
- Recipe Parser not support Value Ranges in Recipe (Gaertemperatur: 20-24) and Unit "ml/l".
- Periph.io onewire, I'm not gonna be able to get it to work. Sorry I used the files in sys

## Roadmap

- Interface: Sensors (sensors.temp, sensors.flow, sensors.*), Control: 433GHz, Relais
- Metric Exporter: Prometheus, NodeExporter and internal
- Tests

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
- [ ] build Makefile
- [ ] convert all yaml packages to one (recipe, config)

0.2

- [ ] try run (check all pins, heat cattle for 2*C)
- [ ] check recipe, all nessecary values set. no negative and creepy values? (tobi)


## Bug Fixes

## Units

- All weight sizes are gramm
- All Time units are Minutes

## general func

- args: parse, check
- config: read, check
- recipe: read, verify (time, temp values)
- sensors: initilize, check, use
- controls: initilize, check, use
- execute recipe plan
- execute metric endpoint

## config

config values of bus is string because periph need a string input: 40, PWM0, GPIO40 are the same

## struct

Valves not supproted yet

(Full Program, All Steps can be independently start from each other)

### Mesh

- Water in Masher (Open/Close Valves), Count Water Amount, Checks (Open Valves -> Water flows)
- Heat Water in Masher (Heater On/Off, Agitator On/Off), Checks (Heater On -> Temp increase) Sensors: Plate#1, Temp#1
- Malt Fill into Masher (Stop: Confirm), Program for Rests
- Läutern (Stop: confirm, Start Hot Tube Water: confirm), Valves to Cooker, Font Hops/Ingredients

### HotTube

- (P) Water in Hot Tube (Same as 1. Other Valves) Sensors Plate#2, Temp#2
- (P) Heat Water in Hot Tube (Same as 2.)

### Cooking

- Heat Water in Cooker (Temp to cook and hold: 97.5, set with cmd while cooking for precision)
- Reach Cook Temp: Cooking Time, Info about Hops/Ingredients get in! (message with beep and terminal)
- Finish Cook Time: Stop(confirm), Whirlpool Info

### Fermentation

- Info About Fermentation

### Yaml

Todo: only use one of the libs.

Using ghodss/yaml and not the pure go-yaml/yaml.
In short, this library first converts YAML to JSON using go-yaml and then uses json.Marshal and json.Unmarshal to convert to or from the struct. This means that it effectively reuses the JSON struct tags as well as the custom JSON methods MarshalJSON and UnmarshalJSON unlike go-yaml

### Implementation

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


// Rührwerk
struct Agitator type
    Name string
    Control: Control Element

struct Heater type
    Name string
    Agitator
    Logic *Logic //Hysteresis (Ursache/Wirking -> Heizen), Overshoot (Überschwingen, Über einen soll zustand hinaus und dann auf diesen einstellt)
    Sensor: Thermometer
    Control: Control Element

// Watchmen
func CheckWaterFlow()
func TempWatcher(Thermometer, Control Element, Output)
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

- 433Utils
- https://github.com/martinohmann/rfoutlet

- [wiringpi](https://tutorials-raspberrypi.de/wiringpi-installieren-pinbelegung/)
- https://github.com/stianeikeland/go-rpio

- https://gobot.io/documentation/platforms/raspi/
- https://godoc.org/periph.io/x/periph/conn/onewire
- https://github.com/kidoman/embd

## Gobot

[doc](https://gobot.io/documentation)

### Installation on Raspberry Pi

Update to latest Raspian Jessie OS and install git and go.
You would normally install Go and Gobot on your workstation.
Once installed, cross compile your program on your workstation, transfer the final executable to your Raspberry Pi.
The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

### Info

- [go-log](https://github.com/go-log/log) logging with interfaces from [this](https://dave.cheney.net/2017/01/23/the-package-level-logger-anti-pattern) article
- [datadog-guide](https://www.datadoghq.com/blog/go-logging/) std logging interface with const values
- [clog](https://github.com/go-clog/clog) use go chan to log back

## Recipe Formats

- [beerxml](http://www.beerxml.com)
- [Beersmith](https://beersmithrecipes.com/)
- [kleiner brauhelfer](https://beersmithrecipes.com/)
- [brewrecipedeveloper](http://www.brewrecipedeveloper.de)


https://github.com/stone/beerxml 5yeas ago -.-

support only one temp at the moment and gpio pin control (ssr)
config values of bus is string because periph need a string input: 40, PWM0, GPIO40 are the same

if you set value "dummy" as device type you are in dummy mode