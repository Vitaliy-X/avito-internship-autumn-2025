APP_NAME = pr-service
CMD_DIR = ./cmd/$(APP_NAME)
BIN_DIR = ./bin
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

.PHONY: build run tidy fmt clean

build:
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_FILE) $(CMD_DIR)

run:
	go run $(CMD_DIR)

tidy:
	go mod tidy

fmt:
	go fmt ./...

clean:
	rm -rf $(BIN_DIR)
