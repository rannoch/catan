GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=catan
LINTER=golangci-lint

all: test lint

test:
	$(GOTEST) ./... -v -coverprofile=coverage.txt -covermode=atomic

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

lint:
	$(LINTER) run