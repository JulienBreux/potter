BINDIR  := $(CURDIR)/bin
BINNAME ?= potter

test: ## [project] Test packages
	go test all -count=1 -cover -coverprofile=coverage.out -v ./...

upgrade-deps: ## [project] Upgrade all dependencies
	go get -u ./...

build: ## [binary] Build local binary
	mkdir -p ./bin
	go build -o ./bin ./cmd/${BINNAME}

run: build ## [binary] Run local binary
	./bin/${BINNAME}

image-build: ## [image] Build local image
	docker build -t ghcr.io/julienbreux/${BINNAME}:latest .

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: image-build
