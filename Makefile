BIN=socks5map

ifeq ($(shell uname -s),Darwin)
INSTALL_DIR := /usr/local/bin
else
INSTALL_DIR := /usr/bin
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Displays this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "make \033[36m%-15s\033[0m %s\n", $$1, $$2}'

build: main.go clean ## Builds project
	@go build -o $(BIN) main.go && strip $(BIN)

.PHONY: run
run: ## Runs project
	@go run main.go

.PHONY: clean
clean: ## Cleans project
	@[ ! -f $(BIN) ] || rm -v $(BIN)

.PHONY: install
install: build ## Installs project
	@install -v -b -m 755 $(BIN) $(INSTALL_DIR)

.PHONY: uninstall $(INSTALL_DIR)/$(BIN)
uninstall: ## Uninstall project
	@rm -v $(INSTALL_DIR)/$(BIN)
