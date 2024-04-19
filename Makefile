BUILD_PATH ?= ./cmd/multichat_bot
WEBSITE_BUILD_PATH ?= ./cmd/serve_website
BIN_PATH ?= ./bin

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

.PHONY: redis
redis:
	docker run -d --name redis-stack-server -p 6379:6379 redis/redis-stack-server:latest

.PHONY: start-tailwind
start-tailwind:
	npx tailwindcss -i ./website/src/main.css -o ./website/src/compiled.css --watch --minify

.PHONY: .postgres
.postgres:
	docker run
