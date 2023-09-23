// Helper function to print-properties configuration properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#configuration-attributes
package main

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get upstream cluster name
//
// Example value: "outbound|80||httpbin.org"
func getXdsClusterName() string {
	xdsClusterName, err := getPropertyString([]string{"xds", "cluster_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading configuration attribute xds.cluster_name: %v", err)
		return ""
	}
	return xdsClusterName
}

type istioFilterMetadata struct {
	config   string
	services []istioService
}

type istioService struct {
	host      string
	name      string
	namespace string
}

// Helper function to parse filter metadata
//
// Example envoy extract:
//
//	"metadata": {
//		"filter_metadata": {
//		 "istio": {
//			"config": "/apis/networking.istio.io/v1alpha3/namespaces/default/destination-rule/tetrate-dr",
//			"services": [
//			 {
//				"name": "tetrate.io",
//				"host": "tetrate.io",
//				"namespace": "default"
//			 }
//			]
//		 }
//		}
func getIstioFilterMetadata(path []string) istioFilterMetadata {
	result := istioFilterMetadata{}

	config, err := getPropertyString(append(path, "config"))
	if err != nil {
		proxywasm.LogWarnf("error reading configuration attribute %v.config: %v", strings.Join(path, "."), err)
		result.config = ""
	} else {
		result.config = config
	}

	services, err := getPropertyByteSliceSlice(append(path, "services"))
	if err != nil {
		proxywasm.LogWarnf("error reading configuration attribute %v.services: %v", strings.Join(path, "."), err)
	}

	for _, service := range services {
		istioService := istioService{}
		istioServiceMap := deserializeToStringMap(service)
		istioService.host = istioServiceMap["host"]
		istioService.name = istioServiceMap["name"]
		istioService.namespace = istioServiceMap["namespace"]
		result.services = append(result.services, istioService)
	}

	return result
}

// Get upstream cluster metadata
func getXdsClusterMetadata() istioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "cluster_metadata", "filter_metadata", "istio"})
}

// Get upstream route name (available in both the request response path, cfr getRouteName())
// This matches the <spec.http.name> in the istio VirtualService CR
func getXdsRouteName() string {
	xdsRouteName, err := getPropertyString([]string{"xds", "route_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading configuration attribute xds.route_name: %v", err)
		return ""
	}
	return xdsRouteName
}

// Get upstream route metadata
func getXdsRouteMetadata() istioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "route_metadata", "filter_metadata", "istio"})
}

// Get upstream host metadata
func getXdsUpstreamHostMetadata() istioFilterMetadata {
	return getIstioFilterMetadata([]string{"xds", "upstream_host_metadata", "filter_metadata", "istio"})
}

// Get listener filter chain name
func getXdsListenerFilterChainName() string {
	pluginName, err := getPropertyString([]string{"xds", "filter_chain_name"})
	if err != nil {
		proxywasm.LogWarnf("error reading configuration attribute xds.filter_chain_name: %v", err)
		return ""
	}
	return pluginName
}
