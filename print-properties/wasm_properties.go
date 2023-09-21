// Helper function to print-properties wasm properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
package main

import (
	"encoding/binary"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get plugin name
// This matches <metadata.name>.<metadata.namespace> in the istio WasmPlugin CR
func getPluginName() string {
	pluginName, err := proxywasm.GetProperty([]string{"plugin_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_name: %v", err)
		return ""
	}
	return string(pluginName)
}

// Get plugin root id
// This matches the <spec.pluginName> in the istio WasmPlugin CR
func getPluginRootId() string {
	pluginRootId, err := proxywasm.GetProperty([]string{"plugin_root_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_root_id: %v", err)
		return ""
	}
	return string(pluginRootId)
}

// Get plugin vm id
//
// TODO: this seems to be always empty?
func getPluginVmId() string {
	pluginVmId, err := proxywasm.GetProperty([]string{"plugin_vm_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_vm_id: %v", err)
		return ""
	}
	return string(pluginVmId)
}

// Get upstream cluster name
//
// Example value: "outbound|80||httpbin.org"
func getClusterName() string {
	clusterName, err := proxywasm.GetProperty([]string{"cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute cluster_name: %v", err)
		return ""
	}
	return string(clusterName)
}

// Get route name (only available in the response path)
// This matches the <spec.http.name> in the istio VirtualService CR
func getRouteName() string {
	routeName, err := proxywasm.GetProperty([]string{"route_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute route_name: %v", err)
		return ""
	}
	return string(routeName)
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
	listenerDirection, err := proxywasm.GetProperty([]string{"listener_direction"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute listener_direction: %v", err)
		return 0
	}

	return trafficDirection(binary.LittleEndian.Uint32(listenerDirection))
}

// Get node id, an opaque node identifier for the envoy node. It also provides the local
// service node name
//
// Example value: router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local
func getNodeId() string {
	nodeId, err := proxywasm.GetProperty([]string{"node", "id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute node.id: %v", err)
	}
	return string(nodeId)
}

// Get node cluster, which defines the local service cluster name where envoy is running
//
// Example value: istio-ingress.istio-ingress
func getNodeCluster() string {
	nodeCluster, err := proxywasm.GetProperty([]string{"node", "cluster"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute node.cluster: %v", err)
	}
	return string(nodeCluster)
}

// struct containing node bootstrap configuration
type nodeMetadata struct {
	annotations         string
	clusterId           string
	envoyPrometheusPort int
	instanceIps         string
	istioProxySha       string
	istioVersion        string
	labels              string
	meshId              string
	name                string
	namespace           string
	nodeName            string
	owner               string
	pilotSan            string
	proxyConfig         string
	serviceAccount      string
	unprivilegedPod     bool
	workloadName        string
}

// Get node metadata, extending the node identifier
func getNodeMetadata() nodeMetadata {
	result := nodeMetadata{}

	nodeMetadataArray, err := proxywasm.GetPropertyMap([]string{"node", "metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata: %v", err)
		return result
	}
	nodeMetadataMap := tupletArrayToMap(nodeMetadataArray)

	result.annotations = getStringValueFromMap(nodeMetadataMap, "ANNOTATIONS")
	result.clusterId = getStringValueFromMap(nodeMetadataMap, "CLUSTER_ID")
	result.envoyPrometheusPort = getIntValueFromMap(nodeMetadataMap, "ENVOY_PROMETHEUS_PORT")
	result.instanceIps = getStringValueFromMap(nodeMetadataMap, "INSTANCE_IPS")
	result.istioProxySha = getStringValueFromMap(nodeMetadataMap, "ISTIO_PROXY_SHA")
	result.istioVersion = getStringValueFromMap(nodeMetadataMap, "ISTIO_VERSION")
	result.labels = getStringValueFromMap(nodeMetadataMap, "LABELS")
	result.meshId = getStringValueFromMap(nodeMetadataMap, "MESH_ID")
	result.name = getStringValueFromMap(nodeMetadataMap, "NAME")
	result.namespace = getStringValueFromMap(nodeMetadataMap, "NAMESPACE")
	result.nodeName = getStringValueFromMap(nodeMetadataMap, "NODE_NAME")
	result.owner = getStringValueFromMap(nodeMetadataMap, "OWNER")
	result.pilotSan = getStringValueFromMap(nodeMetadataMap, "PILOT_SAN")
	result.proxyConfig = getStringValueFromMap(nodeMetadataMap, "PROXY_CONFIG")
	result.serviceAccount = getStringValueFromMap(nodeMetadataMap, "SERVICE_ACCOUNT")
	result.unprivilegedPod = getBoolValueFromMap(nodeMetadataMap, "UNPRIVILEGED_POD")
	result.workloadName = getStringValueFromMap(nodeMetadataMap, "WORKLOAD_NAME")

	return result
}

func getNodeDynamicParams() string {
	nodeDynamicParams, err := proxywasm.GetProperty([]string{"node", "dynamic_parameters"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.dynamic_parameters: %v", err)
	}
	return string(nodeDynamicParams)
}

func getNodeLocality() string {
	nodeLocality, err := proxywasm.GetProperty([]string{"node", "locality"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.locality: %v", err)
	}
	return string(nodeLocality)
}

func getNodeUserAgentName() string {
	nodeUserAgentName, err := proxywasm.GetProperty([]string{"node", "user_agent_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_name: %v", err)
	}
	return string(nodeUserAgentName)
}

func getNodeUserAgentVersion() string {
	nodeUserAgentVersion, err := proxywasm.GetProperty([]string{"node", "user_agent_version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_version: %v", err)
	}
	return string(nodeUserAgentVersion)
}

func getNodeUserAgentBuildVersion() string {
	nodeUserAgentBuildVersion, err := proxywasm.GetProperty([]string{"node", "user_agent_build_version"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.user_agent_build_version: %v", err)
	}
	return string(nodeUserAgentBuildVersion)
}

func getNodeExtensions() string {
	nodeExtensions, err := proxywasm.GetProperty([]string{"node", "extensions"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.extensions: %v", err)
	}
	return string(nodeExtensions)
}

func getNodeClientFeatures() string {
	nodeClientFeatures, err := proxywasm.GetProperty([]string{"node", "client_features"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.extensions: %v", err)
	}
	return string(nodeClientFeatures)
}

func getNodeListeningAddresses() string {
	nodeListeningAddresses, err := proxywasm.GetProperty([]string{"node", "listening_addresses"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.listening_addresses: %v", err)
	}
	return string(nodeListeningAddresses)
}

func getClusterMetadata() string {
	clusterMetadata, err := proxywasm.GetProperty([]string{"node", "cluster_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.cluster_metadata: %v", err)
	}
	return string(clusterMetadata)
}

func getListenerMetadata() string {
	listenerMetadata, err := proxywasm.GetProperty([]string{"node", "listener_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.listener_metadata: %v", err)
	}
	return string(listenerMetadata)
}

func getRouteMetadata() string {
	routeMetadata, err := proxywasm.GetProperty([]string{"node", "route_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.route_metadata: %v", err)
	}
	return string(routeMetadata)
}

func getUpstreamHostMetadata() string {
	upstreamHostMetadata, err := proxywasm.GetProperty([]string{"node", "upstream_host_metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.upstream_host_metadata: %v", err)
	}
	return string(upstreamHostMetadata)
}
