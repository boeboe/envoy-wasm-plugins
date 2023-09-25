// Helper function to retreive node metadata properties from the node wasm property
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node
package properties

import "github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"

func GetNodeMetadataAnnotations() map[string]string {
	annotations, err := getPropertyStringMap([]string{"node", "metadata", "ANNOTATIONS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ANNOTATIONS: %v", err)
		return make(map[string]string)
	}
	return annotations
}

func GetNodeMetadataAppContainers() string {
	appContainers, err := getPropertyString([]string{"node", "metadata", "APP_CONTAINERS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.APP_CONTAINERS: %v", err)
		return ""
	}
	return appContainers
}

func GetNodeMetadataClusterId() string {
	clusterId, err := getPropertyString([]string{"node", "metadata", "CLUSTER_ID"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.CLUSTER_ID: %v", err)
		return ""
	}
	return clusterId
}

func GetNodeMetadataEnvoyPrometheusPort() int {
	envoyPrometheusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_PROMETHEUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ENVOY_PROMETHEUS_PORT: %v", err)
		return 0
	}
	return int(envoyPrometheusPortFloat64)
}

func GetNodeMetadataEnvoyStatusPort() int {
	envoyStatusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_STATUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ENVOY_STATUS_PORT: %v", err)
		return 0
	}
	return int(envoyStatusPortFloat64)
}

func GetNodeMetadataInstanceIps() string {
	instanceIps, err := getPropertyString([]string{"node", "metadata", "INSTANCE_IPS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.INSTANCE_IPS: %v", err)
		return ""
	}
	return instanceIps
}

func GetNodeMetadataInterceptionMode() string {
	interceptionMode, err := getPropertyString([]string{"node", "metadata", "INTERCEPTION_MODE"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.INTERCEPTION_MODE: %v", err)
		return ""
	}
	return interceptionMode
}

func GetNodeMetadataIstioProxySha() string {
	istioProxySha, err := getPropertyString([]string{"node", "metadata", "ISTIO_PROXY_SHA"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ISTIO_PROXY_SHA: %v", err)
		return ""
	}
	return istioProxySha
}

func GetNodeMetadataIstioVersion() string {
	istioVersion, err := getPropertyString([]string{"node", "metadata", "ISTIO_VERSION"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.ISTIO_VERSION: %v", err)
		return ""
	}
	return istioVersion
}

func GetNodeMetadataLabels() map[string]string {
	labels, err := getPropertyStringMap([]string{"node", "metadata", "LABELS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.LABELS: %v", err)
		return make(map[string]string)
	}
	return labels
}

func GetNodeMetadataMeshId() string {
	meshId, err := getPropertyString([]string{"node", "metadata", "MESH_ID"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.MESH_ID: %v", err)
		return ""
	}
	return meshId
}

func GetNodeMetadataName() string {
	name, err := getPropertyString([]string{"node", "metadata", "NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NAME: %v", err)
		return ""
	}
	return name
}

func GetNodeMetadataNamespace() string {
	namespace, err := getPropertyString([]string{"node", "metadata", "NAMESPACE"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NAMESPACE: %v", err)
		return ""
	}
	return namespace
}

func GetNodeMetadataNodeName() string {
	nodeName, err := getPropertyString([]string{"node", "metadata", "NODE_NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.NODE_NAME: %v", err)
		return ""
	}
	return nodeName
}

func GetNodeMetadataOwner() string {
	owner, err := getPropertyString([]string{"node", "metadata", "OWNER"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.OWNER: %v", err)
		return ""
	}
	return owner
}

func GetNodeMetadataPilotSan() []string {
	pilotSan, err := getPropertyStringSlice([]string{"node", "metadata", "PILOT_SAN"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PILOT_SAN: %v", err)
		return make([]string, 0)
	}
	return pilotSan
}

func GetNodeMetadataPodPorts() string {
	podPorts, err := getPropertyString([]string{"node", "metadata", "POD_PORTS"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.POD_PORTS: %v", err)
		return ""
	}
	return podPorts
}

func GetNodeMetadataServiceAccount() string {
	serviceAccount, err := getPropertyString([]string{"node", "metadata", "SERVICE_ACCOUNT"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.SERVICE_ACCOUNT: %v", err)
		return ""
	}
	return serviceAccount
}

func GetNodeMetadataWorkloadName() string {
	workloadName, err := getPropertyString([]string{"node", "metadata", "WORKLOAD_NAME"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.WORKLOAD_NAME: %v", err)
		return ""
	}
	return workloadName
}

func GetNodeProxyConfigBinaryPath() string {
	binaryPath, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "binaryPath"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.binaryPath: %v", err)
		return ""
	}
	return binaryPath
}

func GetNodeProxyConfigConcurrency() int {
	concurrencyFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "concurrency"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.concurrency: %v", err)
		return 0
	}
	return int(concurrencyFloat64)
}

func GetNodeProxyConfigConfigPath() string {
	configPath, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "configPath"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.configPath: %v", err)
		return ""
	}
	return configPath
}

func GetNodeProxyConfigControlPlaneAuthPolicy() string {
	controlPlaneAuthPolicy, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "controlPlaneAuthPolicy"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.controlPlaneAuthPolicy: %v", err)
		return ""
	}
	return controlPlaneAuthPolicy
}

func GetNodeProxyConfigDiscoveryAddress() string {
	discoveryAddress, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "discoveryAddress"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.discoveryAddress: %v", err)
		return ""
	}
	return discoveryAddress
}

func GetNodeProxyConfigDrainDuration() string {
	drainDuration, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "drainDuration"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.drainDuration: %v", err)
		return ""
	}
	return drainDuration
}

func GetNodeProxyConfigExtraStatTags() []string {
	extraStatTags, err := getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "extraStatTags"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.extraStatTags: %v", err)
		return make([]string, 0)
	}
	return extraStatTags
}

func GetNodeProxyConfigProxyAdminPort() int {
	proxyAdminPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "proxyAdminPort"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.proxyAdminPort: %v", err)
		return 0
	}
	return int(proxyAdminPortFloat64)
}

func GetNodeProxyConfigServiceCluster() string {
	serviceCluster, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "serviceCluster"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.serviceCluster: %v", err)
		return ""
	}
	return serviceCluster
}

func GetNodeProxyConfigStatNameLength() int {
	statNameLength, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statNameLength"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.statNameLength: %v", err)
		return 0
	}
	return int(statNameLength)
}

func GetNodeProxyConfigStatusPort() int {
	statusPort, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statusPort"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.statusPort: %v", err)
		return 0
	}
	return int(statusPort)
}

func GetNodeProxyConfigTerminationDrainDuration() string {
	terminationDrainDuration, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "terminationDrainDuration"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.terminationDrainDuration: %v", err)
		return ""
	}
	return terminationDrainDuration
}

func GetNodeProxyConfigTracingZipkinAddress() string {
	tracingZipkinAddress, err := getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "tracing", "zipkin", "address"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.metadata.PROXY_CONFIG.tracing.zipkin.address: %v", err)
		return ""
	}
	return tracingZipkinAddress
}
