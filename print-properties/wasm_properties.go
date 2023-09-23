// Helper function to print-properties wasm properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get plugin name
// This matches <metadata.name>.<metadata.namespace> in the istio WasmPlugin CR
func getPluginName() string {
	pluginName, err := getPropertyString([]string{"plugin_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_name: %v", err)
		return ""
	}
	return pluginName
}

// Get plugin root id
// This matches the <spec.pluginName> in the istio WasmPlugin CR
func getPluginRootId() string {
	pluginRootId, err := getPropertyString([]string{"plugin_root_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_root_id: %v", err)
		return ""
	}
	return pluginRootId
}

// Get plugin vm id
//
// TODO: this seems to be always empty?
func getPluginVmId() string {
	pluginVmId, err := getPropertyString([]string{"plugin_vm_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_vm_id: %v", err)
		return ""
	}
	return pluginVmId
}

// Get upstream cluster name
//
// Example value: "outbound|80||httpbin.org"
func getClusterName() string {
	clusterName, err := getPropertyString([]string{"cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute cluster_name: %v", err)
		return ""
	}
	return clusterName
}

// Get route name (only available in the response path)
// This matches the <spec.http.name> in the istio VirtualService CR
func getRouteName() string {
	routeName, err := getPropertyString([]string{"route_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute route_name: %v", err)
		return ""
	}
	return routeName
}

type trafficDirection int

const (
	Unspecified trafficDirection = iota
	Inbound
	Outbound
)

func (t trafficDirection) String() string {
	switch t {
	case Unspecified:
		return "UNSPECIFIED"
	case Inbound:
		return "INBOUND"
	case Outbound:
		return "OUTBOUND"
	}
	return "UNKNOWN"
}

