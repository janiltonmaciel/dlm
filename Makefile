.SILENT:

SHELL = /bin/bash
.DEFAULT_GOAL := help

PROJECT := dockerfile-gen
TAG := `git describe --tags`
DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
COMMIT := ""

LDFLAGS := -X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(DATE)


build:
	echo "Building $(PROJECT)"
	go build -ldflags "$(LDFLAGS)" -o $(PROJECT) main.go

release:
	goreleaser --rm-dist

## Setup of the project
setup:
	@go get -u github.com/alecthomas/gometalinter
	@go get -u github.com/golang/dep/...
	@brew install goreleaser/tap/goreleaser
	@make vendor-install
	gometalinter --install --update

## Install dependencies of the project
vendor-install:
	@dep ensure -v

## Visualizing dependencies status of the project
vendor-status:
	@dep status

lint: ## Run all the linters
	gometalinter --vendor --disable-all \
		--enable=deadcode \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=staticcheck \
		--enable=gofmt \
		--enable=goimports \
		--enable=dupl \
		--enable=misspell \
		--enable=errcheck \
		--enable=vet \
		--enable=vetshadow \
		--deadline=10m \
		--aggregate \
		./...


COLOR_RESET = \033[0m
COLOR_COMMAND = \033[36m
COLOR_YELLOW = \033[33m

## Prints this help
help:
	printf "${COLOR_YELLOW}dockerfile-gen\n------\n${COLOR_RESET}"
	awk '/^[a-zA-Z\-\_0-9\.%]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "${COLOR_COMMAND}$$ make %s${COLOR_RESET} %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort
	printf "\n"