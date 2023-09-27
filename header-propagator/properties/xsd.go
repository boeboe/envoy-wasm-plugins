// Helper function to retreive xsd configuration properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#configuration-attributes
package properties

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get upstream cluster name
//
// Example value: "outbound|80||httpbin.org"
func GetXdsClusterName() string {
	xdsClusterName, err := getPropertyString([]string{"xds", "cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading xsd configuration attribute xds.cluster_name: %v", err)
		return ""
	}
	return xdsClusterName
}

// Get upstream cluster metadata
func GetXdsClusterMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "cluster_metadata", "filter_metadata", "istio"})
}

// Get upstream route name (available in both the request response path, cfr getRouteName())
// This matches the <spec.http.name> in the istio VirtualService CR
func GetXdsRouteName() string {
	xdsRouteName, err := getPropertyString([]string{"xds", "route_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading xsd configuration attribute xds.route_name: %v", err)
		return ""
	}
	return xdsRouteName
}

// Get upstream route metadata
func GetXdsRouteMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "route_metadata", "filter_metadata", "istio"})
}

// Get upstream host metadata
func GetXdsUpstreamHostMetadata() IstioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "upstream_host_metadata", "filter_metadata", "istio"})
}

// Get listener filter chain name
func GetXdsListenerFilterChainName() string {
	pluginName, err := getPropertyString([]string{"xds", "filter_chain_name"})
	if err != nil {
		proxywasm.LogWarnf("failed reading xsd configuration attribute xds.filter_chain_name: %v", err)
		return ""
	}
	return pluginName
}
