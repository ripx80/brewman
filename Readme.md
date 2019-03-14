# Brewman

- cmd/cli Kingpin vs Cobra
- Interface: Sensors (sensors.temp, sensors.flow, sensors.*)
- Control: 433GHz

- Config: Kingpin Parser like Prometheus
- Recepies: yaml files - converted from mmum (from python to golang)

- Metric Exporter: Prometheus, NodeExporter and internal
- Log: logger to console, file, socket?

- Tests

## Roadmap

0.1

- read config
- read recipe
- temp.Watch() event, error

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
