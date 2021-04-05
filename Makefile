EXECUTABLE ?= clickbaiter-bot
IMAGE ?= quay.io/cbrgm/$(EXECUTABLE)
GO := CGO_ENABLED=0 go
DATE := $(shell date -u '+%FT%T%z')

LDFLAGS += -X main.Version=$(DRONE_TAG)
LDFLAGS += -X main.Revision=$(DRONE_COMMIT)
LDFLAGS += -X "main.BuildDate=$(DATE)"
LDFLAGS += -extldflags '-static'

PACKAGES = $(shell go list ./...)

.PHONY: all
all: build

.PHONY: clean
clean:
	$(GO) clean -i ./...
	rm -rf dist/

.PHONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

$(EXECUTABLE): $(wildcard *.go)
	$(GO) build -v -ldflags '-w $(LDFLAGS)'

.PHONY: build
build: $(EXECUTABLE)

.PHONY: container
container:
	docker build -t $(IMAGE) .