// Helper function to retreive request properties
// https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#request-attributes
package main

import (
	"time"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// Get the path portion of the URL
func getRequestPath() string {
	requestPath, err := getPropertyString([]string{"request", "path"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.path: %v", err)
		return ""
	}
	return requestPath
}

// Get the path portion of the URL without the query string
func getRequestUrlPath() string {
	requestUrlPath, err := getPropertyString([]string{"request", "url_path"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.url_path: %v", err)
		return ""
	}
	return requestUrlPath
}

// Get the host portion of the URL
func getRequestHost() string {
	requestHost, err := getPropertyString([]string{"request", "host"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.host: %v", err)
		return ""
	}
	return requestHost
}

// Get the scheme portion of the URL e.g. “http”
func getRequestScheme() string {
	requestScheme, err := getPropertyString([]string{"request", "scheme"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.scheme: %v", err)
		return ""
	}
	return requestScheme
}

// Get the request method e.g. “GET”
func getRequestMethod() string {
	requestMethod, err := getPropertyString([]string{"request", "method"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.method: %v", err)
		return ""
	}
	return requestMethod
}

// Get all request headers indexed by the lower-cased header name
func getRequestHeaders() map[string]string {
	requestHeaders, err := getPropertyStringMap([]string{"request", "headers"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.headers: %v", err)
		return map[string]string{}
	}
	return requestHeaders
}

// Get the referer request header
func getRequestReferer() string {
	requestReferer, err := getPropertyString([]string{"request", "referer"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.referer: %v", err)
		return ""
	}
	return requestReferer
}

// Get the user agent request header
func getRequestUserAgent() string {
	requestUserAgent, err := getPropertyString([]string{"request", "useragent"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.useragent: %v", err)
		return ""
	}
	return requestUserAgent
}

// Get the time of the first byte received, approximated to nano-seconds
func getRequestTime() time.Time {
	requestTime, err := getPropertTimestamp([]string{"request", "time"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.time: %v", err)
		return time.Now()
	}
	return requestTime
}

// Get the request ID corresponding to x-request-id header value
func getRequestId() string {
	requestId, err := getPropertyString([]string{"request", "id"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.id: %v", err)
		return ""
	}
	return requestId
}

// Get the request protocol (“HTTP/1.0”, “HTTP/1.1”, “HTTP/2”, or “HTTP/3”)
func getRequestProtocol() string {
	requestProtocol, err := getPropertyString([]string{"request", "protocol"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.protocol: %v", err)
		return ""
	}
	return requestProtocol
}

// Get the query portion of the URL in the format of “name1=value1&name2=value2”
func getRequestQuery() string {
	requestQuery, err := getPropertyString([]string{"request", "query"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.query: %v", err)
		return ""
	}
	return requestQuery
}

// Get the total duration of the request, approximated to nano-seconds
func getRequestDuration() int {
	requestDuration, err := getPropertyUint64([]string{"request", "duration"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.duration: %v", err)
		return 0
	}
	return int(requestDuration)
}

// Get the size of the request body. Content length header is used if available
func getRequestSize() int {
	requestSize, err := getPropertyUint64([]string{"request", "size"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.size: %v", err)
		return 0
	}
	return int(requestSize)
}

// Get the total size of the request including the approximate uncompressed size of the headers
func getRequestTotalSize() int {
	requestTotalSize, err := getPropertyUint64([]string{"request", "total_size"})
	if err != nil {
		proxywasm.LogWarnf("failed reading request attribute request.total_size: %v", err)
		return 0
	}
	return int(requestTotalSize)
}
