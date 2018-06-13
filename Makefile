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

dist-clean:
	rm -rf dist
	rm -f $(PROJECT)-*.tar.gz


dist: dist-clean
	mkdir -p dist/alpine-linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/amd64/$(PROJECT) main.go
	mkdir -p dist/alpine-linux/arm64 && GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/arm64/$(PROJECT) main.go
	mkdir -p dist/alpine-linux/armhf && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -a -tags netgo -installsuffix netgo -o dist/alpine-linux/armhf/$(PROJECT) main.go
	mkdir -p dist/linux/amd64 && GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/linux/amd64/$(PROJECT) main.go
	mkdir -p dist/linux/arm64 && GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/linux/arm64/$(PROJECT) main.go
	mkdir -p dist/linux/i386  && GOOS=linux GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/linux/i386/$(PROJECT) main.go
	mkdir -p dist/linux/armel  && GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "$(LDFLAGS)" -o dist/linux/armel/$(PROJECT) main.go
	mkdir -p dist/linux/armhf  && GOOS=linux GOARCH=arm GOARM=6 go build -ldflags "$(LDFLAGS)" -o dist/linux/armhf/$(PROJECT) main.go
	mkdir -p dist/darwin/amd64 && GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/darwin/amd64/$(PROJECT) main.go
	mkdir -p dist/darwin/i386  && GOOS=darwin GOARCH=386 go build -ldflags "$(LDFLAGS)" -o dist/darwin/i386/$(PROJECT) main.go


release: dist
	tar -cvzf $(PROJECT)-alpine-linux-amd64-$(TAG).tar.gz -C dist/alpine-linux/amd64 $(PROJECT)
	tar -cvzf $(PROJECT)-alpine-linux-arm64-$(TAG).tar.gz -C dist/alpine-linux/arm64 $(PROJECT)
	tar -cvzf $(PROJECT)-alpine-linux-armhf-$(TAG).tar.gz -C dist/alpine-linux/armhf $(PROJECT)
	tar -cvzf $(PROJECT)-linux-amd64-$(TAG).tar.gz -C dist/linux/amd64 $(PROJECT)
	tar -cvzf $(PROJECT)-linux-arm64-$(TAG).tar.gz -C dist/linux/arm64 $(PROJECT)
	tar -cvzf $(PROJECT)-linux-i386-$(TAG).tar.gz -C dist/linux/i386 $(PROJECT)
	tar -cvzf $(PROJECT)-linux-armel-$(TAG).tar.gz -C dist/linux/armel $(PROJECT)
	tar -cvzf $(PROJECT)-linux-armhf-$(TAG).tar.gz -C dist/linux/armhf $(PROJECT)
	tar -cvzf $(PROJECT)-darwin-amd64-$(TAG).tar.gz -C dist/darwin/amd64 $(PROJECT)
	tar -cvzf $(PROJECT)-darwin-i386-$(TAG).tar.gz -C dist/darwin/i386 $(PROJECT)


## Setup of the project
setup:
	@go get -u github.com/alecthomas/gometalinter
	@go get -u github.com/golang/dep/...
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