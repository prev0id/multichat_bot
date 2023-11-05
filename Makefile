BUILD_PATH ?= ./cmd/multichat_bot
BIN_PATH ?= ./bin

PHONY: run
run:
	go run $(BUILD_PATH)

.PHONY: build
build:
	go build -o=$(BIN_PATH) $(BUILD_PATH)

