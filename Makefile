.PHONY: help run-cli build build-cli build-api pre-release release release-linux-amd64 release-linux-arm64 release-darwin-amd64 release-darwin-arm64 release-windows-amd64 release-archives release-clean-binaries release-checksums release-sign test test-clean test-parser test-converter test-integration fmt vet tidy verify clean convert-json convert-csv convert-xlsx convert-ofx convert-all convert-manual convert-automatic

GO := go
BIN_DIR := bin
DIST_DIR := dist
PDF ?= ./tmp/manual.pdf
OUT_DIR ?= ./tmp
FORMAT ?= json
PDF_TYPE ?= auto
OUTPUT ?= $(OUT_DIR)/statement.$(FORMAT)
LDFLAGS := -s -w

help:
	@echo "Targets:"
	@echo "  run-cli            Convert PDF using CLI (PDF, FORMAT, OUTPUT, PDF_TYPE vars)"
	@echo "  build              Build both binaries"
	@echo "  release            Run verify, build, archive, checksum, and signing steps"
	@echo "  test               Run all tests"
	@echo "  test-clean         Clear test cache, then test"
	@echo "  test-parser        Run parser tests only"
	@echo "  test-converter     Run converter tests only"
	@echo "  test-integration   Run parser integration test (ENPARA_TEST_PDF_PATH set from PDF var)"
	@echo "  fmt                Run go fmt"
	@echo "  vet                Run go vet"
	@echo "  tidy               Run go mod tidy"
	@echo "  verify             fmt + vet + test-clean"
	@echo "  convert-all        Generate json/csv/xlsx/ofx from PDF"
	@echo "  convert-manual     Generate all outputs from ./tmp/manual.pdf (type1)"
	@echo "  convert-automatic  Generate all outputs from ./tmp/automatic.pdf (type2)"
	@echo "  clean              Remove built binaries and generated outputs"
	@echo ""
	@echo "Variables (override with make VAR=value):"
	@echo "  PDF=$(PDF)"
	@echo "  FORMAT=$(FORMAT)"
	@echo "  PDF_TYPE=$(PDF_TYPE)"
	@echo "  OUT_DIR=$(OUT_DIR)"
	@echo "  OUTPUT=$(OUTPUT)"
	@echo ""
run-cli:
	$(GO) run cmd/cli/main.go "$(PDF)" -f "$(FORMAT)" --type "$(PDF_TYPE)" -o "$(OUTPUT)"

build: build-cli build-api

build-cli:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/enpara-cli cmd/cli/main.go

build-api:
	mkdir -p $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/enpara-api cmd/api/main.go

pre-release: verify
	@echo "Verification completed"

release: pre-release release-linux-amd64 release-linux-arm64 release-darwin-amd64 release-darwin-arm64 release-windows-amd64 release-archives release-clean-binaries release-checksums release-sign
	@echo "Release artifacts generated in $(DIST_DIR)/"
	@echo "Produced: enpara-*.zip, CHECKSUMS.sha256, CHECKSUMS.sha256.asc"

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

release-archives:
	@mkdir -p $(DIST_DIR)
	@cd $(DIST_DIR) && \
	zip -q enpara-linux-amd64.zip enpara-cli-linux-amd64 enpara-api-linux-amd64 && \
	zip -q enpara-linux-arm64.zip enpara-cli-linux-arm64 enpara-api-linux-arm64 && \
	zip -q enpara-darwin-amd64.zip enpara-cli-darwin-amd64 enpara-api-darwin-amd64 && \
	zip -q enpara-darwin-arm64.zip enpara-cli-darwin-arm64 enpara-api-darwin-arm64 && \
	zip -q enpara-windows-amd64.zip enpara-cli-windows-amd64.exe enpara-api-windows-amd64.exe
	@echo "✓ Release archives created in $(DIST_DIR)/"

release-clean-binaries:
	@rm -f \
		$(DIST_DIR)/enpara-cli-linux-amd64 \
		$(DIST_DIR)/enpara-api-linux-amd64 \
		$(DIST_DIR)/enpara-cli-linux-arm64 \
		$(DIST_DIR)/enpara-api-linux-arm64 \
		$(DIST_DIR)/enpara-cli-darwin-amd64 \
		$(DIST_DIR)/enpara-api-darwin-amd64 \
		$(DIST_DIR)/enpara-cli-darwin-arm64 \
		$(DIST_DIR)/enpara-api-darwin-arm64 \
		$(DIST_DIR)/enpara-cli-windows-amd64.exe \
		$(DIST_DIR)/enpara-api-windows-amd64.exe
	@echo "✓ Intermediate release binaries removed from $(DIST_DIR)/"

release-checksums:
	@cd $(DIST_DIR) && \
	sha256sum enpara-*.zip > CHECKSUMS.sha256 && \
	cat CHECKSUMS.sha256
	@echo "✓ Checksums generated: $(DIST_DIR)/CHECKSUMS.sha256"

release-sign:
	@cd $(DIST_DIR) && \
	gpg --batch --yes --detach-sign --armor CHECKSUMS.sha256
	@echo "✓ Release signed: $(DIST_DIR)/CHECKSUMS.sha256.asc"
	@echo "Users can verify with:"
	@echo "  gpg --verify CHECKSUMS.sha256.asc"
	@echo "  sha256sum -c CHECKSUMS.sha256"

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
	$(GO) run cmd/cli/main.go "$(PDF)" -f json --type "$(PDF_TYPE)" -o "$(OUT_DIR)/statement.json"

convert-csv:
	$(GO) run cmd/cli/main.go "$(PDF)" -f csv --type "$(PDF_TYPE)" -o "$(OUT_DIR)/statement.csv"

convert-xlsx:
	$(GO) run cmd/cli/main.go "$(PDF)" -f xlsx --type "$(PDF_TYPE)" -o "$(OUT_DIR)/statement.xlsx"

convert-ofx:
	$(GO) run cmd/cli/main.go "$(PDF)" -f ofx --type "$(PDF_TYPE)" -o "$(OUT_DIR)/statement.ofx"

convert-all: convert-json convert-csv convert-xlsx convert-ofx

convert-manual:
	$(MAKE) convert-all PDF=./tmp/manual.pdf PDF_TYPE=type1

convert-automatic:
	$(MAKE) convert-all PDF=./tmp/automatic.pdf PDF_TYPE=type2

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(DIST_DIR)
	rm -f $(OUT_DIR)/statement.json $(OUT_DIR)/statement.csv $(OUT_DIR)/statement.xlsx $(OUT_DIR)/statement.ofx
