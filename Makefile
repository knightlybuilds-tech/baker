VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/knightlybuilds-tech/baker/internal/cli.Version=$(VERSION)

.PHONY: build install test lint clean bakery-ui

# Bakery UI source (static). Synced into internal/bakery/ui, which is committed
# so `go install` works without the sibling checkout.
BAKERY_UI ?= ../bakery-workspace/dist

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

bakery-ui:
	@test -f "$(BAKERY_UI)/index.html" || (echo "Bakery UI not found at $(BAKERY_UI)"; exit 1)
	rm -rf internal/bakery/ui && mkdir -p internal/bakery/ui
	cp -R "$(BAKERY_UI)"/. internal/bakery/ui/
	@echo "✓ synced Bakery UI from $(BAKERY_UI)"
