# Brewman

## Features

- Log: logger to console, file, socket: Yes with MultipleWriter.

## Not Supported

- At the moment no comments supported in recipe or config file. If you set recipe the config will be overwritten. Parser dont support preserving comments
- Recipe Parser not support Dekoktion and has only kg as unit for malt
- Recipe Parser not support Value Ranges in Recipe (Gaertemperatur: 20-24) and Unit "ml/l".

## Roadmap

- Interface: Sensors (sensors.temp, sensors.flow, sensors.*), Control: 433GHz, Relais
- Metric Exporter: Prometheus, NodeExporter and internal
- Tests

0.1

- [x] read config
- [x] save config
- [x] convert, save and read recipe
- [x] build cmd tool for convert m3 recipes
- [ ] TempWatcher, Temp.Get(), targetTemp, Hold Time

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
