// Helper function to print-properties wasm properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
package main

import (
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get plugin name
func getPluginName() string {
	pluginName, err := proxywasm.GetProperty([]string{"plugin_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_name: %v", err)
		return ""
	}
	return string(pluginName)
}

// Get plugin root id
func getPluginRootId() string {
	pluginRootId, err := proxywasm.GetProperty([]string{"plugin_root_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_root_id: %v", err)
		return ""
	}
	return string(pluginRootId)
}

// Get plugin vm id
func getPluginVmId() string {
	pluginVmId, err := proxywasm.GetProperty([]string{"plugin_vm_id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute plugin_vm_id: %v", err)
		return ""
	}
	return string(pluginVmId)
}

// Get upstream cluster name
func getClusterName() string {
	clusterName, err := proxywasm.GetProperty([]string{"cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute cluster_name: %v", err)
		return ""
	}
	return string(clusterName)
}

// Get route name
func getRouteName() string {
	routeName, err := proxywasm.GetProperty([]string{"route_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute route_name: %v", err)
		return ""
	}
	return string(routeName)
}

// Get listener direction
// Possible values are:
//   - UNSPECIFIED: 0
//   - INBOUND: 1
//   - OUTBOUND: 2
func getListenerDirection() int {
	listenerDirectionString, err := proxywasm.GetProperty([]string{"listener_direction"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute listener_direction: %v", err)
		return 0
	}

	listenerDirection, err := strconv.Atoi(string(listenerDirectionString))
	if err != nil {
		proxywasm.LogWarnf("error converting wasm attribute listener_direction to integer: %v", err)
		return 0
	}

	return listenerDirection
}

// get node id, an opaque node identiefier for the envoy node. it also provides the local
// service node name
func getNodeId() string {
	nodeId, err := proxywasm.GetProperty([]string{"node", "id"})
	if err != nil {
		proxywasm.LogWarnf("error reading wasm attribute node.id: %v", err)
	}
	return string(nodeId)
}

// get node cluster, which defines the local service cluster name where envoy is running
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

func getWasmNodeMetadataStringProperty(nodeProperties [][2]string, propertyName string) string {
	for _, value := range nodeProperties {
		if value[0] == propertyName {
			return string(value[1])
		}
	}
	proxywasm.LogWarnf("unable to find map key '%v' in node.metadata: %v", propertyName, nodeProperties)
	return ""
}

func getWasmNodeMetadataIntProperty(nodeProperties [][2]string, propertyName string) int {
	for _, value := range nodeProperties {
		if value[0] == propertyName {
			intValue, err := strconv.Atoi(string(value[1]))
			if err != nil {
				proxywasm.LogWarnf("error converting wasm attribute node.metadata with map key '%v' to integer: %v", value[0], err)
				return 0
			}
			return intValue
		}
	}
	proxywasm.LogWarnf("unable to find map key '%v' in wasm attribute node.metadata: %v", propertyName, nodeProperties)
	return 0
}

func getWasmNodeMetadataBoolProperty(nodeProperties [][2]string, propertyName string) bool {
	for _, value := range nodeProperties {
		if value[0] == propertyName {
			boolValue, err := strconv.ParseBool(string(value[1]))
			if err != nil {
				proxywasm.LogWarnf("error converting wasm attribute node.metadata with map key '%v' to bool: %v", value[0], err)
				return false
			}
			return boolValue
		}
	}
	proxywasm.LogWarnf("unable to find map key '%v' in wasm attribute node.metadata: %v", propertyName, nodeProperties)
	return false
}

// get node metadata, extending the node identifier
func getNodeMetadata() nodeMetadata {
	result := nodeMetadata{}

	nodeMetadataMap, err := proxywasm.GetPropertyMap([]string{"node", "metadata"})
	if err != nil {
		proxywasm.LogWarnf("error reading node.metadata: %v", err)
		return result
	}

	result.annotations = getWasmNodeMetadataStringProperty(nodeMetadataMap, "ANNOTATIONS")
	result.clusterId = getWasmNodeMetadataStringProperty(nodeMetadataMap, "CLUSTER_ID")
	result.envoyPrometheusPort = getWasmNodeMetadataIntProperty(nodeMetadataMap, "ENVOY_PROMETHEUS_PORT")
	result.instanceIps = getWasmNodeMetadataStringProperty(nodeMetadataMap, "INSTANCE_IPS")
	result.istioProxySha = getWasmNodeMetadataStringProperty(nodeMetadataMap, "ISTIO_PROXY_SHA")
	result.istioVersion = getWasmNodeMetadataStringProperty(nodeMetadataMap, "ISTIO_VERSION")
	result.labels = getWasmNodeMetadataStringProperty(nodeMetadataMap, "LABELS")
	result.meshId = getWasmNodeMetadataStringProperty(nodeMetadataMap, "MESH_ID")
	result.name = getWasmNodeMetadataStringProperty(nodeMetadataMap, "NAME")
	result.namespace = getWasmNodeMetadataStringProperty(nodeMetadataMap, "NAMESPACE")
	result.nodeName = getWasmNodeMetadataStringProperty(nodeMetadataMap, "NODE_NAME")
	result.owner = getWasmNodeMetadataStringProperty(nodeMetadataMap, "OWNER")
	result.pilotSan = getWasmNodeMetadataStringProperty(nodeMetadataMap, "PILOT_SAN")
	result.proxyConfig = getWasmNodeMetadataStringProperty(nodeMetadataMap, "PROXY_CONFIG")
	result.serviceAccount = getWasmNodeMetadataStringProperty(nodeMetadataMap, "SERVICE_ACCOUNT")
	result.unprivilegedPod = getWasmNodeMetadataBoolProperty(nodeMetadataMap, "UNPRIVILEGED_POD")
	result.workloadName = getWasmNodeMetadataStringProperty(nodeMetadataMap, "WORKLOAD_NAME")

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
