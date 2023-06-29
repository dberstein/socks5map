BIN=socks5map

.PHONY: help
help: ## Displays this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: run
run: ## Runs project
	@go run main.go

.PHONY: clean
clean: ## Cleans project
	@[ ! -f $(BIN) ] || rm -v $(BIN)

.PHONY: build
build: clean ## Builds project
	@go build -o $(BIN) main.go

.PHONY: install
install: build ## Installs project
	@echo TODO
