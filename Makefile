.PHONY: help run-cli build build-cli build-api release release-linux-amd64 release-linux-arm64 release-darwin-amd64 release-darwin-arm64 release-windows-amd64 test test-clean test-parser test-converter test-integration fmt vet tidy verify clean convert-json convert-csv convert-xlsx convert-ofx convert-all

GO := go
BIN_DIR := bin
DIST_DIR := dist
PDF ?= ./tmp/transaction.pdf
OUT_DIR ?= ./tmp
FORMAT ?= json
OUTPUT ?= $(OUT_DIR)/statement.$(FORMAT)
LDFLAGS := -s -w

help:
	@echo "Targets:"
	@echo "  run-cli            Convert PDF using CLI (PDF, FORMAT, OUTPUT vars)"
	@echo "  build              Build both binaries"
	@echo "  release            Build CLI and API binaries for Linux/macOS/Windows"
	@echo "  release-<platform> Build release binaries for one platform"
	@echo "  test               Run all tests"
	@echo "  test-clean         Clear test cache, then test"
	@echo "  test-parser        Run parser tests only"
	@echo "  test-converter     Run converter tests only"
	@echo "  test-integration   Run parser integration test (ENPARA_TEST_PDF_PATH set from PDF var)"
	@echo "  fmt                Run go fmt"
	@echo "  vet                Run go vet"
	@echo "  tidy               Run go mod tidy"
	@echo "  verify             fmt + vet + test-clean"
	@echo "  convert-all        Generate json/csv/xlsx/ofx from sample PDF"
	@echo "  clean              Remove built binaries and generated outputs"

run-cli:
	$(GO) run cmd/cli/main.go "$(PDF)" -f "$(FORMAT)" -o "$(OUTPUT)"

build: build-cli build-api

build-cli:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/enpara-cli cmd/cli/main.go

build-api:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/enpara-api cmd/api/main.go

release: release-linux-amd64 release-linux-arm64 release-darwin-amd64 release-darwin-arm64 release-windows-amd64

release-linux-amd64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-cli-linux-amd64 ./cmd/cli
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-api-linux-amd64 ./cmd/api

release-linux-arm64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-cli-linux-arm64 ./cmd/cli
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-api-linux-arm64 ./cmd/api

release-darwin-amd64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-cli-darwin-amd64 ./cmd/cli
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-api-darwin-amd64 ./cmd/api

release-darwin-arm64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-cli-darwin-arm64 ./cmd/cli
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-api-darwin-arm64 ./cmd/api

release-windows-amd64:
	mkdir -p $(DIST_DIR)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-cli-windows-amd64.exe ./cmd/cli
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST_DIR)/enpara-api-windows-amd64.exe ./cmd/api

test:
	$(GO) test ./...

test-clean:
	$(GO) clean -testcache
	$(GO) test ./...

test-parser:
	$(GO) test ./tests/parser

test-converter:
	$(GO) test ./tests/converter

test-integration:
	ENPARA_TEST_PDF_PATH="$(abspath $(PDF))" $(GO) test ./tests/parser -run TestParseRealPDFIntegration -v

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

tidy:
	$(GO) mod tidy

verify: fmt vet test-clean

convert-json:
	$(GO) run cmd/cli/main.go "$(PDF)" -f json -o "$(OUT_DIR)/statement.json"

convert-csv:
	$(GO) run cmd/cli/main.go "$(PDF)" -f csv -o "$(OUT_DIR)/statement.csv"

convert-xlsx:
	$(GO) run cmd/cli/main.go "$(PDF)" -f xlsx -o "$(OUT_DIR)/statement.xlsx"

convert-ofx:
	$(GO) run cmd/cli/main.go "$(PDF)" -f ofx -o "$(OUT_DIR)/statement.ofx"

convert-all: convert-json convert-csv convert-xlsx convert-ofx

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	rm -f $(OUT_DIR)/statement.json $(OUT_DIR)/statement.csv $(OUT_DIR)/statement.xlsx $(OUT_DIR)/statement.ofx
