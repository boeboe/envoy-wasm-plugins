package main

import (
	"print-properties/properties"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

// main initializes the VM context.
func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// NewPluginContext is used for creating PluginContext for each plugin configuration.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	pluginConfig pluginConfig
}

type pluginConfig struct {
	OnPluginStart          propertiesPrinting `json:"onPluginStart"`
	OnHttpRequestHeaders   propertiesPrinting `json:"onHttpRequestHeaders"`
	OnHttpRequestBody      propertiesPrinting `json:"onHttpRequestBody"`
	OnHttpRequestTrailers  propertiesPrinting `json:"onHttpRequestTrailers"`
	OnHttpResponseHeaders  propertiesPrinting `json:"onHttpResponseHeaders"`
	OnHttpResponseBody     propertiesPrinting `json:"onHttpResponseBody"`
	OnHttpResponseTrailers propertiesPrinting `json:"onHttpResponseTrailers"`
	OnHttpStreamDone       propertiesPrinting `json:"onHttpStreamDone"`
}

type propertiesPrinting struct {
	PrintWasmProperties            bool `json:"printWasmProperties"`
	PrintNodeMetadataProperties    bool `json:"printNodeMetadataProperties"`
	PrintNodeProxyConfigProperties bool `json:"printNodeProxyConfigProperties"`
	PrintXdsProperties             bool `json:"printXdsProperties"`
	PrintUpstreamProperties        bool `json:"printUpstreamProperties"`
	PrintConnectionProperties      bool `json:"printConnectionProperties"`
	PrintResponseProperties        bool `json:"printResponseProperties"`
	PrintRequestProperties         bool `json:"printRequestProperties"`
}

// OnPluginStart is called for all plugin contexts (after OnVmStart if this is the VM context).
// During this call, GetPluginConfiguration is available and can be used to
// retrieve the configuration set at config.configuration in the host configuration.
func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	if pluginConfigurationSize == 0 {
		return types.OnPluginStartStatusOK
	}

	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogErrorf("Failed to get plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	if !gjson.ValidBytes(configData) {
		proxywasm.LogError("Invalid JSON configuration")
		return types.OnPluginStartStatusFailed
	}

	p.pluginConfig = parseConfigData(configData)
	if anyPrintBoolTrue(p.pluginConfig.OnPluginStart) {
		proxywasm.LogInfo("********** OnPluginStart **********")
		printProperties(p.pluginConfig.OnPluginStart)
	}
	return types.OnPluginStartStatusOK
}

// NewHttpContext is used for creating HttpContext for each Http stream.
// Return nil to indicate this PluginContext is not for HttpContext
func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID:    contextID,
		pluginConfig: p.pluginConfig,
	}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID    uint32
	pluginConfig pluginConfig
}

// *********************************************
// REQUEST PATH
// *********************************************

// OnHttpRequestHeaders is called when request headers arrive.
// Return types.ActionPause if you want to stop sending headers to the upstream.
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpRequestHeaders) {
		proxywasm.LogInfo("********** OnHttpRequestHeaders **********")
		printProperties(ctx.pluginConfig.OnHttpRequestHeaders)
	}
	return types.ActionContinue
}

// OnHttpRequestBody is called when a request body *frame* arrives.
// Note that this is potentially called multiple times until we see end_of_stream = true.
// Return types.ActionPause if you want to buffer the body and stop sending body to the upstream.
// Even after returning types.ActionPause, this will be called when an unseen frame arrives.
func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpRequestBody) {
		proxywasm.LogInfo("********** OnHttpRequestBody **********")
		printProperties(ctx.pluginConfig.OnHttpRequestBody)
	}
	return types.ActionContinue
}

// OnHttpRequestTrailers is called when request trailers arrive.
// Return types.ActionPause if you want to stop sending trailers to the upstream.
func (ctx *httpContext) OnHttpRequestTrailers(bodySize int) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpRequestTrailers) {
		proxywasm.LogInfo("********** OnHttpRequestTrailers **********")
		printProperties(ctx.pluginConfig.OnHttpRequestTrailers)
	}
	return types.ActionContinue
}

// *********************************************
// RESPONSE PATH
// *********************************************

// OnHttpResponseHeaders is called when response headers arrive.
// Return types.ActionPause if you want to stop sending headers to downstream.
func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpResponseHeaders) {
		proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
		printProperties(ctx.pluginConfig.OnHttpResponseHeaders)
	}
	return types.ActionContinue
}

