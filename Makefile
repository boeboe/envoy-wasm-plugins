# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# https://www.gnu.org/software/make/manual/make.html#One-Shell
# If the .ONESHELL special target appears anywhere in the makefile then all recipe lines for each target will be provided to a single invocation of the shell. Newlines between recipe lines will be preserved.
# .ONESHELL:
.PHONY: kind-up istio-up kind-down
.SILENT: kind-up istio-up kind-down

up: istio-up ## Bring up full demo scenario
clean: kind-down
down: kind-down ## Bring down full demo scenario

kind-up: ## Spin up a local kind cluster
	@/bin/sh -c './make.sh $@'

istio-up: kind-up ## Install istio in the kind cluster
	@/bin/sh -c './make.sh $@'

kind-down: ## Remove the local kind cluster
	@/bin/sh -c './make.sh $@'