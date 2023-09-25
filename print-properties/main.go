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
	proxywasm.LogInfof(">> getPluginName: %v", properties.GetPluginName())
	proxywasm.LogInfof(">> getPluginRootId: %v", properties.GetPluginRootId())
	proxywasm.LogInfof(">> getPluginVmId: %v", properties.GetPluginVmId())
	proxywasm.LogInfof(">> getClusterName: %v", properties.GetClusterName())
	proxywasm.LogInfof(">> getRouteName: %v", properties.GetRouteName())
	proxywasm.LogInfof(">> getListenerDirection: %v", properties.GetListenerDirection())
	proxywasm.LogInfof(">> getNodeId: %v", properties.GetNodeId())
	proxywasm.LogInfof(">> getNodeCluster: %v", properties.GetNodeCluster())

	proxywasm.LogInfof(">> getNodeMetadata.annotations: %v", properties.GetNodeMetadataAnnotations())
	proxywasm.LogInfof(">> getNodeMetadata.appContainers: %v", properties.GetNodeMetadataAppContainers())
	proxywasm.LogInfof(">> getNodeMetadata.clusterId: %v", properties.GetNodeMetadataClusterId())
	proxywasm.LogInfof(">> getNodeMetadata.envoyPrometheusPort: %v", properties.GetNodeMetadataEnvoyPrometheusPort())
	proxywasm.LogInfof(">> getNodeMetadata.envoyStatusPort: %v", properties.GetNodeMetadataEnvoyStatusPort())
	proxywasm.LogInfof(">> getNodeMetadata.instanceIps: %v", properties.GetNodeMetadataInstanceIps())
	proxywasm.LogInfof(">> getNodeMetadata.interceptionMode: %v", properties.GetNodeMetadataInterceptionMode())
	proxywasm.LogInfof(">> getNodeMetadata.istioProxySha: %v", properties.GetNodeMetadataIstioProxySha())
	proxywasm.LogInfof(">> getNodeMetadata.istioVersion: %v", properties.GetNodeMetadataIstioVersion())
	proxywasm.LogInfof(">> getNodeMetadata.labels: %v", properties.GetNodeMetadataLabels())
	proxywasm.LogInfof(">> getNodeMetadata.meshId: %v", properties.GetNodeMetadataMeshId())
	proxywasm.LogInfof(">> getNodeMetadata.name: %v", properties.GetNodeMetadataName())
	proxywasm.LogInfof(">> getNodeMetadata.namespace: %v", properties.GetNodeMetadataNamespace())
	proxywasm.LogInfof(">> getNodeMetadata.nodeName: %v", properties.GetNodeMetadataNodeName())
	proxywasm.LogInfof(">> getNodeMetadata.owner: %v", properties.GetNodeMetadataOwner())
	proxywasm.LogInfof(">> getNodeMetadata.pilotSan: %v", properties.GetNodeMetadataPilotSan())
	proxywasm.LogInfof(">> getNodeMetadata.podPorts: %v", properties.GetNodeMetadataPodPorts())
	proxywasm.LogInfof(">> getNodeMetadata.serviceAccount: %v", properties.GetNodeMetadataServiceAccount())
	proxywasm.LogInfof(">> getNodeMetadata.workloadName: %v", properties.GetNodeMetadataWorkloadName())

	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.binaryPath: %v", properties.GetNodeProxyConfigBinaryPath())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.concurrency: %v", properties.GetNodeProxyConfigConcurrency())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.configPath: %v", properties.GetNodeProxyConfigConfigPath())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.controlPlaneAuthPolicy: %v", properties.GetNodeProxyConfigControlPlaneAuthPolicy())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.discoveryAddress: %v", properties.GetNodeProxyConfigDiscoveryAddress())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.drainDuration: %v", properties.GetNodeProxyConfigDrainDuration())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.extraStatTags: %v", properties.GetNodeProxyConfigExtraStatTags())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.proxyAdminPort: %v", properties.GetNodeProxyConfigProxyAdminPort())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.serviceCluster: %v", properties.GetNodeProxyConfigServiceCluster())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.statNameLength: %v", properties.GetNodeProxyConfigStatNameLength())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.statusPort: %v", properties.GetNodeProxyConfigStatusPort())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.terminationDrainDuration: %v", properties.GetNodeProxyConfigTerminationDrainDuration())
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.tracing.zipkin.address: %v", properties.GetNodeProxyConfigTracingZipkinAddress())

	proxywasm.LogInfof(">> getNodeDynamicParams: %v", properties.GetNodeDynamicParams())
	proxywasm.LogInfof(">> getNodeLocality: %+v", properties.GetNodeLocality())
	proxywasm.LogInfof(">> getNodeUserAgentName: %v", properties.GetNodeUserAgentName())
	proxywasm.LogInfof(">> getNodeUserAgentVersion: %v", properties.GetNodeUserAgentVersion())
	proxywasm.LogInfof(">> getNodeUserAgentBuildVersion: %v", properties.GetNodeUserAgentBuildVersion())
	proxywasm.LogInfof(">> getNodeExtensions: %v", properties.GetNodeExtensions())
	proxywasm.LogInfof(">> getNodeClientFeatures: %v", properties.GetNodeClientFeatures())
	proxywasm.LogInfof(">> getNodeListeningAddresses: %v", properties.GetNodeListeningAddresses())
	proxywasm.LogInfof(">> getClusterMetadata: %+v", properties.GetClusterMetadata())
	proxywasm.LogInfof(">> getListenerMetadata: %+v", properties.GetListenerMetadata())
	proxywasm.LogInfof(">> getRouteMetadata: %+v", properties.GetRouteMetadata())
	proxywasm.LogInfof(">> getUpstreamHostMetadata: %+v", properties.GetUpstreamHostMetadata())
}

