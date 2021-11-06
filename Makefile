PROJECTNAME="TetraControl"

# Go related variables.
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/build
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

go-full-build: go-clean go-get go-build

go-build:
	@echo "  >  Building binary..."
	@mkdir -p $(GOBIN)
	@GOOS=linux CGO_ENABLED=0 go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	@chmod 755 $(GOBIN)/$(PROJECTNAME)

go-generate:
	@echo "  >  Generating dependency files..."
	@go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@go get $(get)

go-install:
	@echo "  >  Running go install..."
	@go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@go clean

go-test:
	@echo "  >  Running tests..."
	@go test

go-run:
	@echo "  >  Running ${PROJECTNAME}"
	@-(cd $(GOBIN); ./$(PROJECTNAME))

## install: Download and install dependencies
install: go-get

# clean: Runs go clean
clean: go-clean

## full-build: cleans project, installs dependencies, and builds project
full-build: go-full-build

## build: Runs go build
build: go-build

## package: Builds lambda zip
package: go-full-build
	@echo "  >  Zipping package..."
	@cd $(GOBIN) && zip $(PROJECTNAME).zip $(PROJECTNAME)

## clean: Runs go clean
clean:
	@rm -rf build

## run: full-builds and executes project binary
run: go-build go-run

## test: Run unit tests
test: go-test

## help: Displays help text for make commands
.DEFAULT_GOAL := help
all: help
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'