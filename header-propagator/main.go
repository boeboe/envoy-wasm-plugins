package main

import (
	"header-propagator/properties"
	"header-propagator/utils"
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

const (
	Inbound  = "INBOUND"
	Outbound = "OUTBOUND"
)

type pluginConfig struct {
	CorrelationHeader   string            `json:"correlationHeader"`
	PropagationHeader   propagationHeader `json:"propagationHeader"`
	RequestPropagation  bool              `json:"requestPropagation"`
	ResponsePropagation bool              `json:"responsePropagation"`
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
	if pluginConfigurationSize == 0 {
		return types.OnPluginStartStatusOK
	}

	configData, err := proxywasm.GetPluginConfiguration()
	if err != nil || !gjson.ValidBytes(configData) {
		return types.OnPluginStartStatusFailed
	}

	jsonData := gjson.ParseBytes(configData)
	p.config = pluginConfig{
		RequestPropagation:  jsonData.Get("requestPropagation").Bool(),
		ResponsePropagation: jsonData.Get("responsePropagation").Bool(),
		CorrelationHeader:   jsonData.Get("correlationHeader").String(),
		PropagationHeader: propagationHeader{
			Name:    jsonData.Get("propagationHeader.name").String(),
			Default: jsonData.Get("propagationHeader.default").String(),
		},
	}

	return types.OnPluginStartStatusOK
}

func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID:           contextID,
		correlationHeader:   p.config.CorrelationHeader,
		propagationHeader:   p.config.PropagationHeader,
		requestPropagation:  p.config.RequestPropagation,
		responsePropagation: p.config.ResponsePropagation,
	}
}

type httpContext struct {
	types.DefaultHttpContext
	contextID           uint32
	correlationHeader   string
	propagationHeader   propagationHeader
	requestPropagation  bool
	responsePropagation bool
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	if !ctx.requestPropagation {
		return types.ActionContinue
	}

	reqHeaders := properties.GetRequestHeaders()
	direction := properties.GetListenerDirection()

	cHeaderVal, ok := reqHeaders[strings.ToLower(ctx.correlationHeader)]
	if !ok {
		return types.ActionContinue
	}

	switch direction.String() {
	case Inbound:
		handleInboundRequest(ctx, reqHeaders, cHeaderVal)
	case Outbound:
		handleOutboundRequest(ctx, reqHeaders, cHeaderVal)
	}

	return types.ActionContinue
}

func handleInboundRequest(ctx *httpContext, reqHeaders map[string]string, cHeaderVal string) {
	pHeaderName := strings.ToLower(ctx.propagationHeader.Name)
	pHeaderVal, ok := reqHeaders[pHeaderName]

	if !ok {
		setRequestHeader(pHeaderName, ctx.propagationHeader.Default)
		setSharedData(cHeaderVal, ctx.propagationHeader.Default)
	} else {
		setSharedData(cHeaderVal, pHeaderVal)
	}
}

func handleOutboundRequest(ctx *httpContext, reqHeaders map[string]string, cHeaderVal string) {
	pHeaderName := strings.ToLower(ctx.propagationHeader.Name)

	if _, ok := reqHeaders[pHeaderName]; !ok {
		pHeaderVal, err := getSharedData(cHeaderVal)
		if err == nil {
			setRequestHeader(pHeaderName, pHeaderVal)
		}
	}
}

func (ctx *httpContext) OnHttpResponseHeaders(numHeaders int, endOfStream bool) types.Action {
	if !ctx.responsePropagation {
		return types.ActionContinue
	}

	resHeaders := properties.GetResponseHeaders()
	direction := properties.GetListenerDirection()

	cHeaderVal, ok := resHeaders[strings.ToLower(ctx.correlationHeader)]
	if !ok {
		return types.ActionContinue
	}

	switch direction.String() {
	case Inbound:
		handleInboundResponse(ctx, resHeaders, cHeaderVal)
	case Outbound:
		handleOutboundResponse(ctx, resHeaders, cHeaderVal)
	}

	return types.ActionContinue
}

func handleInboundResponse(ctx *httpContext, resHeaders map[string]string, cHeaderVal string) {
	pHeaderName := strings.ToLower(ctx.propagationHeader.Name)

	if _, ok := resHeaders[pHeaderName]; !ok {
		pHeaderVal, err := getSharedData(cHeaderVal)
		if err == nil {
			setResponseHeader(pHeaderName, pHeaderVal)
		}
	}
}

func handleOutboundResponse(ctx *httpContext, resHeaders map[string]string, cHeaderVal string) {
	pHeaderName := strings.ToLower(ctx.propagationHeader.Name)

	if _, ok := resHeaders[pHeaderName]; !ok {
		pHeaderVal, err := getSharedData(cHeaderVal)
		if err == nil {
			setResponseHeader(pHeaderName, pHeaderVal)
		}
	}
}

// Helper functions to reduce repetition
func setRequestHeader(name, value string) {
	if err := proxywasm.AddHttpRequestHeader(name, value); err != nil {
		proxywasm.LogErrorf("Failed to set request header: %v", err)
	}
}

func setResponseHeader(name, value string) {
	if err := proxywasm.AddHttpResponseHeader(name, value); err != nil {
		proxywasm.LogErrorf("Failed to set response header: %v", err)
	}
}

func setSharedData(key, value string) {
	if err := utils.SetSharedDataSafe(key, []byte(value), 0); err != nil {
		proxywasm.LogErrorf("Failed to set shared data: %v", err)
	}
}

func getSharedData(key string) (string, error) {
	data, _, err := utils.GetSharedDataSafe(key)
	return string(data), err
}
