VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/knightlybuilds-tech/baker/internal/cli.Version=$(VERSION)

.PHONY: build install test lint clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/baker ./cmd/baker

install:
	go install -ldflags "$(LDFLAGS)" ./cmd/baker
	@bin="$$(go env GOBIN)"; [ -n "$$bin" ] || bin="$$(go env GOPATH)/bin"; \
	echo "✓ baker installed to $$bin/baker"; \
	case ":$$PATH:" in *":$$bin:"*) ;; *) echo; "$$bin/baker" --version >/dev/null;; esac

test:
	go test ./...

lint:
	go vet ./...
	@test -z "$$(gofmt -l .)" || (echo "gofmt needed:"; gofmt -l .; exit 1)

clean:
	rm -rf bin
