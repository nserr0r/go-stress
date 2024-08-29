BINARY_NAME := go-stress
SRC := $(wildcard *.go)
INSTALL_DIR := /usr/local/bin

LDFLAGS :=

YELLOW := $(shell tput setaf 3)
GREEN := $(shell tput setaf 2)
RESET := $(shell tput sgr0)

.PHONY: all
all: build

.PHONY: build
build:
	@echo "$(YELLOW)Building the project...$(RESET)"
	@go mod init main
	@go tidy
	@go build -ldflags="$(LDFLAGS)" -o $(BINARY_NAME) $(SRC)
	@echo "$(GREEN)Build complete!$(RESET)"

.PHONY: clean
clean:
	@echo "$(YELLOW)Cleaning up...$(RESET)"
	@rm -f $(BINARY_NAME)
	@echo "$(GREEN)Cleanup complete!$(RESET)"

.PHONY: run
run: build
	@echo "$(YELLOW)Running the project...$(RESET)"
	@./$(BINARY_NAME) -host=localhost:3001 -path=/crypt/ws -conn=10 -msg=5

.PHONY: install
install: build
	@echo "$(YELLOW)Installing to $(INSTALL_DIR)...$(RESET)"
	@install -m 0755 $(BINARY_NAME) $(INSTALL_DIR)
	@echo "$(GREEN)Installation complete!$(RESET)"

.PHONY: test
test:
	@echo "$(YELLOW)Running tests...$(RESET)"
	@go test ./...
	@echo "$(GREEN)Tests complete!$(RESET)" 

.PHONY: coverage
coverage:
	@echo "$(YELLOW)Generating coverage report...$(RESET)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out
	@echo "$(GREEN)Coverage report generated!$(RESET)"

.PHONY: deps-clean
deps-clean:
	@echo "$(YELLOW)Cleaning Go module cache...$(RESET)"
	@go clean -modcache
	@echo "$(GREEN)Go module cache cleaned!$(RESET)"

.PHONY: deps-update
deps-update:
	@echo "$(YELLOW)Updating project dependencies...$(RESET)"
	@go get -u ./...
	@echo "$(GREEN)Dependencies updated!$(RESET)"
