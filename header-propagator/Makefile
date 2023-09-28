# Copyright (c) Tetrate, Inc 2022 All Rights Reserved.

# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## This help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

.ONESHELL:
.PHONY: compile release-docker release-github deploy dump-config enable-wasm-debug curl
.SILENT: compile release-docker release-github deploy dump-config enable-wasm-debug curl

NAME := envoy-wasm-plugins
WASM := $(NAME).wasm

OUT  := $(WASM)
TAG  ?= header-propagator-0.1
HUB  ?= docker.io/boeboe
REPO ?= boeboe/envoy-wasm-plugins

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

all: clean test lint release-docker release-github deploy

compile:
	@echo "Compiling..."
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test:
	@echo "Running tests..."
	go test -v -tags=proxytest ./...

lint:
	@echo "Linting code..."
	go run $(LINTER) run --verbose --build-tags proxytest

clean:
	@echo "Cleaning up..."
	rm -f *.wasm

release-docker: compile
	@echo "Releasing Docker image..."
	docker build --platform linux/amd64 --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .
	docker push $(HUB)/$(NAME):$(TAG)

release-github: compile
	@echo "Releasing to GitHub..."
	@echo "Checking if release $(TAG) already exists..."
	@if gh release view $(TAG) -R $(REPO) > /dev/null 2>&1; then \
		echo "Release $(TAG) exists. Deleting it..."; \
		gh release delete $(TAG) -R $(REPO) --yes; \
	fi
	@echo "Creating a new release on GitHub..."
	gh release create $(TAG) $(WASM) --title "Release $(TAG)" --notes "Release notes for $(TAG)" --repo $(REPO)

deploy:
	@echo "Deploying..."
	kubectl apply -f istio/zipkin.yaml
	kubectl apply -f istio/gateway.yaml
	kubectl apply -f istio/app-a.yaml
	kubectl apply -f istio/app-b.yaml
	kubectl apply -f istio/httpbin.yaml
	kubectl apply -f istio/wasmplugin.yaml
	@echo "Force restarting pods..."
	kubectl delete pod $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a) -n ns-a
	kubectl delete pod $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b) -n ns-b
	kubectl delete pod $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin) -n ns-httpbin

dump-config:
	@echo "Dumping config..."
	kubectl exec $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress -- pilot-agent request GET config_dump > ingress_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a) -n ns-a --output json > app_a_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b) -n ns-b --output json > app_b_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin) -n ns-httpbin --output json > httpbin_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json

enable-full-debug:
	@echo "Enabling full debug..."
	istioctl -n istio-ingress pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level debug
	istioctl -n ns-a pc log $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a).ns-a --level debug
	istioctl -n ns-b pc log $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b).ns-b --level debug
	istioctl -n ns-httpbin pc log $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin).ns-httpbin --level debug

enable-full-info:
	@echo "Enabling full info..."
	istioctl -n istio-ingress pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level info
	istioctl -n ns-a pc log $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a).ns-a --level info
	istioctl -n ns-b pc log $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b).ns-b --level info
	istioctl -n ns-httpbin pc log $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin).ns-httpbin --level info

enable-wasm-debug:
	@echo "Enabling WASM debug..."
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level wasm:debug
	istioctl pc log $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a).ns-a --level wasm:debug
	istioctl pc log $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b).ns-b --level wasm:debug
	istioctl pc log $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin).ns-httpbin --level wasm:debug

curl:
	@echo "Single request..."
	@echo curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-a\" \
		--resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\"
	@echo curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-b\" \
		--resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\"
	@echo "Execute in a loop..."
	@echo "while true; do curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-a\" \
		 --resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\" ; sleep 1 ; done"
	@echo "while true; do curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-a\" \
		 --resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\" ; sleep 1 ; done"