// OnHttpResponseBody is called when a response body *frame* arrives.
// Note that this is potentially called multiple times until we see end_of_stream = true.
// Return types.ActionPause if you want to buffer the body and stop sending body to the downtream.
// Even after returning types.ActionPause, this will be called when an unseen frame arrives.
func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpResponseBody) {
		proxywasm.LogInfo("********** OnHttpResponseBody **********")
		printProperties(ctx.pluginConfig.OnHttpResponseBody)
	}
	return types.ActionContinue
}

// OnHttpResponseTrailers is called when response trailers arrive.
// Return types.ActionPause if you want to stop sending trailers to the downstream.
func (ctx *httpContext) OnHttpResponseTrailers(bodySize int) types.Action {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpResponseTrailers) {
		proxywasm.LogInfo("********** OnHttpResponseTrailers **********")
		printProperties(ctx.pluginConfig.OnHttpResponseTrailers)
	}
	return types.ActionContinue
}

// OnHttpStreamDone is called before the host deletes this context.
// You can retrieve the HTTP request/response information (such as headers, etc.) during this call.
// This can be used to implement logging features.
func (ctx *httpContext) OnHttpStreamDone() {
	if anyPrintBoolTrue(ctx.pluginConfig.OnHttpStreamDone) {
		proxywasm.LogInfo("********** OnHttpStreamDone **********")
		printProperties(ctx.pluginConfig.OnHttpStreamDone)
	}
}

// *********************************************
// HELPER FUNCTIONS
// *********************************************

