BUILD_PATH ?= ./cmd/multichat_bot
BIN_PATH ?= ./bin

.PHONY: run
run:
	go run $(BUILD_PATH) -config="local.json"

.PHONY: build
build:
	go build -o=$(BIN_PATH) $(BUILD_PATH)

.PHONY: .install_golangci_lint
.install_golangci_lin:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.55.2

.PHONY: lint
lint: .install_golangci_lin
	golangci-lint run

.PHONY: start-tailwind
start-tailwind:
	npx tailwindcss -c ./website/tailwind.config.js -i ./website/src/css/main.css -o ./website/src/css/main.min.css --watch --minify

