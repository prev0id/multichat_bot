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

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	oapi-codegen -config ".oapi-codegen.yaml" ./api/specification.yml