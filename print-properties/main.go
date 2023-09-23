package main

import (
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

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestHeaders **********")
	printWasmProperties()
	printConfigurationProperties()
	printUpstreamProperties()
	printConnectionProperties()

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
	printWasmProperties()
	printConfigurationProperties()
	printUpstreamProperties()
	printConnectionProperties()

	return types.ActionContinue
}

func printWasmProperties() {
	proxywasm.LogInfof(">> getPluginName: %v", getPluginName())
	proxywasm.LogInfof(">> getPluginRootId: %v", getPluginRootId())
	proxywasm.LogInfof(">> getPluginVmId: %v", getPluginVmId())
	proxywasm.LogInfof(">> getClusterName: %v", getClusterName())
	proxywasm.LogInfof(">> getRouteName: %v", getRouteName())
	proxywasm.LogInfof(">> getListenerDirection: %v", getListenerDirection())
	proxywasm.LogInfof(">> getNodeId: %v", getNodeId())
	proxywasm.LogInfof(">> getNodeCluster: %v", getNodeCluster())

	proxywasm.LogInfof(">> getNodeMetadata.annotations: %v", getNodeMetadata().annotations)
	proxywasm.LogInfof(">> getNodeMetadata.appContainers: %v", getNodeMetadata().annotations)
	proxywasm.LogInfof(">> getNodeMetadata.clusterId: %v", getNodeMetadata().clusterId)
	proxywasm.LogInfof(">> getNodeMetadata.envoyPrometheusPort: %v", getNodeMetadata().envoyPrometheusPort)
	proxywasm.LogInfof(">> getNodeMetadata.envoyStatusPort: %v", getNodeMetadata().envoyStatusPort)
	proxywasm.LogInfof(">> getNodeMetadata.instanceIps: %v", getNodeMetadata().instanceIps)
	proxywasm.LogInfof(">> getNodeMetadata.interceptionMode: %v", getNodeMetadata().instanceIps)
	proxywasm.LogInfof(">> getNodeMetadata.istioProxySha: %v", getNodeMetadata().istioProxySha)
	proxywasm.LogInfof(">> getNodeMetadata.istioVersion: %v", getNodeMetadata().istioVersion)
	proxywasm.LogInfof(">> getNodeMetadata.labels: %v", getNodeMetadata().labels)
	proxywasm.LogInfof(">> getNodeMetadata.meshId: %v", getNodeMetadata().meshId)
	proxywasm.LogInfof(">> getNodeMetadata.name: %v", getNodeMetadata().name)
	proxywasm.LogInfof(">> getNodeMetadata.namespace: %v", getNodeMetadata().namespace)
	proxywasm.LogInfof(">> getNodeMetadata.nodeName: %v", getNodeMetadata().nodeName)
	proxywasm.LogInfof(">> getNodeMetadata.owner: %v", getNodeMetadata().owner)
	proxywasm.LogInfof(">> getNodeMetadata.pilotSan: %v", getNodeMetadata().pilotSan)
	proxywasm.LogInfof(">> getNodeMetadata.podPorts: %v", getNodeMetadata().podPorts)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.binaryPath: %v", getNodeMetadata().proxyConfig.binaryPath)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.concurrency: %v", getNodeMetadata().proxyConfig.concurrency)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.configPath: %v", getNodeMetadata().proxyConfig.configPath)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.controlPlaneAuthPolicy: %v", getNodeMetadata().proxyConfig.controlPlaneAuthPolicy)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.discoveryAddress: %v", getNodeMetadata().proxyConfig.discoveryAddress)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.drainDuration: %v", getNodeMetadata().proxyConfig.drainDuration)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.extraStatTags: %v", getNodeMetadata().proxyConfig.extraStatTags)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.proxyAdminPort: %v", getNodeMetadata().proxyConfig.proxyAdminPort)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.serviceCluster: %v", getNodeMetadata().proxyConfig.serviceCluster)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.statNameLength: %v", getNodeMetadata().proxyConfig.statNameLength)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.statusPort: %v", getNodeMetadata().proxyConfig.statusPort)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.terminationDrainDuration: %v", getNodeMetadata().proxyConfig.terminationDrainDuration)
	proxywasm.LogInfof(">> getNodeMetadata.proxyConfig.tracing.address: %v", getNodeMetadata().proxyConfig.tracing.address)
	proxywasm.LogInfof(">> getNodeMetadata.serviceAccount: %v", getNodeMetadata().serviceAccount)
	proxywasm.LogInfof(">> getNodeMetadata.workloadName: %v", getNodeMetadata().workloadName)

	proxywasm.LogInfof(">> getNodeDynamicParams: %v", getNodeDynamicParams())
	proxywasm.LogInfof(">> getNodeLocality: %v", getNodeLocality())
	proxywasm.LogInfof(">> getNodeUserAgentName: %v", getNodeUserAgentName())
	proxywasm.LogInfof(">> getNodeUserAgentVersion: %v", getNodeUserAgentVersion())
	proxywasm.LogInfof(">> getNodeUserAgentBuildVersion: %v", getNodeUserAgentBuildVersion())
	proxywasm.LogInfof(">> getNodeExtensions: %v", getNodeExtensions())
	proxywasm.LogInfof(">> getNodeClientFeatures: %v", getNodeClientFeatures())
	proxywasm.LogInfof(">> getNodeListeningAddresses: %v", getNodeListeningAddresses())
	proxywasm.LogInfof(">> getClusterMetadata: %+v", getClusterMetadata())
	proxywasm.LogInfof(">> getListenerMetadata: %+v", getListenerMetadata())
	proxywasm.LogInfof(">> getRouteMetadata: %+v", getRouteMetadata())
	proxywasm.LogInfof(">> getUpstreamHostMetadata: %+v", getUpstreamHostMetadata())
}

