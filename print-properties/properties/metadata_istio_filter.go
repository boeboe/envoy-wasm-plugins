// Helper function and structs to parse istio filter metadata
package properties

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Metadata provides additional inputs to filters based on matched listeners,
// filter chains, routes and endpoints. It is structured as a map, usually from
// filter name (in reverse DNS format) to metadata specific to the filter. Metadata
// key-values for a filter are merged as connection and request handling occurs,
// with later values for the same key overriding earlier values
//
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#config-core-v3-metadata
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
		proxywasm.LogWarnf("failed reading configuration attribute %v.config: %v", strings.Join(path, "."), err)
		result.config = ""
	} else {
		result.config = config
	}

	services, err := getPropertyByteSliceSlice(append(path, "services"))
	if err != nil {
		proxywasm.LogWarnf("failed reading configuration attribute %v.services: %v", strings.Join(path, "."), err)
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
