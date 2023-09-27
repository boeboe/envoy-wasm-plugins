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
TAG  ?= geo-tagger-0.1
HUB  ?= docker.io/boeboe
REPO ?= boeboe/envoy-wasm-plugins

LINTER := github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0

# all: clean test lint docker-build docker-push deploy
all: clean test lint release deploy

compile: $(OUT)

$(OUT):
	tinygo build -o $(OUT) -scheduler=none -target=wasi ./...

test:
	go test -v -tags=proxytest ./...

lint:
	go run $(LINTER) run --verbose --build-tags proxytest

clean:
	rm -f *.wasm

release-local: $(OUT)
	@echo "Releasing $(TAG) into the local wasm-repo..."
	curl --resolve "wasm-repo.tetrate.io:80:$$(kubectl get services --namespace wasm-repo wasm-repo --output jsonpath='{.status.loadBalancer.ingress[0].ip}')" -X POST -F "file=@$(OUT)" http://wasm-repo.tetrate.io/wasm-plugins/$(TAG).wasm
	@echo "\nAvailable plugins: "
	curl --resolve "wasm-repo.tetrate.io:80:$$(kubectl get services --namespace wasm-repo wasm-repo --output jsonpath='{.status.loadBalancer.ingress[0].ip}')" http://wasm-repo.tetrate.io/list

release-github: $(OUT)
	@echo "Checking if release $(TAG) already exists..."
	@if gh release view $(TAG) -R $(REPO) > /dev/null 2>&1; then \
		echo "Release $(TAG) exists. Deleting it..."; \
		gh release delete $(TAG) -R $(REPO) --yes; \
	fi
	@echo "Creating a new release on GitHub..."
	gh release create $(TAG) $(WASM) --title "Release $(TAG)" --notes "Release notes for $(TAG)" --repo $(REPO)

deploy:
	if ! $$(kubectl get configmap istio -n istio-system -o yaml | grep -q "BOOTSTRAP_XDS_AGENT") ; then echo "patch/merge istio configmap manually with ./istio/istio-configmap-patch.yaml content" ; exit 1 ; fi
	wasm_plugin_sha256=$$(sha256sum ${WASM} | awk '{print $$1}') envsubst < istio/envoyfilter.template.yaml > istio/envoyfilter.yaml
	kubectl apply -f istio/envoyfilter.yaml
	kubectl delete pod $$(kubectl get pods -n istio-system -o custom-columns=:metadata.name -l istio=pilot) -n istio-system
	kubectl delete pod $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress
	sleep 5
	kubectl wait --for=condition=Ready pod -l istio=pilot -n istio-system --timeout=300s
	kubectl wait --for=condition=Ready pod -l istio=ingress -n istio-ingress --timeout=300s

dump-config:
	kubectl exec $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) -n istio-ingress -- pilot-agent request GET config_dump > ingress_config_dump-$$(date "+%Y.%m.%d-%H.%M.%S").json

dump-logs:
	kubectl -n istio-ingress logs $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress) > ingress_logs-$$(date "+%Y.%m.%d-%H.%M.%S").log

enable-wasm-debug:
	istioctl pc log $$(kubectl get pods -n istio-ingress -o custom-columns=:metadata.name -l istio=ingress).istio-ingress --level wasm:debug

curl:
	echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.httpbin.org:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://wasm.httpbin.org/headers\"
	echo curl -v -H \"X-B3-Sampled: 1\" --resolve \"wasm.tetrate.io:80:$$(kubectl get services --namespace istio-ingress istio-ingress --output jsonpath='{.status.loadBalancer.ingress[0].ip}')\" \"http://wasm.tetrate.io/\"
	