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
.PHONY: compile docker-build docker-push deploy dump-config enable-wasm-debug curl
.SILENT: compile docker-build docker-push deploy dump-config enable-wasm-debug curl

NAME := envoy-wasm-plugins
WASM := $(NAME).wasm

OUT := $(WASM)
TAG ?= header-propagator-0.1
HUB ?= docker.io/boeboe

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

all: clean test lint docker-build docker-push deploy

compile:
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test:
	go test -v -tags=proxytest ./...

lint:
	go run $(LINTER) run --verbose --build-tags proxytest

clean:
	rm -f *.wasm

docker-build: compile
	docker build --platform linux/amd64 --build-arg WASM_BINARY_PATH=$(OUT) -t $(HUB)/$(NAME):$(TAG) .

docker-push:
	docker push $(HUB)/$(NAME):$(TAG)

deploy:
	kubectl apply -f istio/zipkin.yaml
	kubectl apply -f istio/gateway.yaml
	kubectl apply -f istio/app-a.yaml
	kubectl apply -f istio/app-b.yaml
	kubectl apply -f istio/httpbin.yaml
	kubectl apply -f istio/wasmplugin.yaml
	kubectl delete pod $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a) -n ns-a
	kubectl delete pod $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b) -n ns-b
	kubectl delete pod $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin) -n ns-httpbin

dump-config:
	kubectl exec $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress -- pilot-agent request GET config_dump > ingress_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a) -n ns-a --output json > app_a_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b) -n ns-b --output json > app_b_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json
	istioctl pc dump $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin) -n ns-httpbin --output json > httpbin_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json

dump-logs:
	kubectl -n istio-ingress logs $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) > ingress_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log
	kubectl -n ns-a logs $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a) -c istio-proxy > app_a_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log
	kubectl -n ns-b logs $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b) -c istio-proxy > app_b_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log
	kubectl -n ns-httpbin logs $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin) -c istio-proxy > httpbin_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log

enable-full-debug:
	istioctl -n istio-ingress pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level debug
	istioctl -n ns-a pc log $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a).ns-a --level debug
	istioctl -n ns-b pc log $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b).ns-b --level debug
	istioctl -n ns-httpbin pc log $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin).ns-httpbin --level debug

enable-full-info:
	istioctl -n istio-ingress pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level info
	istioctl -n ns-a pc log $$(kubectl get pods -n ns-a -o custom-columns=:metadata.name -l app=app-a).ns-a --level info
	istioctl -n ns-b pc log $$(kubectl get pods -n ns-b -o custom-columns=:metadata.name -l app=app-b).ns-b --level info
	istioctl -n ns-httpbin pc log $$(kubectl get pods -n ns-httpbin -o custom-columns=:metadata.name -l app=httpbin).ns-httpbin --level info

enable-wasm-debug:
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
	
	
	@echo "Executing in a loop..."
	@echo "while true; do curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-a\" \
		 --resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\" ; sleep 1 ; done"
	@echo "while true; do curl -v -H \"X-B3-Sampled: 1\" -H \"X-B3-Traceid: $$(openssl rand -hex 16)\" -H \"X-Tetrate-Swimlaneid: lane-a\" \
		 --resolve \"propagate.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers\" ; sleep 1 ; done"


