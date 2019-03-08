# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GORUN=$(GOCMD) run
BINARY_NAME=npmdepcopy

all: test build
install:
	$(GOGET) ./...
build:
	mkdir -p builds
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)_win.exe -v ./cmd/npmdepcopy/...
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)_osx -v ./cmd/npmdepcopy/...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o builds/$(BINARY_NAME)_unix -v ./cmd/npmdepcopy/...
test:
	$(GOTEST) -v ./src/...
clean:
	$(GOCLEAN) ./src/...
	rm -f scripts
run:
	$(GORUN) ./cmd/npmdepcopy/main.go