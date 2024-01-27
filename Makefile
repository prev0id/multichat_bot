BUILD_PATH ?= ./cmd/multichat_bot
WEBSITE_BUILD_PATH ?= ./cmd/serve_website
BIN_PATH ?= ./bin
SERVE_PORT ?= 8080

.PHONY: serve
serve:
	go run $(WEBSITE_BUILD_PATH) port=$(SERVE_PORT)

.PHONY: run
run:
	go run $(BUILD_PATH) -config="./configs/local.json"

.PHONY: build
build:
	go build -o=$(BIN_PATH) $(BUILD_PATH)

.PHONY: .install_golangci_lint
.install_golangci_lin:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: lint
lint: .install_golangci_lin
	golangci-lint run

.PHONY: .install_oapi_codegen
.install_oapi_codegen:
	go install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0

.PHONY: generate
generate: .install_oapi_codegen
	oapi-codegen -config ".oapi-codegen.yaml" ./api/specification.yml