// Get listener direction (enum value)
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#enum-config-core-v3-trafficdirection
//
// Possible values are:
//   - UNSPECIFIED: 0
//   - INBOUND: 1
//   - OUTBOUND: 2
func getListenerDirection() trafficDirection {
	listenerDirection, err := getPropertyUint64([]string{"listener_direction"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute listener_direction: %v", err)
		return 0
	}

	return trafficDirection(int(listenerDirection))
}

// Get node id, an opaque node identifier for the envoy node. It also provides the local
// service node name
//
// Example value: router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local
func getNodeId() string {
	nodeId, err := getPropertyString([]string{"node", "id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute node.id: %v", err)
	}
	return nodeId
}

// Get node cluster, which defines the local service cluster name where envoy is running
//
// Example value: istio-ingress.istio-ingress
func getNodeCluster() string {
	nodeCluster, err := getPropertyString([]string{"node", "cluster"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute node.cluster: %v", err)
	}
	return nodeCluster
}

// struct containing node bootstrap configuration
type nodeMetadata struct {
	annotations         map[string]string
	appContainers       string
	clusterId           string
	envoyPrometheusPort int
	envoyStatusPort     int
	instanceIps         string
	interceptionMode    string
	istioProxySha       string
	istioVersion        string
	labels              map[string]string
	meshId              string
	name                string
	namespace           string
	nodeName            string
	owner               string
	pilotSan            []string
	podPorts            string
	proxyConfig         proxyConfig
	serviceAccount      string
	workloadName        string
}

type proxyConfig struct {
	binaryPath               string
	concurrency              int
	configPath               string
	controlPlaneAuthPolicy   string
	discoveryAddress         string
	drainDuration            string
	extraStatTags            []string
	proxyAdminPort           int
	serviceCluster           string
	statNameLength           int
	statusPort               int
	terminationDrainDuration string
	tracing                  tracing
}

type tracing struct {
	address string
}

// Get node metadata, extending the node identifier
func getNodeMetadata() nodeMetadata {
	result := nodeMetadata{}
	var err error

	result.annotations, err = getPropertyStringMap([]string{"node", "metadata", "ANNOTATIONS"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.ANNOTATIONS: %v", err)
	}
	result.appContainers, err = getPropertyString([]string{"node", "metadata", "APP_CONTAINERS"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.APP_CONTAINERS: %v", err)
	}
	result.clusterId, err = getPropertyString([]string{"node", "metadata", "CLUSTER_ID"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.CLUSTER_ID: %v", err)
	}
	envoyPrometheusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_PROMETHEUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.ENVOY_PROMETHEUS_PORT: %v", err)
	}
	result.envoyPrometheusPort = int(envoyPrometheusPortFloat64)
	envoyStatusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "ENVOY_STATUS_PORT"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.ENVOY_STATUS_PORT: %v", err)
	}
	result.envoyStatusPort = int(envoyStatusPortFloat64)
	result.instanceIps, err = getPropertyString([]string{"node", "metadata", "INSTANCE_IPS"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.INSTANCE_IPS: %v", err)
	}
	result.interceptionMode, err = getPropertyString([]string{"node", "metadata", "INTERCEPTION_MODE"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.INTERCEPTION_MODE: %v", err)
	}
	result.istioProxySha, err = getPropertyString([]string{"node", "metadata", "ISTIO_PROXY_SHA"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.ISTIO_PROXY_SHA: %v", err)
	}
	result.istioVersion, err = getPropertyString([]string{"node", "metadata", "ISTIO_VERSION"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.ISTIO_VERSION: %v", err)
	}
	result.labels, err = getPropertyStringMap([]string{"node", "metadata", "LABELS"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.LABELS: %v", err)
	}
	result.meshId, err = getPropertyString([]string{"node", "metadata", "MESH_ID"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.MESH_ID: %v", err)
	}
	result.name, err = getPropertyString([]string{"node", "metadata", "NAME"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.NAME: %v", err)
	}
	result.namespace, err = getPropertyString([]string{"node", "metadata", "NAMESPACE"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.NAMESPACE: %v", err)
	}
	result.nodeName, err = getPropertyString([]string{"node", "metadata", "NODE_NAME"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.NODE_NAME: %v", err)
	}
	result.owner, err = getPropertyString([]string{"node", "metadata", "OWNER"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.OWNER: %v", err)
	}
	result.pilotSan, err = getPropertyStringSlice([]string{"node", "metadata", "PILOT_SAN"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PILOT_SAN: %v", err)
	}
	result.podPorts, err = getPropertyString([]string{"node", "metadata", "POD_PORTS"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.POD_PORTS: %v", err)
	}
	result.serviceAccount, err = getPropertyString([]string{"node", "metadata", "SERVICE_ACCOUNT"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.SERVICE_ACCOUNT: %v", err)
	}
	result.workloadName, err = getPropertyString([]string{"node", "metadata", "WORKLOAD_NAME"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.WORKLOAD_NAME: %v", err)
	}

	result.proxyConfig = proxyConfig{}
	result.proxyConfig.binaryPath, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "binaryPath"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.binaryPath: %v", err)
	}
	proxyConfigConcurrencyFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "concurrency"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.concurrency: %v", err)
	}
	result.proxyConfig.concurrency = int(proxyConfigConcurrencyFloat64)
	result.proxyConfig.configPath, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "configPath"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.configPath: %v", err)
	}
	result.proxyConfig.controlPlaneAuthPolicy, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "controlPlaneAuthPolicy"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.controlPlaneAuthPolicy: %v", err)
	}
	result.proxyConfig.discoveryAddress, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "discoveryAddress"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.discoveryAddress: %v", err)
	}
	result.proxyConfig.drainDuration, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "drainDuration"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.drainDuration: %v", err)
	}
	result.proxyConfig.extraStatTags, err = getPropertyStringSlice([]string{"node", "metadata", "PROXY_CONFIG", "extraStatTags"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.extraStatTags: %v", err)
	}
	proxyConfigProxyAdminPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "proxyAdminPort"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.proxyAdminPort: %v", err)
	}
	result.proxyConfig.proxyAdminPort = int(proxyConfigProxyAdminPortFloat64)
	result.proxyConfig.serviceCluster, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "serviceCluster"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.serviceCluster: %v", err)
	}
	proxyConfigStatNameLengthFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statNameLength"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.statNameLength: %v", err)
	}
	result.proxyConfig.statNameLength = int(proxyConfigStatNameLengthFloat64)
	proxyConfigStatusPortFloat64, err := getPropertyFloat64([]string{"node", "metadata", "PROXY_CONFIG", "statusPort"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.statusPort: %v", err)
	}
	result.proxyConfig.statusPort = int(proxyConfigStatusPortFloat64)
	result.proxyConfig.terminationDrainDuration, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "terminationDrainDuration"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.terminationDrainDuration: %v", err)
	}
	result.proxyConfig.tracing.address, err = getPropertyString([]string{"node", "metadata", "PROXY_CONFIG", "tracing", "zipkin", "address"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata.PROXY_CONFIG.tracing.zipkin.address: %v", err)
	}

	return result
}

func getNodeDynamicParams() string {
	nodeDynamicParams, err := getPropertyString([]string{"node", "dynamic_parameters", "params"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.dynamic_parameters.params: %v", err)
	}
	return nodeDynamicParams
}

func getNodeLocality() string {
	nodeLocality, err := getPropertyString([]string{"node", "locality", "region"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.locality.region: %v", err)
	}
	return nodeLocality
}

func getNodeUserAgentName() string {
	nodeUserAgentName, err := getPropertyString([]string{"node", "user_agent_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_name: %v", err)
	}
	return nodeUserAgentName
}

// Not used by istio
func getNodeUserAgentVersion() string {
	nodeUserAgentVersion, err := getPropertyString([]string{"node", "user_agent_version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_version: %v", err)
	}
	return nodeUserAgentVersion
}

func getNodeUserAgentBuildVersion() map[string]string {
	nodeUserAgentBuildVersion, err := getPropertyStringMap([]string{"node", "user_agent_build_version", "metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_build_version: %v", err)
	}
	return nodeUserAgentBuildVersion
}

func getNodeExtensions() string {
	nodeExtensions, err := getPropertyString([]string{"node", "extensions"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.extensions: %v", err)
	}
	return nodeExtensions
}

func getNodeClientFeatures() string {
	nodeClientFeatures, err := getPropertyString([]string{"node", "client_features"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.extensions: %v", err)
	}
	return nodeClientFeatures
}

func getNodeListeningAddresses() string {
	nodeListeningAddresses, err := getPropertyString([]string{"node", "listening_addresses"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.listening_addresses: %v", err)
	}
	return nodeListeningAddresses
}

func getClusterMetadata() string {
	clusterMetadata, err := getPropertyString([]string{"node", "cluster_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.cluster_metadata: %v", err)
	}
	return clusterMetadata
}

func getListenerMetadata() string {
	listenerMetadata, err := getPropertyString([]string{"node", "listener_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.listener_metadata: %v", err)
	}
	return listenerMetadata
}

func getRouteMetadata() string {
	routeMetadata, err := getPropertyString([]string{"node", "route_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.route_metadata: %v", err)
	}
	return routeMetadata
}

func getUpstreamHostMetadata() string {
	upstreamHostMetadata, err := getPropertyString([]string{"node", "upstream_host_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.upstream_host_metadata: %v", err)
	}
	return upstreamHostMetadata
}
