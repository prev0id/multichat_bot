BUILD_PATH ?= ./cmd/multichat_bot
BIN_PATH ?= ./bin

.PHONY: run
run: .install_air
	air --build.cmd "go build -o=$(BIN_PATH) ./cmd/multichat_bot" \
 		--build.bin "./$(BIN_PATH)/multichat_bot" \
 		--build.args_bin "-config='./configs/local.json'" \
 		--build.exclude_dir "website,bin,db" \
 		--build.include_file "./configs/local.json"

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

.PHONY: .install_air
.install_air:
	go install github.com/cosmtrek/air@v1.52.0
