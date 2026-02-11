VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)

.PHONY: build dev clean lint format format-check ci

build:
	npm run build
	go build -ldflags "$(LDFLAGS)" -o leno .

dev:
	npm run dev

lint:
	npm run lint
	go vet ./...

format:
	npm run format

format-check:
	npm run format:check

ci: format-check lint build

clean:
	rm -rf public/build
	rm -f leno
