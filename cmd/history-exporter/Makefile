GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=history-exporter

all: test build

test:
	$(GOTEST) -v -cover ./...
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist
	rm -rf output
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

release-test:
	echo [INFO] Creating non-validated build
	go build .
	mkdir -p output
	mv $(BINARY_NAME) output/$(BINARY_NAME)
	goreleaser release --key $(GORELEASER_KEY) --rm-dist --snapshot

build:
	@echo [INFO] Creating build
	goreleaser build --rm-dist

release:
	@echo [INFO] Updating templates
	goreleaser release --key $(GORELEASER_KEY) --rm-dist
