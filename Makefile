# Self documented Makefile
# http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help: ## Show list of make targets and their description
	@grep -E '^[%.a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help

PROJECT_ROOT=$(shell pwd)

.PHONY: all
all:
	make clean
	make setup
	make build.build

.PHONY: setup
setup: ## Run setup scripts to prepare development environment
	@scripts/setup.sh

.PHONY: build
build: ## Show build.sh help
	@scripts/build.sh

.PHONY: build.
build.%: ## Build artifact defined by '%'. '%' can be client, server, or testserver. Ex: 'make build.server` will trigger ./scripts/build.sh server
	@scripts/build.sh $*

.PHONY: lint
lint: ## Run linter with --fast option, for most common issues
	@scripts/lint.sh

.PHONY: clean
clean:
	rm -rf ./bin