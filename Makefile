BUILD_PATH ?= ./cmd/multichat_bot
WEBSITE_BUILD_PATH ?= ./cmd/serve_website
BIN_PATH ?= ./bin

PHONY: serve
serve:
	go run $(WEBSITE_BUILD_PATH)

PHONY: run
run:
	go run $(BUILD_PATH) -config="./configs/local.json"


.PHONY: build
build:
	go build -o=$(BIN_PATH) $(BUILD_PATH)

