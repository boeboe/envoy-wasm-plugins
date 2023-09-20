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

OUT := $(WASM)
TAG ?= print-properties-0.1
HUB ?= docker.io/boeboe

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

all: clean test lint docker-build docker-push deploy

compile: $(OUT)

$(OUT):
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test:
	go test -v -tags=proxytest ./...

lint:
	go run $(LINTER) run --verbose --build-tags proxytest

clean:
	rm -f *.wasm

docker-build: $(OUT)
	docker build --platform linux/amd64 --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)

deploy:
	kubectl apply -f istio/wasmplugin.yaml
	kubectl apply -f istio/httpbin-org.yaml
	kubectl apply -f istio/tetrate-io.yaml

dump-config:
	kubectl exec $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingressgateway) -n istio-ingress -- pilot-agent request GET config_dump > config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json

dump-logs:
	kubectl -n istio-ingress logs $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingressgateway)

enable-wasm-debug:
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingressgateway).istio-ingress --level wasm:debug

curl:
	ingress_ip=$$(kubectl get services --namespace istio-ingress istio-ingressgateway --output jsonpath='{.status.loadBalancer.ingress[0].ip}')
	echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.httpbin.org:80:$${ingress_ip}\" \"http://wasm.httpbin.org/headers\"
	echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.tetrate.io:80:$${ingress_ip}\" \"http://wasm.tetrate.io/\"
	