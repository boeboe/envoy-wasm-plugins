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

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
	printWasmProperties()

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
	proxywasm.LogInfof(">> getNodeMetadata.meshId: %v", getNodeMetadata().meshId)
	proxywasm.LogInfof(">> getNodeMetadata.istioVersion: %v", getNodeMetadata().istioVersion)
	proxywasm.LogInfof(">> getNodeMetadata.envoyPrometheusPort: %v", getNodeMetadata().envoyPrometheusPort)
	proxywasm.LogInfof(">> getNodeMetadata.envoyStatusPort: %v", getNodeMetadata().envoyStatusPort)
	proxywasm.LogInfof(">> getNodeMetadata.annotations: %v", getNodeMetadata().annotations)
	proxywasm.LogInfof(">> getNodeMetadata: %+v", getNodeMetadata())
	proxywasm.LogInfof(">> getNodeDynamicParams: %v", getNodeDynamicParams())
	proxywasm.LogInfof(">> getNodeLocality: %v", getNodeLocality())
	proxywasm.LogInfof(">> getNodeUserAgentName: %v", getNodeUserAgentName())
	proxywasm.LogInfof(">> getNodeUserAgentVersion: %v", getNodeUserAgentVersion())
	proxywasm.LogInfof(">> getNodeUserAgentBuildVersion: %v", getNodeUserAgentBuildVersion())
	proxywasm.LogInfof(">> getNodeExtensions: %v", getNodeExtensions())
	proxywasm.LogInfof(">> getNodeClientFeatures: %v", getNodeClientFeatures())
	proxywasm.LogInfof(">> getNodeListeningAddresses: %v", getNodeListeningAddresses())
	proxywasm.LogInfof(">> getClusterMetadata: %v", getClusterMetadata())
	proxywasm.LogInfof(">> getListenerMetadata: %v", getListenerMetadata())
	proxywasm.LogInfof(">> getRouteMetadata: %v", getRouteMetadata())
	proxywasm.LogInfof(">> getUpstreamHostMetadata: %v", getUpstreamHostMetadata())
}
