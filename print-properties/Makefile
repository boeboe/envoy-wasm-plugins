# Copyright (c) Tetrate, Inc 2022 All Rights Reserved.

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# https://www.gnu.org/software/make/manual/make.html#One-Shell
# If the .ONESHELL special target appears anywhere in the makefile then all recipe lines for each target will be provided to a single invocation of the shell. Newlines between recipe lines will be preserved.
.ONESHELL:
.PHONY: docker-build docker-push deploy dump-config enable-wasm-debug curl
.SILENT: docker-build docker-push deploy dump-config enable-wasm-debug curl

NAME := envoy-wasm-plugins
WASM := $(NAME).wasm

OUT  := $(WASM)
TAG  ?= print-properties-0.1
HUB  ?= docker.io/boeboe
REPO ?= boeboe/envoy-wasm-plugins

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

all: clean test lint release-docker release-github deploy ## Rebuild, rerelease and redeploy

compile: ## Compile the wasm binary
	@echo "Compiling..."
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test: ## Run golang tests
	@echo "Running tests..."
	go test -v -tags=proxytest ./...

lint: # Lint golang code
	@echo "Linting code..."
	go run $(LINTER) run --verbose --build-tags proxytest

clean:## Remove output artifats
	@echo "Cleaning up..."
	rm -f *.wasm

release-docker: compile ## Release docker container on dockerhub
	@echo "Releasing Docker image..."
	DOCKER_BUILDKIT=0 docker build --platform linux/amd64 --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .
	DOCKER_BUILDKIT=0 docker push $(HUB)/$(NAME):$(TAG)

release-github: compile ## Release wasm binary on github
	@echo "Releasing to GitHub..."
	@echo "Checking if release $(TAG) already exists..."
	@if gh release view $(TAG) -R $(REPO) > /dev/null 2>&1; then \
		echo "Release $(TAG) exists. Deleting it..."; \
		gh release delete $(TAG) -R $(REPO) --yes; \
	fi
	@echo "Creating a new release on GitHub..."
	gh release create $(TAG) $(WASM) --title "Release $(TAG)" --notes "Release notes for $(TAG)" --repo $(REPO)

deploy: ## Deploy the demo applications and wasm filters
	@echo "Deploying..."
	kubectl apply -f istio/wasmplugin.yaml
	kubectl apply -f istio/httpbin-org.yaml
	kubectl apply -f istio/tetrate-io.yaml
	@echo "Force restarting pods..."
	kubectl delete pod $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress

reboot-pods: ## Force reboot workload pods
	@echo "Force restarting pods..."
	kubectl delete pod $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress

dump-config: ## Dump ingress and sidecar envoy configs
	@echo "Dumping config..."
	kubectl exec $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress -- pilot-agent request GET config_dump > ingress_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json

dump-logs: ## Dump ingress and sidecar logs
	@echo "Dumping logs..."
	kubectl -n istio-ingress logs $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) > ingress_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log

enable-full-debug: ## Enable debug:all on ingress and sidecar
	@echo "Enabling full debug..."
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level debug

enable-full-info: ## Enable info:all on ingress and sidecar
	@echo "Enabling full info..."
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level info

enable-http-debug: ## Enable debug:http on ingress and sidecar
	@echo "Enabling HTTP debug..."
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level http:debug

enable-wasm-debug: ## Enable debug:wasm on ingress and sidecar
	@echo "Enabling WASM debug..."
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level wasm:debug

curl: ## Print some sample curl commands for manual verification
	@echo "Call to ingress that forwards to http.org on the internet"
	@echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.httpbin.org:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://wasm.httpbin.org/headers\"
	@echo "Call to ingress that forwards to tetrate.io on the internet"
	@echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://wasm.tetrate.io/\"
	