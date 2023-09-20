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
.PHONY: kind-up istio-up kind-down
.SILENT: kind-up istio-up kind-down


K8S_VERSION := 1.27.3
CLUSTER_NAME := istio-wasm
ISTIO_VERSION := 1.18.2

KUBECTL := kubectl --context kind-${CLUSTER_NAME}
HELM := helm --kube-context kind-${CLUSTER_NAME}

up: istio-up ## Bring up full demo scenario
clean: kind-down ## Bring down full demo scenario


kind-up: ## Spin up a local kind cluster
	if ! $$(kind get clusters | grep -q ${CLUSTER_NAME}) ; then 
		kind create cluster --name ${CLUSTER_NAME} --image kindest/node:v${K8S_VERSION} ; 
	else 
		echo "Kind cluster ${CLUSTER_NAME} already exists" ; 
	fi
	k8s_api_ip=$$(${KUBECTL} -n kube-system get pod -l component=kube-apiserver -o=jsonpath="{.items[0].metadata.annotations.kubeadm\.kubernetes\.io/kube-apiserver\.advertise-address\.endpoint}")
	export metallb_startip=$$(echo $${k8s_api_ip} | awk -F '.' "{ print \$$1\".\"\$$2\".\"\$$3\".100\";}")
	export metallb_stopip=$$(echo $${k8s_api_ip} | awk -F '.' "{ print \$$1\".\"\$$2\".\"\$$3\".200\";}")
	if ! $$(${HELM} ls -n metallb-system | grep -q metallb) ; then 
		${HELM} install metallb metallb/metallb -n metallb-system --create-namespace
	else
		${HELM} upgrade metallb metallb/metallb -n metallb-system
	fi
	${KUBECTL} --context kind-${CLUSTER_NAME} apply -f - << EOF
	$$(envsubst < kubernetes/metallb-poolconfig.yaml)
	EOF


istio-up: kind-up ## Install istio in the kind cluster
	if ! $$(${HELM} ls -n istio-system | grep -q istio-base) ; then 
		${HELM} install istio-base tetratelabs/base -n istio-system --create-namespace --version ${ISTIO_VERSION}
	else
		${HELM} upgrade istio-base tetratelabs/base -n istio-system --version ${ISTIO_VERSION}
	fi
	if ! $$(${HELM} ls -n istio-system | grep -q istiod) ; then 
		${HELM} install istiod tetratelabs/istiod -n istio-system --version ${ISTIO_VERSION}
	else
		${HELM} upgrade istiod tetratelabs/istiod -n istio-system --version ${ISTIO_VERSION}
	fi
	if ! $$(${HELM} ls -n istio-ingress | grep -q istio-ingress) ; then 
		${HELM} install istio-ingress tetratelabs/istio-ingress -n istio-ingress --create-namespace --version ${ISTIO_VERSION}
	else
		${HELM} upgrade istio-ingress tetratelabs/istio-ingress -n istio-ingress --version ${ISTIO_VERSION}
	fi


kind-down: ## Remove the local kind cluster
	kind delete cluster --name ${CLUSTER_NAME}