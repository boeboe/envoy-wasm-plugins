// Helper function to retreive response properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#response-attributes
package properties

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get response HTTP status code
func GetResponseCode() int {
	responseCode, err := getPropertyUint64([]string{"response", "code"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.code: %v", err)
		return 0
	}
	return int(responseCode)
}

// Get internal response code details (subject to change)
func GetResponseCodeDetails() string {
	responseCodeDetails, err := getPropertyString([]string{"response", "code_details"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.code_details: %v", err)
		return ""
	}
	return responseCodeDetails
}

// Get additional details about the response beyond the standard response code encoded as a bit-vector
func GetResponseFlags() int {
	responseFlags, err := getPropertyUint64([]string{"response", "flags"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.flags: %v", err)
		return 0
	}
	return int(responseFlags)
}

// Get response gRPC status code
func GetResponseGrpcStatusCode() int {
	responseGrpcStatusCode, err := getPropertyUint64([]string{"response", "grpc_status"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.grpc_status: %v", err)
		return 0
	}
	return int(responseGrpcStatusCode)
}

// Get all response headers indexed by the lower-cased header name
func GetResponseHeaders() map[string]string {
	responseHeaders, err := getPropertyStringMap([]string{"response", "headers"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.headers: %v", err)
		return map[string]string{}
	}
	return responseHeaders
}

// Get all response trailers indexed by the lower-cased trailer name
func GetResponseTrailers() map[string]string {
	responseTrailers, err := getPropertyStringMap([]string{"response", "trailers"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.trailers: %v", err)
		return map[string]string{}
	}
	return responseTrailers
}

// Get size of the response body
func GetResponseSize() int {
	responseSize, err := getPropertyUint64([]string{"response", "size"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.size: %v", err)
		return 0
	}
	return int(responseSize)
}

// Get total size of the response including the approximate uncompressed size of the headers and the trailers
func GetResponseTotalSize() int {
	responseTotalSize, err := getPropertyUint64([]string{"response", "total_size"})
	if err != nil {
		proxywasm.LogWarnf("failed reading response attribute response.total_size: %v", err)
		return 0
	}
	return int(responseTotalSize)
}