// parseConfigData parses the configuration data into the pluginConfig.
func parseConfigData(data []byte) pluginConfig {
	type eventNames struct {
		name   string
		fields []string
	}
	events := []eventNames{
		{"onPluginStart", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpRequestHeaders", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpRequestBody", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpRequestTrailers", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpResponseHeaders", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpResponseBody", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpResponseTrailers", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
		{"onHttpStreamDone", []string{"printWasmProperties", "printNodeMetadataProperties", "printNodeProxyConfigProperties", "printXdsProperties", "printUpstreamProperties", "printConnectionProperties", "printResponseProperties", "printRequestProperties"}},
	}
	jsonData := gjson.ParseBytes(data)
	config := pluginConfig{}

	for _, event := range events {
		propPrint := propertiesPrinting{}
		for _, field := range event.fields {
			path := event.name + "." + field
			value := jsonData.Get(path).Bool()
			switch field {
			case "printWasmProperties":
				propPrint.PrintWasmProperties = value
			case "printNodeMetadataProperties":
				propPrint.PrintNodeMetadataProperties = value
			case "printNodeProxyConfigProperties":
				propPrint.PrintNodeProxyConfigProperties = value
			case "printXdsProperties":
				propPrint.PrintXdsProperties = value
			case "printUpstreamProperties":
				propPrint.PrintConnectionProperties = value
			case "printConnectionProperties":
				propPrint.PrintNodeProxyConfigProperties = value
			case "printResponseProperties":
				propPrint.PrintResponseProperties = value
			case "printRequestProperties":
				propPrint.PrintRequestProperties = value
			}
		}
		switch event.name {
		case "onPluginStart":
			config.OnPluginStart = propPrint
		case "onHttpRequestHeaders":
			config.OnHttpRequestHeaders = propPrint
		case "onHttpRequestBody":
			config.OnHttpRequestBody = propPrint
		case "onHttpRequestTrailers":
			config.OnHttpRequestTrailers = propPrint
		case "onHttpResponseHeaders":
			config.OnHttpResponseHeaders = propPrint
		case "onHttpResponseBody":
			config.OnHttpResponseBody = propPrint
		case "onHttpResponseTrailers":
			config.OnHttpResponseTrailers = propPrint
		case "onHttpStreamDone":
			config.OnHttpStreamDone = propPrint
		}
	}
	return config
}

// anyPrintBoolTrue checks if any of the properties should be printed.
func anyPrintBoolTrue(p propertiesPrinting) bool {
	return p.PrintWasmProperties ||
		p.PrintNodeMetadataProperties ||
		p.PrintNodeProxyConfigProperties ||
		p.PrintXdsProperties ||
		p.PrintUpstreamProperties ||
		p.PrintConnectionProperties ||
		p.PrintResponseProperties ||
		p.PrintRequestProperties
}

// printProperties logs the properties based on the configuration.
func printProperties(p propertiesPrinting) {
	if p.PrintWasmProperties {
		printWasmProperties()
	}
	if p.PrintNodeMetadataProperties {
		printNodeMetadataProperties()
	}
	if p.PrintNodeProxyConfigProperties {
		printNodeProxyConfigProperties()
	}
	if p.PrintXdsProperties {
		printXdsProperties()
	}
	if p.PrintUpstreamProperties {
		printUpstreamProperties()
	}
	if p.PrintConnectionProperties {
		printConnectionProperties()
	}
	if p.PrintResponseProperties {
		printResponseProperties()
	}
	if p.PrintRequestProperties {
		printRequestProperties()
	}
}

func printWasmProperties() {
	proxywasm.LogInfof(">> GetPluginName: %v", properties.GetPluginName())
	proxywasm.LogInfof(">> GetPluginRootId: %v", properties.GetPluginRootId())
	proxywasm.LogInfof(">> GetPluginVmId: %v", properties.GetPluginVmId())
	proxywasm.LogInfof(">> GetClusterName: %v", properties.GetClusterName())
	proxywasm.LogInfof(">> GetRouteName: %v", properties.GetRouteName())
	proxywasm.LogInfof(">> GetListenerDirection: %v", properties.GetListenerDirection())
	proxywasm.LogInfof(">> GetNodeId: %v", properties.GetNodeId())
	proxywasm.LogInfof(">> GetNodeCluster: %v", properties.GetNodeCluster())
	proxywasm.LogInfof(">> GetNodeDynamicParams: %v", properties.GetNodeDynamicParams())
	proxywasm.LogInfof(">> GetNodeLocality: %+v", properties.GetNodeLocality())
	proxywasm.LogInfof(">> GetNodeUserAgentName: %v", properties.GetNodeUserAgentName())
	proxywasm.LogInfof(">> GetNodeUserAgentVersion: %v", properties.GetNodeUserAgentVersion())
	proxywasm.LogInfof(">> GetNodeUserAgentBuildVersion: %v", properties.GetNodeUserAgentBuildVersion())
	proxywasm.LogInfof(">> GetNodeExtensions: %+v", properties.GetNodeExtensions())
	proxywasm.LogInfof(">> GetNodeClientFeatures: %v", properties.GetNodeClientFeatures())
	proxywasm.LogInfof(">> GetNodeListeningAddresses: %v", properties.GetNodeListeningAddresses())
	proxywasm.LogInfof(">> GetClusterMetadata: %+v", properties.GetClusterMetadata())
	proxywasm.LogInfof(">> GetListenerMetadata: %+v", properties.GetListenerMetadata())
	proxywasm.LogInfof(">> GetRouteMetadata: %+v", properties.GetRouteMetadata())
	proxywasm.LogInfof(">> GetUpstreamHostMetadata: %+v", properties.GetUpstreamHostMetadata())
}

func printNodeMetadataProperties() {
	proxywasm.LogInfof(">> GetNodeMetadataAnnotations: %v", properties.GetNodeMetadataAnnotations())
	proxywasm.LogInfof(">> GetNodeMetadataAppContainers: %v", properties.GetNodeMetadataAppContainers())
	proxywasm.LogInfof(">> GetNodeMetadataClusterId: %v", properties.GetNodeMetadataClusterId())
	proxywasm.LogInfof(">> GetNodeMetadataEnvoyPrometheusPort: %v", properties.GetNodeMetadataEnvoyPrometheusPort())
	proxywasm.LogInfof(">> GetNodeMetadataEnvoyStatusPort: %v", properties.GetNodeMetadataEnvoyStatusPort())
	proxywasm.LogInfof(">> GetNodeMetadataInstanceIps: %v", properties.GetNodeMetadataInstanceIps())
	proxywasm.LogInfof(">> GetNodeMetadataInterceptionMode: %v", properties.GetNodeMetadataInterceptionMode())
	proxywasm.LogInfof(">> GetNodeMetadataIstioProxySha: %v", properties.GetNodeMetadataIstioProxySha())
	proxywasm.LogInfof(">> GetNodeMetadataIstioVersion: %v", properties.GetNodeMetadataIstioVersion())
	proxywasm.LogInfof(">> GetNodeMetadataLabels: %v", properties.GetNodeMetadataLabels())
	proxywasm.LogInfof(">> GetNodeMetadataMeshId: %v", properties.GetNodeMetadataMeshId())
	proxywasm.LogInfof(">> GetNodeMetadataName: %v", properties.GetNodeMetadataName())
	proxywasm.LogInfof(">> GetNodeMetadataNamespace: %v", properties.GetNodeMetadataNamespace())
	proxywasm.LogInfof(">> GetNodeMetadataNodeName: %v", properties.GetNodeMetadataNodeName())
	proxywasm.LogInfof(">> GetNodeMetadataOwner: %v", properties.GetNodeMetadataOwner())
	proxywasm.LogInfof(">> GetNodeMetadataPilotSan: %v", properties.GetNodeMetadataPilotSan())
	proxywasm.LogInfof(">> GetNodeMetadataPodPorts: %v", properties.GetNodeMetadataPodPorts())
	proxywasm.LogInfof(">> GetNodeMetadataServiceAccount: %v", properties.GetNodeMetadataServiceAccount())
	proxywasm.LogInfof(">> GetNodeMetadataWorkloadName: %v", properties.GetNodeMetadataWorkloadName())
}

func printNodeProxyConfigProperties() {
	proxywasm.LogInfof(">> GetNodeProxyConfigBinaryPath: %v", properties.GetNodeProxyConfigBinaryPath())
	proxywasm.LogInfof(">> GetNodeProxyConfigConcurrency: %v", properties.GetNodeProxyConfigConcurrency())
	proxywasm.LogInfof(">> GetNodeProxyConfigConfigPath: %v", properties.GetNodeProxyConfigConfigPath())
	proxywasm.LogInfof(">> GetNodeProxyConfigControlPlaneAuthPolicy: %v", properties.GetNodeProxyConfigControlPlaneAuthPolicy())
	proxywasm.LogInfof(">> GetNodeProxyConfigDiscoveryAddress: %v", properties.GetNodeProxyConfigDiscoveryAddress())
	proxywasm.LogInfof(">> GetNodeProxyConfigDrainDuration: %v", properties.GetNodeProxyConfigDrainDuration())
	proxywasm.LogInfof(">> GetNodeProxyConfigExtraStatTags: %v", properties.GetNodeProxyConfigExtraStatTags())
	proxywasm.LogInfof(">> GetNodeProxyConfigHoldApplicationUntilProxyStarts: %v", properties.GetNodeProxyConfigHoldApplicationUntilProxyStarts())
	proxywasm.LogInfof(">> GetNodeProxyConfigProxyAdminPort: %v", properties.GetNodeProxyConfigProxyAdminPort())
	proxywasm.LogInfof(">> GetNodeProxyConfigProxyStatsMatcher: %v", properties.GetNodeProxyConfigProxyStatsMatcher())
	proxywasm.LogInfof(">> GetNodeProxyConfigServiceCluster: %v", properties.GetNodeProxyConfigServiceCluster())
	proxywasm.LogInfof(">> GetNodeProxyConfigStatNameLength: %v", properties.GetNodeProxyConfigStatNameLength())
	proxywasm.LogInfof(">> GetNodeProxyConfigStatusPort: %v", properties.GetNodeProxyConfigStatusPort())
	proxywasm.LogInfof(">> GetNodeProxyConfigTerminationDrainDuration: %v", properties.GetNodeProxyConfigTerminationDrainDuration())
	proxywasm.LogInfof(">> GetNodeProxyConfigTracingZipkinAddress: %v", properties.GetNodeProxyConfigTracingZipkinAddress())
}

func printXdsProperties() {
	proxywasm.LogInfof(">> GetXdsClusterName: %v", properties.GetXdsClusterName())
	proxywasm.LogInfof(">> GetXdsClusterMetadata: %+v", properties.GetXdsClusterMetadata())
	proxywasm.LogInfof(">> GetXdsRouteName: %v", properties.GetXdsRouteName())
	proxywasm.LogInfof(">> GetXdsRouteMetadata: %+v", properties.GetXdsRouteMetadata())
	proxywasm.LogInfof(">> GetXdsUpstreamHostMetadata: %+v", properties.GetXdsUpstreamHostMetadata())
	proxywasm.LogInfof(">> GetXdsListenerFilterChainName: %v", properties.GetXdsListenerFilterChainName())
}

func printUpstreamProperties() {
	proxywasm.LogInfof(">> GetUpstreamAddress: %v", properties.GetUpstreamAddress())
	proxywasm.LogInfof(">> GetUpstreamPort: %v", properties.GetUpstreamPort())
	proxywasm.LogInfof(">> GetUpstreamTlsVersion: %v", properties.GetUpstreamTlsVersion())
	proxywasm.LogInfof(">> GetUpstreamSubjectLocalCertificate: %v", properties.GetUpstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> GetUpstreamSubjectPeerCertificate: %v", properties.GetUpstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> GetUpstreamDnsSanLocalCertificate: %v", properties.GetUpstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> GetUpstreamDnsSanPeerCertificate: %v", properties.GetUpstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> GetUpstreamUriSanLocalCertificate: %v", properties.GetUpstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> GetUpstreamUriSanPeerCertificate: %v", properties.GetUpstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> GetUpstreamSha256PeerCertificateDigest: %v", properties.GetUpstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> GetUpstreamLocalAddress: %v", properties.GetUpstreamLocalAddress())
	proxywasm.LogInfof(">> GetUpstreamTransportFailureReason: %v", properties.GetUpstreamTransportFailureReason())
}

func printConnectionProperties() {
	proxywasm.LogInfof(">> GetDownstreamRemoteAddress: %v", properties.GetDownstreamRemoteAddress())
	proxywasm.LogInfof(">> GetDownstreamRemotePort: %v", properties.GetDownstreamRemotePort())
	proxywasm.LogInfof(">> GetDownstreamLocalAddress: %v", properties.GetDownstreamLocalAddress())
	proxywasm.LogInfof(">> GetDownstreamLocalPort: %v", properties.GetDownstreamLocalPort())
	proxywasm.LogInfof(">> GetDownstreamConnectionId: %v", properties.GetDownstreamConnectionId())
	proxywasm.LogInfof(">> IsDownstreamConnectionTls: %v", properties.IsDownstreamConnectionTls())
	proxywasm.LogInfof(">> GetDownstreamRequestedServerName: %v", properties.GetDownstreamRequestedServerName())
	proxywasm.LogInfof(">> GetDownstreamTlsVersion: %v", properties.GetDownstreamTlsVersion())
	proxywasm.LogInfof(">> GetDownstreamSubjectLocalCertificate: %v", properties.GetDownstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> GetDownstreamSubjectPeerCertificate: %v", properties.GetDownstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> GetDownstreamDnsSanLocalCertificate: %v", properties.GetDownstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> GetDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> GetDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> GetDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> GetDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> GetDownstreamTerminationDetails: %v", properties.GetDownstreamTerminationDetails())
}

func printResponseProperties() {
	proxywasm.LogInfof(">> GetResponseCode: %v", properties.GetResponseCode())
	proxywasm.LogInfof(">> GetResponseCodeDetails: %v", properties.GetResponseCodeDetails())
	proxywasm.LogInfof(">> GetResponseFlags: %v", properties.GetResponseFlags())
	proxywasm.LogInfof(">> GetResponseGrpcStatusCode: %v", properties.GetResponseGrpcStatusCode())
	proxywasm.LogInfof(">> GetResponseHeaders: %+v", properties.GetResponseHeaders())
	proxywasm.LogInfof(">> GetResponseTrailers: %+v", properties.GetResponseTrailers())
	proxywasm.LogInfof(">> GetResponseSize: %v", properties.GetResponseSize())
	proxywasm.LogInfof(">> GetResponseTotalSize: %v", properties.GetResponseTotalSize())
}

func printRequestProperties() {
	proxywasm.LogInfof(">> GetRequestPath: %v", properties.GetRequestPath())
	proxywasm.LogInfof(">> GetRequestUrlPath: %v", properties.GetRequestUrlPath())
	proxywasm.LogInfof(">> GetRequestHost: %v", properties.GetRequestHost())
	proxywasm.LogInfof(">> GetRequestScheme: %v", properties.GetRequestScheme())
	proxywasm.LogInfof(">> GetRequestMethod: %v", properties.GetRequestMethod())
	proxywasm.LogInfof(">> GetRequestHeaders: %+v", properties.GetRequestHeaders())
	proxywasm.LogInfof(">> GetRequestReferer: %v", properties.GetRequestReferer())
	proxywasm.LogInfof(">> GetRequestUserAgent: %v", properties.GetRequestUserAgent())
	proxywasm.LogInfof(">> GetRequestTime: %v", properties.GetRequestTime())
	proxywasm.LogInfof(">> GetRequestId: %v", properties.GetRequestId())
	proxywasm.LogInfof(">> GetRequestProtocol: %v", properties.GetRequestProtocol())
	proxywasm.LogInfof(">> GetRequestQuery: %v", properties.GetRequestQuery())
	proxywasm.LogInfof(">> GetRequestDuration: %v", properties.GetRequestDuration())
	proxywasm.LogInfof(">> GetRequestSize: %v", properties.GetRequestSize())
	proxywasm.LogInfof(">> GetRequestTotalSize: %v", properties.GetRequestTotalSize())
}
