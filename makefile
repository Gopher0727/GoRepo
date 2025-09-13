# Makefile for GoRepo

GO := go
MOD := github.com/Gopher0727/GoRepo
BIN_DIR := bin
BACKEND_BIN := $(BIN_DIR)/backend
JSONR_BIN := $(BIN_DIR)/jsonrepair
BACKEND_PKG := ./backend/cmd
JSONR_PKG := ./tools/JsonRepair

.PHONY: all build run run-dev gen-compose run-jsonrepair test test-mysql test-redis test-logger fmt vet tidy mod-download clean

all: build

build: $(BACKEND_BIN) $(JSONR_BIN)

$(BACKEND_BIN):
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(BACKEND_BIN) $(BACKEND_PKG)

$(JSONR_BIN):
	@mkdir -p $(BIN_DIR)
	$(GO) build -o $(JSONR_BIN) $(JSONR_PKG)

run: build
	@echo "运行后端：./$(BACKEND_BIN)"
	./$(BACKEND_BIN)

run-dev:
	@echo "使用 go run 直接运行后端（用于开发）"
	$(GO) run $(BACKEND_PKG)

gen-compose:
	@echo "通过运行 [backend/cmd/main.go] 生成 docker-compose.yml"
	$(GO) run $(BACKEND_PKG)

run-jsonrepair:
	@echo "运行 tools/JsonRepair 示例程序"
	@$(GO) run $(JSONR_PKG)

test:
	@$(GO) test ./...

test-mysql:
	@echo "注意：需要本机运行 MySQL 并确保 [config.toml] 中配置正确"
	@$(GO) test ./tests -run TestMySQL -v

test-redis:
	@echo "注意：需要本机运行 Redis 并确保 [config.toml] 中配置正确"
	@$(GO) test ./tests -run TestRedis -v

test-logger:
	@$(GO) test ./tests -run TestLogger -v

fmt:
	@goimports-reviser.exe -rm-unused -format -recursive .
	@# $(GO) fmt ./... > /dev/null 2>&1

vet:
	@$(GO) vet ./...

tidy:
	@$(GO) mod tidy

mod-download:
	@$(GO) mod download

clean:
	@echo "清理 bin/ 和 docker-compose.yml"
	rm -rf $(BIN_DIR) docker-compose.yml