package main

import (
	"header-propagator/properties"
	"header-propagator/utils"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/tidwall/gjson"
)

type pluginConfig struct {
	CorrelationHeader string            `json:"correlationHeader"`
	PropagationHeader propagationHeader `json:"propagationHeader"`
	Enabled           bool              `json:"enabled"`
}

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

func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	types.DefaultPluginContext
	config pluginConfig
}

func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("********** OnPluginStart **********")
	if pluginConfigurationSize == 0 {
		return types.OnPluginStartStatusOK
	}

	// Read the plugin configuration.
	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil {
		proxywasm.LogErrorf("failed to get plugin configuration: %v", err.Error())
		return types.OnPluginStartStatusFailed
	}

	// Check if the plugin configuration is valid json
	if !gjson.ValidBytes(configData) {
		proxywasm.LogErrorf("invalid JSON in plugin configuration: %v", string(configData))
		return types.OnPluginStartStatusFailed
	}
	jsonData := gjson.ParseBytes(configData)

	// Parsing the enabled field.
	p.config.Enabled = jsonData.Get("enabled").Bool()
	p.config.CorrelationHeader = jsonData.Get("correlationHeader").String()
	p.config.PropagationHeader = propagationHeader{
		Name:    jsonData.Get("propagationHeader.name").String(),
		Default: jsonData.Get("propagationHeader.default").String(),
	}

	proxywasm.LogDebugf("Parsed plugin config: %v", p.config)
	return types.OnPluginStartStatusOK
}

func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("********** NewHttpContext **********")
	return &httpContext{contextID: contextID, correlationHeader: p.config.CorrelationHeader, propagationHeader: p.config.PropagationHeader}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID         uint32
	correlationHeader string
	propagationHeader propagationHeader
}

// *********************************************
// REQUEST PATH
// *********************************************

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestHeaders **********")
	reqHeaders := properties.GetRequestHeaders()
	direction := properties.GetListenerDirection()
	proxywasm.LogInfof("GetRequestHeaders: %+v", reqHeaders)
	proxywasm.LogInfof("direction: %v", direction)

	cHeaderVal, ok := reqHeaders[strings.ToLower(ctx.correlationHeader)]
	if !ok {
		proxywasm.LogInfof("correlationHeader '%v' not available in the request headers", strings.ToLower(ctx.correlationHeader))
		return types.ActionContinue
	}
	proxywasm.LogInfof("correlationHeader '%v' available in the request headers with value '%v'", strings.ToLower(ctx.correlationHeader), cHeaderVal)

	switch direction.String() {
	case "UNSPECIFIED":
	case "INBOUND":
		pHeaderVal, ok := reqHeaders[strings.ToLower(ctx.propagationHeader.Name)]
		if !ok {
			if err := proxywasm.AddHttpRequestHeader(strings.ToLower(ctx.propagationHeader.Name), ctx.propagationHeader.Default); err != nil {
				proxywasm.LogErrorf("failed to set request header: %v", err)
			}
			proxywasm.LogInfof("inbound header '%v' set to default value: '%v'", strings.ToLower(ctx.propagationHeader.Name), ctx.propagationHeader.Default)
		} else {
			if err := utils.SetSharedDataSafe(cHeaderVal, []byte(pHeaderVal), 0); err != nil {
				proxywasm.LogErrorf("failed to set shared data key '%v' with value '%v': %v", cHeaderVal, pHeaderVal, err)
			}
		}
	case "OUTBOUND":
		_, ok := reqHeaders[strings.ToLower(ctx.propagationHeader.Name)]
		if !ok {
			pHeaderVal, _, err := utils.GetSharedDataSafe(cHeaderVal)
			if err != nil {
				proxywasm.LogErrorf("failed to get shared data key '%v': %v", cHeaderVal, err)
			}
			if err := proxywasm.AddHttpRequestHeader(strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal)); err != nil {
				proxywasm.LogErrorf("failed to set request header '%v' with value '%v': %v", strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal), err)
			}
			proxywasm.LogInfof("outbound header '%v' added with value '%v'", strings.ToLower(ctx.propagationHeader.Name), string(pHeaderVal))
		}
	}
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestTrailers(bodySize int) types.Action {
	return types.ActionContinue
}

// *********************************************
// RESPONSE PATH
// *********************************************

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")
	resHeaders := properties.GetResponseHeaders()
	direction := properties.GetListenerDirection()
	proxywasm.LogInfof("GetResponseHeaders: %+v", resHeaders)
	proxywasm.LogInfof("direction: %v", direction)

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseTrailers(bodySize int) types.Action {
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogInfo("********** OnHttpStreamDone **********")
}
