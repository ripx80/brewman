#!/bin/sh
GOOS=linux GOARCH=arm go build -o brewman-arm cmd/brewman/main.go && scp brewman-arm pi@pi:~/go/src/github.com/ripx80/brewman/cmd/brewman/
exit 0
