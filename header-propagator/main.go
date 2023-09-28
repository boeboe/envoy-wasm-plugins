package main

import (
	"header-propagator/properties"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/tidwall/gjson"
)

type PluginConfig struct {
	CorrelationHeaders []struct {
		Name string `json:"name"`
	} `json:"correlationHeaders"`
	PropagationHeaders []struct {
		Name    string `json:"name"`
		Default string `json:"default"`
	} `json:"propagationHeaders"`
	Enabled bool `json:"enabled"`
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
	config PluginConfig
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

	// For arrays, we need to handle them separately.
	jsonData.Get("correlationHeaders").ForEach(func(_, v gjson.Result) bool {
		header := struct {
			Name string `json:"name"`
		}{
			Name: v.Get("name").String(),
		}
		p.config.CorrelationHeaders = append(p.config.CorrelationHeaders, header)
		return true
	})

	jsonData.Get("propagationHeaders").ForEach(func(_, v gjson.Result) bool {
		header := struct {
			Name    string `json:"name"`
			Default string `json:"default"`
		}{
			Name:    v.Get("name").String(),
			Default: v.Get("default").String(),
		}
		p.config.PropagationHeaders = append(p.config.PropagationHeaders, header)
		return true
	})

	proxywasm.LogDebugf("Parsed plugin config: %v", p.config)

	return types.OnPluginStartStatusOK
}

func (*pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	proxywasm.LogInfo("********** NewHttpContext **********")
	return &httpContext{contextID: contextID}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID uint32
}

// *********************************************
// REQUEST PATH
// *********************************************

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestHeaders **********")

	reqHeaders := properties.GetRequestHeaders()
	proxywasm.LogInfof("GetRequestHeaders: %+v", reqHeaders)

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestBody **********")

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpRequestTrailers(bodySize int) types.Action {
	proxywasm.LogInfo("********** OnHttpRequestTrailers **********")

	return types.ActionContinue
}

// *********************************************
// RESPONSE PATH
// *********************************************

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseHeaders **********")

	resHeaders := properties.GetResponseHeaders()
	proxywasm.LogInfof("GetResponseHeaders: %+v", resHeaders)

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseBody(bodySize int, endOfStream bool) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseBody **********")

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseTrailers(bodySize int) types.Action {
	proxywasm.LogInfo("********** OnHttpResponseTrailers **********")

	return types.ActionContinue
}

func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogInfo("********** OnHttpStreamDone **********")
}
