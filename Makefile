.SILENT: help

COLOR_RESET = \033[0m
COLOR_COMMAND = \033[36m
COLOR_YELLOW = \033[33m
COLOR_GREEN = \033[32m
COLOR_RED = \033[31m

SHELL = /bin/bash
.DEFAULT_GOAL := help

PROJECT := dockerfile-gen

GITHUB_TOKEN := $(shell git config --get github.token || echo $$GITHUB_TOKEN)

TAG := `git describe --tags`
DATE := `date -u +"%Y-%m-%dT%H:%M:%SZ"`
COMMIT := ""

LDFLAGS := -X main.version=$(TAG) -X main.commit=$(COMMIT) -X main.date=$(DATE)

LANGUAGES := $(shell grep language config/languages.yml | awk '{print $$3}')


## Generate versions language yml. Ex: make generate-node
generate-%:
	@printf "\nGenerate $*\n"
	python versions/official-images.py /tmp/all_versions_exported/$*/ $* False save

generate-all:
	@printf "Generate: $(LANGUAGES)\n"
	@for lang in $(LANGUAGES); do make generate-$$lang; done;

packr:
	@packr clean && packr

git-tag:
	@printf "\n"; \
	read -p "Tag ($(TAG)): "; \
	if [ ! "$$REPLY" ]; then \
		printf "\n${COLOR_RED}"; \
		echo "Invalid tag."; \
		exit 1; \
	fi; \
	TAG=$$REPLY; \
	sed -i.bak -r "s/[0-9]+.[0-9]+.[0-9]+/$$TAG/g" README.md && rm README.md.bak 2>/dev/null; \
	sed -i.bak -r "s/[0-9]+.[0-9]+.[0-9]+$$/$$TAG/g" Dockerfile && rm Dockerfile.bak 2>/dev/null; \
	git commit README.md Dockerfile -m "Update README.md and Dockerfile with release $$TAG"; \
	git tag -s $$TAG -m "$$TAG"

## Build project
build: packr
	echo "Building $(PROJECT)"
	go build -ldflags "$(LDFLAGS)" -o $(PROJECT) main.go

## Release of the project
release: packr git-tag
	@if [ ! "$(GITHUB_TOKEN)" ]; then \
		echo "github token should be configurated."; \
		exit 1; \
	fi; \
	export GITHUB_TOKEN=$(GITHUB_TOKEN); \
	goreleaser release --rm-dist; \
	goreleaser release -f .goreleaser-docker.yml --rm-dist; \
	echo "Release - OK"


## Setup of the project
setup:
	@go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	@go get -u github.com/golang/dep/...
	@go get -u github.com/gobuffalo/packr/packr
	@brew install -u goreleaser/tap/goreleaser
	@make vendor-install

## Install dependencies of the project
vendor-install:
	@dep ensure -v

## Update dependencies of the project
vendor-update:
	@dep ensure -update

## Visualizing dependencies status of the project
vendor-status:
	@dep status

## Visualizing dependencies 
vendor-view:
	@brew install graphviz
	@dep status -dot | dot -T png | open -f -a /Applications/Preview.app

lint: ## Run all the linters
	golangci-lint run --skip-dirs official-images

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