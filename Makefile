SOURCES  = $(shell find . -name '*.go')

.PHONY: default
default: lint

.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: lib
lib: $(SOURCES) ## build  library
	go build ./...

.PHONY: deps
deps: ## install dependencies to run everything
	go env
	@go install honnef.co/go/tools/cmd/staticcheck@latest

.PHONY: lint
lint: vet staticcheck ## run all linters

.PHONY: vet
vet: $(SOURCES) ## run Go vet
	go vet ./...

.PHONY: staticcheck
# -ST1000 missing package doc in many packages
# -ST1003 wrong naming convention Api vs API, Id vs ID
# -ST1020 too many wrong comments on exported functions to fix right away
staticcheck: $(SOURCES) ## run staticcheck
	staticcheck -checks "all" ./...

.PHONY: check-fmt
check-fmt: $(SOURCES) ## check format code
	@if [ "$$(gofmt -s -d $(SOURCES))" != "" ]; then false; else true; fi

.PHONY: check-race
check-race: lib ## run all tests with race checker
	go test -race ./...