func printConfigurationProperties() {
	proxywasm.LogInfof(">> getXdsClusterName: %v", properties.GetXdsClusterName())
	proxywasm.LogInfof(">> getXdsClusterMetadata: %+v", properties.GetXdsClusterMetadata())
	proxywasm.LogInfof(">> getXdsRouteName: %v", properties.GetXdsRouteName())
	proxywasm.LogInfof(">> getXdsRouteMetadata: %+v", properties.GetXdsRouteMetadata())
	proxywasm.LogInfof(">> getXdsUpstreamHostMetadata: %+v", properties.GetXdsUpstreamHostMetadata())
	proxywasm.LogInfof(">> getXdsListenerFilterChainName: %v", properties.GetXdsListenerFilterChainName())
}

func printUpstreamProperties() {
	proxywasm.LogInfof(">> getUpstreamAddress: %v", properties.GetUpstreamAddress())
	proxywasm.LogInfof(">> getUpstreamPort: %v", properties.GetUpstreamPort())
	proxywasm.LogInfof(">> getUpstreamTlsVersion: %v", properties.GetUpstreamTlsVersion())
	proxywasm.LogInfof(">> getUpstreamSubjectLocalCertificate: %v", properties.GetUpstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamSubjectPeerCertificate: %v", properties.GetUpstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamDnsSanLocalCertificate: %v", properties.GetUpstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamDnsSanPeerCertificate: %v", properties.GetUpstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamUriSanLocalCertificate: %v", properties.GetUpstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamUriSanPeerCertificate: %v", properties.GetUpstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamSha256PeerCertificateDigest: %v", properties.GetUpstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> getUpstreamLocalAddress: %v", properties.GetUpstreamLocalAddress())
	proxywasm.LogInfof(">> getUpstreamTransportFailureReason: %v", properties.GetUpstreamTransportFailureReason())
}

func printConnectionProperties() {
	proxywasm.LogInfof(">> getDownstreamRemoteAddress: %v", properties.GetDownstreamRemoteAddress())
	proxywasm.LogInfof(">> getDownstreamRemotePort: %v", properties.GetDownstreamRemotePort())
	proxywasm.LogInfof(">> getDownstreamLocalAddress: %v", properties.GetDownstreamLocalAddress())
	proxywasm.LogInfof(">> getDownstreamLocalPort: %v", properties.GetDownstreamLocalPort())
	proxywasm.LogInfof(">> getDownstreamConnectionId: %v", properties.GetDownstreamConnectionId())
	proxywasm.LogInfof(">> isDownstreamConnectionTls: %v", properties.IsDownstreamConnectionTls())
	proxywasm.LogInfof(">> getDownstreamRequestedServerName: %v", properties.GetDownstreamRequestedServerName())
	proxywasm.LogInfof(">> getDownstreamTlsVersion: %v", properties.GetDownstreamTlsVersion())
	proxywasm.LogInfof(">> getDownstreamSubjectLocalCertificate: %v", properties.GetDownstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamSubjectPeerCertificate: %v", properties.GetDownstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanLocalCertificate: %v", properties.GetDownstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", properties.GetDownstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> getDownstreamTerminationDetails: %v", properties.GetDownstreamTerminationDetails())
}

func printResponseProperties() {
	proxywasm.LogInfof(">> getResponseCode: %v", properties.GetResponseCode())
	proxywasm.LogInfof(">> getResponseCodeDetails: %v", properties.GetResponseCodeDetails())
	proxywasm.LogInfof(">> getResponseFlags: %v", properties.GetResponseFlags())
	proxywasm.LogInfof(">> getResponseGrpcStatusCode: %v", properties.GetResponseGrpcStatusCode())
	proxywasm.LogInfof(">> getResponseHeaders: %+v", properties.GetResponseHeaders())
	proxywasm.LogInfof(">> getResponseTrailers: %+v", properties.GetResponseTrailers())
	proxywasm.LogInfof(">> getResponseSize: %v", properties.GetResponseSize())
	proxywasm.LogInfof(">> getResponseTotalSize: %v", properties.GetResponseTotalSize())
}

func printRequestProperties() {
	proxywasm.LogInfof(">> getRequestPath: %v", properties.GetRequestPath())
	proxywasm.LogInfof(">> getRequestUrlPath: %v", properties.GetRequestUrlPath())
	proxywasm.LogInfof(">> getRequestHost: %v", properties.GetRequestHost())
	proxywasm.LogInfof(">> getRequestScheme: %v", properties.GetRequestScheme())
	proxywasm.LogInfof(">> getRequestMethod: %v", properties.GetRequestMethod())
	proxywasm.LogInfof(">> getRequestHeaders: %+v", properties.GetRequestHeaders())
	proxywasm.LogInfof(">> getRequestReferer: %v", properties.GetRequestReferer())
	proxywasm.LogInfof(">> getRequestUserAgent: %v", properties.GetRequestUserAgent())
	proxywasm.LogInfof(">> getRequestTime: %v", properties.GetRequestTime())
	proxywasm.LogInfof(">> getRequestId: %v", properties.GetRequestId())
	proxywasm.LogInfof(">> getRequestProtocol: %v", properties.GetRequestProtocol())
	proxywasm.LogInfof(">> getRequestQuery: %v", properties.GetRequestQuery())
	proxywasm.LogInfof(">> getRequestDuration: %v", properties.GetRequestDuration())
	proxywasm.LogInfof(">> getRequestSize: %v", properties.GetRequestSize())
	proxywasm.LogInfof(">> getRequestTotalSize: %v", properties.GetRequestTotalSize())
}
