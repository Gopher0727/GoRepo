# Makefile

GO := go
BIN_DIR := bin
BACKEND_BIN := $(BIN_DIR)/backend
BACKEND_PKG := ./backend/cmd

.PHONY: build run test clean

build:
	@mkdir -p $(BIN_DIR)
	@$(GO) build -o $(BACKEND_BIN) $(BACKEND_PKG)
	@./$(BACKEND_BIN)

run:
	@$(GO) run $(BACKEND_PKG)

test:
	@$(GO) test ./...

clean:
	@$(GO) mod tidy
	@goimports-reviser.exe -format -recursive .
	@rm -rf $(BIN_DIR)
