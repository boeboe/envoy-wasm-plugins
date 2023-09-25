// Helper function to retreive wasm properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes
package properties

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get plugin name
// This matches <metadata.name>.<metadata.namespace> in the istio WasmPlugin CR
func GetPluginName() string {
	pluginName, err := getPropertyString([]string{"plugin_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute plugin_name: %v", err)
		return ""
	}
	return pluginName
}

// Get plugin root id
// This matches the <spec.pluginName> in the istio WasmPlugin CR
func GetPluginRootId() string {
	pluginRootId, err := getPropertyString([]string{"plugin_root_id"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute plugin_root_id: %v", err)
		return ""
	}
	return pluginRootId
}

// Get plugin vm id
//
// TODO: this seems to be always empty?
func GetPluginVmId() string {
	pluginVmId, err := getPropertyString([]string{"plugin_vm_id"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute plugin_vm_id: %v", err)
		return ""
	}
	return pluginVmId
}

// Get upstream cluster name
//
// Example value: "outbound|80||httpbin.org"
func GetClusterName() string {
	clusterName, err := getPropertyString([]string{"cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute cluster_name: %v", err)
		return ""
	}
	return clusterName
}

// Get route name (only available in the response path, cfr getXdsRouteName())
// This matches the <spec.http.name> in the istio VirtualService CR
func GetRouteName() string {
	routeName, err := getPropertyString([]string{"route_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute route_name: %v", err)
		return ""
	}
	return routeName
}

// Identifies the direction of the traffic relative to the local Envoy
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#enum-config-core-v3-trafficdirection
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

// Get listener direction, an enum value of the listener traffic direction
//
// Possible values are:
//   - UNSPECIFIED: 0 (default option is unspecified)
//   - INBOUND: 1 (⁣the transport is used for incoming traffic)
//   - OUTBOUND: 2 (the transport is used for outgoing traffic)
func GetListenerDirection() trafficDirection {
	listenerDirection, err := getPropertyUint64([]string{"listener_direction"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute listener_direction: %v", err)
		return 0
	}

	return trafficDirection(int(listenerDirection))
}

// Get an opaque node identifier for the Envoy node. This also provides the local
// service node name. It should be set if any of the following features are used:
// statsd, CDS, and HTTP tracing, either in this message or via --service-node
//
// Example value: router~10.244.0.22~istio-ingress-6d78c67d85-qsbtz.istio-ingress~istio-ingress.svc.cluster.local
func GetNodeId() string {
	nodeId, err := getPropertyString([]string{"node", "id"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute node.id: %v", err)
	}
	return nodeId
}

// Get node cluster, which defines the local service cluster name where envoy
// is running. Though optional, it should be set if any of the following features
// are used: statsd, health check cluster verification, runtime override directory,
// user agent addition, HTTP global rate limiting, CDS, and HTTP tracing, either
// in this message or via --service-cluster
//
// Example value: istio-ingress.istio-ingress
func GetNodeCluster() string {
	nodeCluster, err := getPropertyString([]string{"node", "cluster"})
	if err != nil {
		proxywasm.LogWarnf("failed reading wasm attribute node.cluster: %v", err)
	}
	return nodeCluster
}

// Get map from xDS resource type URL to dynamic context parameters. These may vary at
// runtime (unlike other fields in this message). For example, the xDS client may have a
// shared identifier that changes during the lifetime of the xDS client. In Envoy, this
// would be achieved by updating the dynamic context on the Server::Instance’s LocalInfo
// context provider. The shard ID dynamic parameter then appears in this field during
// future discovery requests
func GetNodeDynamicParams() string {
	nodeDynamicParams, err := getPropertyString([]string{"node", "dynamic_parameters", "params"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.dynamic_parameters.params: %v", err)
	}
	return nodeDynamicParams
}

// Identifies location of where either Envoy runs or where upstream hosts run
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-locality
type Locality struct {
	region  string
	zone    string
	subzone string
}

// Get locality specifying where the Envoy instance is running
func GetNodeLocality() Locality {
	result := Locality{}

	region, err := getPropertyString([]string{"node", "locality", "region"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.locality.region: %v", err)
		result.region = ""
	}
	result.region = region

	zone, err := getPropertyString([]string{"node", "locality", "zone"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.locality.zone: %v", err)
	}
	result.zone = zone

	subzone, err := getPropertyString([]string{"node", "locality", "subzone"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.locality.subzone: %v", err)
	}
	result.subzone = subzone

	return result
}

// Get free-form string that identifies the entity requesting config
//
// Example: “envoy” or “grpc”
func GetNodeUserAgentName() string {
	nodeUserAgentName, err := getPropertyString([]string{"node", "user_agent_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.user_agent_name: %v", err)
	}
	return nodeUserAgentName
}

// Get free-form string that identifies the version of the entity requesting config
//
// Example “1.12.2” or “abcd1234”, or “SpecialEnvoyBuild”
// Not used by istio
func GetNodeUserAgentVersion() string {
	nodeUserAgentVersion, err := getPropertyString([]string{"node", "user_agent_version"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.user_agent_version: %v", err)
	}
	return nodeUserAgentVersion
}

// Get structured version of the entity requesting config
func GetNodeUserAgentBuildVersion() map[string]string {
	nodeUserAgentBuildVersion, err := getPropertyStringMap([]string{"node", "user_agent_build_version", "metadata"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.user_agent_build_version: %v", err)
	}
	return nodeUserAgentBuildVersion
}

// Version and identification for an Envoy extension
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-extension
type Extension struct {
	name      string
	category  string
	type_urls []string
}

// Get list of extensions and their versions supported by the node
func GetNodeExtensions() []Extension {
	result := make([]Extension, 0)
	extensionsRawSlice, err := getPropertyByteSliceSlice([]string{"node", "extensions"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.extensions: %v", err)
	}

	for _, extensionRawSlice := range extensionsRawSlice {
		extensionStringSlice := deserializeProtobufToStringSlice(extensionRawSlice)
		extension := Extension{}
		extension.name = string(extensionStringSlice[0])
		extension.category = string(extensionStringSlice[1])
		extenstionTypeUrls := []string{}
		extenstionTypeUrls = append(extenstionTypeUrls, extensionStringSlice[2:]...)
		extension.type_urls = extenstionTypeUrls
		result = append(result, extension)
	}

	return result
}

// Get client feature support list. These are well known features described in the Envoy API
// repository for a given major version of an API. Client features use reverse DNS naming
// scheme, for example "com.acme.feature"
func GetNodeClientFeatures() []string {
	nodeClientFeatures, err := proxywasm.GetProperty([]string{"node", "client_features"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.client_features: %v", err)
	}
	return deserializeProtobufToStringSlice(nodeClientFeatures)
}

// Get known listening ports on the node as a generic hint to the management server for filtering
// listeners to be returned. For example, if there is a listener bound to port 80, the list can
// optionally contain the SocketAddress (0.0.0.0,80). The field is optional and just a hint
//
// Not used by istio
func GetNodeListeningAddresses() []string {
	nodeListeningAddresses, err := getPropertyStringSlice([]string{"node", "listening_addresses"})
	if err != nil {
		proxywasm.LogWarnf("failed reading node.listening_addresses: %v", err)
	}
	return nodeListeningAddresses
}

// Get cluster metadata
func GetClusterMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"node", "cluster_metadata", "filter_metadata", "istio"})
}

// Get listener metadata
func GetListenerMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"node", "listener_metadata", "filter_metadata", "istio"})
}

// Get route metadata
func GetRouteMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"node", "route_metadata", "filter_metadata", "istio"})
}

// Get upstream host metadata
func GetUpstreamHostMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"node", "upstream_host_metadata", "filter_metadata", "istio"})
}
