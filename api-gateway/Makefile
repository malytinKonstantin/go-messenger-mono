# Переменные
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet
GOLINT=golangci-lint

# Пути
SRC_DIRS := ./...

# Команды
.PHONY: all build clean test coverage deps fmt lint vet

all: build

build:
	$(GOBUILD) $(SRC_DIRS)

clean:
	$(GOCLEAN)

test:
	$(GOTEST) -v $(SRC_DIRS)

coverage:
	$(GOTEST) -v -coverprofile=coverage.out $(SRC_DIRS)
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

deps:
	$(GOGET) -v -t -d $(SRC_DIRS)
	$(GOMOD) tidy

fmt:
	$(GOFMT) -w .

lint:
	$(GOLINT) run

vet:
	$(GOVET) $(SRC_DIRS)

check: fmt lint vet test