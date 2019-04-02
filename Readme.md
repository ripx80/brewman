# Brewman

- cmd/cli Kingpin vs Cobra, Kingpin have simple configuration. Cobra has more complexity and power

- Interface: Sensors (sensors.temp, sensors.flow, sensors.*), Control: 433GHz, Relais

- Config: Kingpin Parser like Prometheus
- Recepies: yaml files - converted from mmum (from python to golang)

- Metric Exporter: Prometheus, NodeExporter and internal
- Log: logger to console, file, socket: Yes with MultipleWriter.

- Tests

- At the moment no comments supported. If you set recipe the config will be overwritten. Parser dont support preserving comments

- Recipe Parser not support Dekoktion and has only kg as unit for malt
- Recipe Parser not support Value Ranges in Recipe (Gaertemperatur: 20-24)

Ideas:
    - set the max volume of mesher, cooker, fermenter. So you can auto calc max Outcome with MainCast and Grouting!

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

## Roadmap

0.1

- [x] read config
- [x] save config
- convert, read recipe
- temp.Watch() event, error

1.x

- recipe:Download recipe, search recipe

## Dependencies

- 433Utils
- https://github.com/martinohmann/rfoutlet

- [wiringpi](https://tutorials-raspberrypi.de/wiringpi-installieren-pinbelegung/)
- https://github.com/stianeikeland/go-rpio

- https://gobot.io/documentation/platforms/raspi/
- https://godoc.org/periph.io/x/periph/conn/onewire

## Gobot

[doc](https://gobot.io/documentation)

### Installation on Raspberry Pi

Update to latest Raspian Jessie OS and install git and go.
You would normally install Go and Gobot on your workstation.
Once installed, cross compile your program on your workstation, transfer the final executable to your Raspberry Pi.
The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.
