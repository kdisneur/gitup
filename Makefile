GO := go
GOFMT := gofmt
GOLINT := golint
GOCOMPILE_OPTIONS :=
JQ := jq
GOTEST_OPTIONS :=
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git branch --no-color | awk '/^\* / { print $$2 }')

OS := darwin
ARCH := 386
BINARY_NAME := "gitup"
FULL_BINARY_NAME := $(BINARY_NAME)-$(OS)-$(ARCH)

GIT_STATE := $(shell if [ -z "$(shell git status --short)" ]; then echo clean; else echo dirty; fi)

VERSION := v0.1.0
ALREADY_RELEASED := $(shell if [ $$(curl --silent --output /dev/null --write-out "%{http_code}" https://api.github.com/repos/kdisneur/gitup/releases/tags/$(VERSION)) -eq 200 ]; then echo "true"; else echo "false"; fi)
IS_RELEASE := $(shell if git log --format=%B -n1 HEAD | grep -q '\[release\]'; then echo "true"; fi)
PRE_RELEASE_VERSION := $(shell git log --oneline  HEAD...$$(git log --oneline --format='%H' --grep '\[release\]' -n 1) | wc -l | awk '{print $$1}')
BUILD_NUMBER=

clean:
	@$(GO) clean
	@rm -rf build

compile: _dependencies
	@if [ "$(ALREADY_RELEASED)" = "true" ]; then \
		echo "$(VERSION) already released."; \
		exit 1; \
	fi

	@touch pkg/version/version.go
	@GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GOCOMPILE_OPTIONS) \
		-ldflags \
			"-X github.com/kdisneur/gitup/pkg/version.gitBranch=$(GIT_BRANCH) \
			 -X github.com/kdisneur/gitup/pkg/version.gitCommit=$(GIT_COMMIT) \
			 -X github.com/kdisneur/gitup/pkg/version.gitState=$(GIT_STATE) \
			 -X github.com/kdisneur/gitup/pkg/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ') \
			 -X github.com/kdisneur/gitup/pkg/version.version=$(VERSION) \
			 -X github.com/kdisneur/gitup/pkg/version.isRelease=$(IS_RELEASE) \
			 -X github.com/kdisneur/gitup/pkg/version.prereleaseNumber=$(PRE_RELEASE_VERSION) \
			 -X github.com/kdisneur/gitup/pkg/version.buildNumber=$(BUILD_NUMBER)" \
		-o build/$(BINARY_NAME)

	@tar czf build/$(FULL_BINARY_NAME).tgz -C build  $(BINARY_NAME)
	@echo "archive generated at build/$(FULL_BINARY_NAME).tgz"

	@mv build/$(BINARY_NAME) build/$(FULL_BINARY_NAME)
	@echo "archive generated at build/$(FULL_BINARY_NAME)"

test-style: _gofmt _golint

test-unit:
	@$(GO) test $(GOTEST_OPTIONS) ./...

_dependencies:
	@$(GO) mod download

_gofmt:
	@data=$$($(GOFMT) -l main.go pkg);\
	if [ -n "$$data" ]; then \
		echo $$data; \
		exit 1; \
	fi

_golint:
	@$(GOLINT) -set_exit_status pkg/... .
