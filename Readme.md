# Brewman

## Setup

```bash
make; ./bin/brewman
```

You will find your executables in bin

- brewman - Brewery Helper

build with upx

./build.sh

## Features

describe comming soon

## Not Supported

- At the moment no comments supported in recipe or config file. If you set recipe the config will be overwritten. Parser dont support preserving comments
- Recipe Parser not support Dekoktion and has only kg as unit for malt
- Recipe Parser not support Value Ranges in Recipe (Gaertemperatur: 20-24) and Unit "ml/l".
- Periph.io onewire, I'm not gonna be able to get it to work. Sorry I used the files in sys

## Units

- all weight sizes are gramm
- all time units are minutes
