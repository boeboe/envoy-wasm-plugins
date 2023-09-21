#!/usr/bin/env bash


K8S_VERSION=1.27.3
CLUSTER_NAME=istio-wasm
ISTIO_VERSION=1.18.2
NODE_NAME=${CLUSTER_NAME}-control-plane

KUBECTL="kubectl --context kind-${CLUSTER_NAME}"
HELM="helm --kube-context kind-${CLUSTER_NAME}"

if [[ $1 = "kind-up" ]]; then

	if ! $(kind get clusters | grep -q ${CLUSTER_NAME}) ; then 
		kind create cluster --name ${CLUSTER_NAME} --image kindest/node:v${K8S_VERSION}
	else 
		echo "Kind cluster ${CLUSTER_NAME} already exists"
	fi
	k8s_api_ip=$(${KUBECTL} -n kube-system get pod -l component=kube-apiserver -o=jsonpath="{.items[0].metadata.annotations.kubeadm\.kubernetes\.io/kube-apiserver\.advertise-address\.endpoint}")
	export metallb_startip=$(echo ${k8s_api_ip} | awk -F '.' "{ print \$1\".\"\$2\".\"\$3\".100\";}")
	export metallb_stopip=$(echo ${k8s_api_ip} | awk -F '.' "{ print \$1\".\"\$2\".\"\$3\".200\";}")
	${KUBECTL} apply -f kubernetes/metallb-0.12.1.yaml
	${KUBECTL} apply -f - << EOF
$(envsubst < kubernetes/metallb-poolconfig.yaml)
EOF

	${KUBECTL} label node ${NODE_NAME} topology.kubernetes.io/region=region1 --overwrite=true ;
	${KUBECTL} label node ${NODE_NAME} topology.kubernetes.io/zone=zone1a --overwrite=true ;
	${KUBECTL} label node ${NODE_NAME} topology.istio.io/subzone=subzone1a1 --overwrite=true ;
  exit 0
fi


if [[ $1 = "istio-up" ]]; then
	if ! $(${HELM} ls -n istio-system | grep -q istio-base) ; then 
		${HELM} install istio-base tetratelabs/base -n istio-system --create-namespace --version ${ISTIO_VERSION} --wait
	else
		${HELM} upgrade istio-base tetratelabs/base -n istio-system --version ${ISTIO_VERSION} --wait
	fi
	if ! $(${HELM} ls -n istio-system | grep -q istiod) ; then 
		${HELM} install istiod tetratelabs/istiod -n istio-system --version ${ISTIO_VERSION} --wait
	else
		${HELM} upgrade istiod tetratelabs/istiod -n istio-system --version ${ISTIO_VERSION} --wait
	fi
	if ! $(${HELM} ls -n istio-ingress | grep -q istio-ingress) ; then 
		${HELM} install istio-ingress tetratelabs/gateway -n istio-ingress --create-namespace --version ${ISTIO_VERSION} --wait \
			--set labels.istio-locality=region1
	else
		${HELM} upgrade istio-ingress tetratelabs/gateway -n istio-ingress --version ${ISTIO_VERSION} --wait \
			--set labels.istio-locality=region1
	fi

  exit 0
fi


if [[ $1 = "kind-down" ]]; then
  kind delete cluster --name ${CLUSTER_NAME}
  exit 0
fi

echo "please specify action ./make.sh kind-up/istio-up/kind-down"
exit 1