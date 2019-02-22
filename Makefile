GO := go
GOFMT := gofmt
GOLINT := golint
GOCOMPILE_OPTIONS :=
GOTEST_OPTIONS :=
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git branch --no-color | awk '/^\* / { print $$2 }')
BINARY_NAME := "gitup"
GIT_STATE := $(shell if [ -z "$(shell git status --short)" ]; then echo clean; else echo dirty; fi)

clean:
	@$(GO) clean
	@rm -rf build

compile: _dependencies
	@touch pkg/version/version.go
	@$(GO) build $(GOCOMPILE_OPTIONS) \
		-ldflags \
			"-X github.com/kdisneur/gitup/pkg/version.gitBranch=$(GIT_BRANCH) \
			 -X github.com/kdisneur/gitup/pkg/version.gitCommit=$(GIT_COMMIT) \
			 -X github.com/kdisneur/gitup/pkg/version.gitState=$(GIT_STATE) \
			 -X github.com/kdisneur/gitup/pkg/version.buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		-o build/$(BINARY_NAME)
	@echo "file generated at build/$(BINARY_NAME)"

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