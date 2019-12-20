# GOARM=6 (Raspberry Pi A, A+, B, B+, Zero) GOARM=7 (Raspberry Pi 2, 3)
# change to GOBIN if set
BIN = $(CURDIR)/bin
 # Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOARMOPTS="GOARM=6 GOARCH=arm GOOS=linux"


all: test build
build:
		# unix
		$(GOBUILD) cmd/brewman/main.go -o $($BIN)/brewman -v
		# arm
		$(GOARMOPTS) $(GOBUILD) cmd/brewman/main.go -o $(BIN)/brewman-arm -v
test:
		$(GOTEST) -v ./...
clean:
		$(GOCLEAN)
		rm -f $(BIN)/bin/*

deps:
		$(GOGET) github.com/markbates/goth
		$(GOGET) github.com/markbates/pop


# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v