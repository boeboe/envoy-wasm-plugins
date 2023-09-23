// Helper function and structs to parse istio filter metadata
package main

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

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
