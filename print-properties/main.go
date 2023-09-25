package main

import (
	"wasm-minimal/properties"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
)

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("********** OnPluginStart **********")
	// printWasmProperties()

	return types.OnPluginStartStatusOK
}

// Override types.DefaultPluginContext.
func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("********** NewHttpContext **********")
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	// Embed the default http context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

// *********************************************
// REQUEST PATH
// *********************************************

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestHeaders **********")
	printWasmProperties()
	printConfigurationProperties()
	printUpstreamProperties()
	printConnectionProperties()
	printRequestProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestBody **********")
	printRequestProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestTrailers(bodySize int) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestTrailers **********")
	printRequestProperties()

	return types.ActionContinue
}

// *********************************************
// RESPONSE PATH
// *********************************************

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
	printResponseProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseBody **********")
	printResponseProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseTrailers(bodySize int) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseTrailers **********")
	printResponseProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogInfo("********** OnHttpStreamDone **********")
	printRequestProperties()
	printResponseProperties()
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

	proxywasm.LogInfof(">> GetNodeMetadata.annotations: %v", properties.GetNodeMetadataAnnotations())
	proxywasm.LogInfof(">> GetNodeMetadata.appContainers: %v", properties.GetNodeMetadataAppContainers())
	proxywasm.LogInfof(">> GetNodeMetadata.clusterId: %v", properties.GetNodeMetadataClusterId())
	proxywasm.LogInfof(">> GetNodeMetadata.envoyPrometheusPort: %v", properties.GetNodeMetadataEnvoyPrometheusPort())
	proxywasm.LogInfof(">> GetNodeMetadata.envoyStatusPort: %v", properties.GetNodeMetadataEnvoyStatusPort())
	proxywasm.LogInfof(">> GetNodeMetadata.instanceIps: %v", properties.GetNodeMetadataInstanceIps())
	proxywasm.LogInfof(">> GetNodeMetadata.interceptionMode: %v", properties.GetNodeMetadataInterceptionMode())
	proxywasm.LogInfof(">> GetNodeMetadata.istioProxySha: %v", properties.GetNodeMetadataIstioProxySha())
	proxywasm.LogInfof(">> GetNodeMetadata.istioVersion: %v", properties.GetNodeMetadataIstioVersion())
	proxywasm.LogInfof(">> GetNodeMetadata.labels: %v", properties.GetNodeMetadataLabels())
	proxywasm.LogInfof(">> GetNodeMetadata.meshId: %v", properties.GetNodeMetadataMeshId())
	proxywasm.LogInfof(">> GetNodeMetadata.name: %v", properties.GetNodeMetadataName())
	proxywasm.LogInfof(">> GetNodeMetadata.namespace: %v", properties.GetNodeMetadataNamespace())
	proxywasm.LogInfof(">> GetNodeMetadata.nodeName: %v", properties.GetNodeMetadataNodeName())
	proxywasm.LogInfof(">> GetNodeMetadata.owner: %v", properties.GetNodeMetadataOwner())
	proxywasm.LogInfof(">> GetNodeMetadata.pilotSan: %v", properties.GetNodeMetadataPilotSan())
	proxywasm.LogInfof(">> GetNodeMetadata.podPorts: %v", properties.GetNodeMetadataPodPorts())
	proxywasm.LogInfof(">> GetNodeMetadata.serviceAccount: %v", properties.GetNodeMetadataServiceAccount())
	proxywasm.LogInfof(">> GetNodeMetadata.workloadName: %v", properties.GetNodeMetadataWorkloadName())

	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.binaryPath: %v", properties.GetNodeProxyConfigBinaryPath())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.concurrency: %v", properties.GetNodeProxyConfigConcurrency())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.configPath: %v", properties.GetNodeProxyConfigConfigPath())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.controlPlaneAuthPolicy: %v", properties.GetNodeProxyConfigControlPlaneAuthPolicy())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.discoveryAddress: %v", properties.GetNodeProxyConfigDiscoveryAddress())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.drainDuration: %v", properties.GetNodeProxyConfigDrainDuration())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.extraStatTags: %v", properties.GetNodeProxyConfigExtraStatTags())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.holdApplicationUntilProxyStarts: %v", properties.GetNodeProxyConfigHoldApplicationUntilProxyStarts())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.proxyAdminPort: %v", properties.GetNodeProxyConfigProxyAdminPort())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.proxyStatsMatcher: %+v", properties.GetNodeProxyConfigProxyStatsMatcher())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.serviceCluster: %v", properties.GetNodeProxyConfigServiceCluster())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.statNameLength: %v", properties.GetNodeProxyConfigStatNameLength())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.statusPort: %v", properties.GetNodeProxyConfigStatusPort())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.terminationDrainDuration: %v", properties.GetNodeProxyConfigTerminationDrainDuration())
	proxywasm.LogInfof(">> GetNodeMetadata.proxyConfig.tracing.zipkin.address: %v", properties.GetNodeProxyConfigTracingZipkinAddress())

	proxywasm.LogInfof(">> GetNodeDynamicParams: %v", properties.GetNodeDynamicParams())
	proxywasm.LogInfof(">> GetNodeLocality: %+v", properties.GetNodeLocality())
	proxywasm.LogInfof(">> GetNodeUserAgentName: %v", properties.GetNodeUserAgentName())
	proxywasm.LogInfof(">> GetNodeUserAgentVersion: %v", properties.GetNodeUserAgentVersion())
	proxywasm.LogInfof(">> GetNodeUserAgentBuildVersion: %v", properties.GetNodeUserAgentBuildVersion())
	proxywasm.LogInfof(">> GetNodeExtensions: %v", properties.GetNodeExtensions())
	proxywasm.LogInfof(">> GetNodeClientFeatures: %v", properties.GetNodeClientFeatures())
	proxywasm.LogInfof(">> GetNodeListeningAddresses: %v", properties.GetNodeListeningAddresses())
	proxywasm.LogInfof(">> GetClusterMetadata: %+v", properties.GetClusterMetadata())
	proxywasm.LogInfof(">> GetListenerMetadata: %+v", properties.GetListenerMetadata())
	proxywasm.LogInfof(">> GetRouteMetadata: %+v", properties.GetRouteMetadata())
	proxywasm.LogInfof(">> GetUpstreamHostMetadata: %+v", properties.GetUpstreamHostMetadata())
}

func printConfigurationProperties() {
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
