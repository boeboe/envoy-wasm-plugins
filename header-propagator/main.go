package main

import (
	"header-propagator/properties"
	"header-propagator/utils"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

// pluginConfig represents the configuration for the plugin.
type pluginConfig struct {
	CorrelationHeader string            `json:"correlationHeader"`
	PropagationHeader propagationHeader `json:"propagationHeader"`
	Enabled           bool              `json:"enabled"`
}

// propagationHeader represents the structure of the propagation header.
type propagationHeader struct {
	Name    string `json:"name"`
	Default string `json:"default"`
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

type vmContext struct {
	types.DefaultVMContext
}

// NewPluginContext creates a new plugin context.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	config pluginConfig
}

// OnPluginStart initializes the plugin.
func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("Starting plugin...")

	if pluginConfigurationSize == 0 {
		return types.OnPluginStartStatusOK
	}

	// Read the plugin configuration.
	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogErrorf("Failed to get plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}

	// Validate the plugin configuration JSON.
	if !gjson.ValidBytes(configData) {
		proxywasm.LogErrorf("Invalid JSON in plugin configuration: %v", string(configData))
		return types.OnPluginStartStatusFailed
	}
	jsonData := gjson.ParseBytes(configData)

	// Parse the configuration.
	p.config.Enabled = jsonData.Get("enabled").Bool()
	p.config.CorrelationHeader = jsonData.Get("correlationHeader").String()
	p.config.PropagationHeader = propagationHeader{
		Name:    jsonData.Get("propagationHeader.name").String(),
		Default: jsonData.Get("propagationHeader.default").String(),
	}

	proxywasm.LogDebugf("Parsed plugin config: %v", p.config)
	return types.OnPluginStartStatusOK
}

// NewHttpContext creates a new HTTP context.
func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("Creating new HTTP context...")
	return &httpContext{
		contextID:         contextID,
		correlationHeader: p.config.CorrelationHeader,
		propagationHeader: p.config.PropagationHeader,
	}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID         uint32
	correlationHeader string
	propagationHeader propagationHeader
}

// OnHttpRequestHeaders handles incoming request headers.
func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("Processing request headers...")

	reqHeaders := properties.GetRequestHeaders()
	direction := properties.GetListenerDirection()

	cHeaderVal, ok := reqHeaders[strings.ToLower(ctx.correlationHeader)]
	if !ok {
		proxywasm.LogInfof("Correlation header '%v' not available in the request headers", strings.ToLower(ctx.correlationHeader))
		return types.ActionContinue
	}

	switch direction.String() {
	case "INBOUND":
		handleInboundRequest(ctx, reqHeaders, cHeaderVal)
	case "OUTBOUND":
		handleOutboundRequest(ctx, reqHeaders, cHeaderVal)
	}

	return types.ActionContinue
}

// handleInboundRequest processes an inbound request.
func handleInboundRequest(ctx *httpContext, reqHeaders map[string]string, cHeaderVal string) {
	pHeaderVal, ok := reqHeaders[strings.ToLower(ctx.propagationHeader.Name)]
	if !ok {
		if err := proxywasm.AddHttpRequestHeader(strings.ToLower(ctx.propagationHeader.Name), ctx.propagationHeader.Default); err != nil {
			proxywasm.LogErrorf("Failed to set inbound request header: %v", err)
		}
	} else {
		if err := utils.SetSharedDataSafe(cHeaderVal, []byte(pHeaderVal), 0); err != nil {
			proxywasm.LogErrorf("Failed to set shared data key '%v' with value '%v': %v", cHeaderVal, pHeaderVal, err)
		}
	}
}

// handleOutboundRequest processes an outbound request.
func handleOutboundRequest(ctx *httpContext, reqHeaders map[string]string, cHeaderVal string) {
	_, ok := reqHeaders[strings.ToLower(ctx.propagationHeader.Name)]
	if !ok {
		pHeaderVal, _, err := utils.GetSharedDataSafe(cHeaderVal)
		if err != nil {
			proxywasm.LogErrorf("Failed to get shared data key '%v': %v", cHeaderVal, err)
		}
		if err := proxywasm.AddHttpRequestHeader(strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal)); err != nil {
			proxywasm.LogErrorf("Failed to set outbound request header '%v' with value '%v': %v", strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal), err)
		}
	}
}

// OnHttpResponseHeaders handles incoming response headers.
func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("Processing response headers...")

	resHeaders := properties.GetResponseHeaders()
	direction := properties.GetListenerDirection()

	cHeaderVal, ok := resHeaders[strings.ToLower(ctx.correlationHeader)]
	if !ok {
		proxywasm.LogInfof("Correlation header '%v' not available in the response headers", strings.ToLower(ctx.correlationHeader))
		return types.ActionContinue
	}

	switch direction.String() {
	case "INBOUND":
		handleInboundResponse(ctx, resHeaders, cHeaderVal)
	case "OUTBOUND":
		handleOutboundResponse(ctx, resHeaders, cHeaderVal)
	}

	return types.ActionContinue
}

// handleInboundResponse processes an inbound response.
func handleInboundResponse(ctx *httpContext, resHeaders map[string]string, cHeaderVal string) {
	pHeaderVal, _, err := utils.GetSharedDataSafe(cHeaderVal)
	if err != nil {
		proxywasm.LogErrorf("Failed to get shared data key '%v': %v", cHeaderVal, err)
	}
	if err := proxywasm.AddHttpResponseHeader(strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal)); err != nil {
		proxywasm.LogErrorf("Failed to set inbound response header '%v' with value '%v': %v", strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal), err)
	}
}

// handleOutboundResponse processes an outbound response.
func handleOutboundResponse(ctx *httpContext, resHeaders map[string]string, cHeaderVal string) {
	pHeaderVal, _, err := utils.GetSharedDataSafe(cHeaderVal)
	if err != nil {
		proxywasm.LogErrorf("Failed to get shared data key '%v': %v", cHeaderVal, err)
	}
	if err := proxywasm.AddHttpResponseHeader(strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal)); err != nil {
		proxywasm.LogErrorf("Failed to set outbound response header '%v' with value '%v': %v", strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal), err)
	}
}