func printConfigurationProperties() {
	proxywasm.LogInfof(">> getXdsClusterName: %v", getXdsClusterName())
	proxywasm.LogInfof(">> getXdsClusterMetadata: %+v", getXdsClusterMetadata())
	proxywasm.LogInfof(">> getXdsRouteName: %v", getXdsRouteName())
	proxywasm.LogInfof(">> getXdsRouteMetadata: %+v", getXdsRouteMetadata())
	proxywasm.LogInfof(">> getXdsUpstreamHostMetadata: %+v", getXdsUpstreamHostMetadata())
	proxywasm.LogInfof(">> getXdsListenerFilterChainName: %v", getXdsListenerFilterChainName())
}

func printUpstreamProperties() {
	proxywasm.LogInfof(">> getUpstreamAddress: %v", getUpstreamAddress())
	proxywasm.LogInfof(">> getUpstreamPort: %v", getUpstreamPort())
	proxywasm.LogInfof(">> getUpstreamTlsVersion: %v", getUpstreamTlsVersion())
	proxywasm.LogInfof(">> getUpstreamSubjectLocalCertificate: %v", getUpstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamSubjectPeerCertificate: %v", getUpstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamDnsSanLocalCertificate: %v", getUpstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamDnsSanPeerCertificate: %v", getUpstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamUriSanLocalCertificate: %v", getUpstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> getUpstreamUriSanPeerCertificate: %v", getUpstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> getUpstreamSha256PeerCertificateDigest: %v", getUpstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> getUpstreamLocalAddress: %v", getUpstreamLocalAddress())
	proxywasm.LogInfof(">> getUpstreamTransportFailureReason: %v", getUpstreamTransportFailureReason())
}

func printConnectionProperties() {
	proxywasm.LogInfof(">> getDownstreamRemoteAddress: %v", getDownstreamRemoteAddress())
	proxywasm.LogInfof(">> getDownstreamRemotePort: %v", getDownstreamRemotePort())
	proxywasm.LogInfof(">> getDownstreamLocalAddress: %v", getDownstreamLocalAddress())
	proxywasm.LogInfof(">> getDownstreamLocalPort: %v", getDownstreamLocalPort())
	proxywasm.LogInfof(">> getDownstreamConnectionId: %v", getDownstreamConnectionId())
	proxywasm.LogInfof(">> isDownstreamConnectionTls: %v", isDownstreamConnectionTls())
	proxywasm.LogInfof(">> getDownstreamRequestedServerName: %v", getDownstreamRequestedServerName())
	proxywasm.LogInfof(">> getDownstreamTlsVersion: %v", getDownstreamTlsVersion())
	proxywasm.LogInfof(">> getDownstreamSubjectLocalCertificate: %v", getDownstreamSubjectLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamSubjectPeerCertificate: %v", getDownstreamSubjectPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanLocalCertificate: %v", getDownstreamDnsSanLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", getDownstreamDnsSanPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", getDownstreamUriSanLocalCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", getDownstreamUriSanPeerCertificate())
	proxywasm.LogInfof(">> getDownstreamDnsSanPeerCertificate: %v", getDownstreamSha256PeerCertificateDigest())
	proxywasm.LogInfof(">> getDownstreamTerminationDetails: %v", getDownstreamTerminationDetails())
}
