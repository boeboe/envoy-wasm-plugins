package main

import (
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"geo-tagger/utils"
)

const (
	geoDBKey         = "geolocation_db"
	geoDBUpdateQueue = "geolocation_update_queue"
	geoDBFetcherVMID = "geo-fetcher"
)

// vmContext is the main context for the VM.
type vmContext struct {
	types.DefaultVMContext
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

// pluginContext represents the context for the plugin.
type pluginContext struct {
	types.DefaultPluginContext
	queueID uint32
}

// NewPluginContext creates a new plugin context.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// OnPluginStart reads the shared geolocation data and registers to the shared queue.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogInfo("********** OnPluginStart *********")

	// Read the shared geolocation data
	data, _, err := proxywasm.GetSharedData(geoDBKey)
	if err != nil {
		proxywasm.LogCriticalf("error reading shared data: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("successfully read shared data of size: %d", len(data))

	// Register to the shared queue
	queueID, err := proxywasm.ResolveSharedQueue(geoDBFetcherVMID, geoDBUpdateQueue)
	if err != nil {
		proxywasm.LogCriticalf("failed to register shared queue: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.queueID = queueID
	proxywasm.LogInfof("successfully registered to shared queue: %v", geoDBUpdateQueue)

	return types.OnPluginStartStatusOK
}

// OnQueueReady handles messages from the shared queue.
func (ctx *pluginContext) OnQueueReady(queueID uint32) {
	proxywasm.LogInfo("********** OnQueueReady *********")
	for {
		data, err := proxywasm.DequeueSharedQueue(ctx.queueID)
		if err != nil {
			if err == types.ErrorStatusEmpty {
				// No more data in the queue
				break
			}
			proxywasm.LogErrorf("error dequeuing shared queue: %v", err)
			return
		}

		if string(data) == "update_available" {
			geoData, _, err := utils.GetSharedDataSafe(geoDBKey)
			if err != nil {
				proxywasm.LogErrorf("error reading updated shared data: %v", err)
				return
			}
			proxywasm.LogInfof("successfully read updated shared data of size: %d", len(geoData))
			// Handle the updated geolocation data as needed
		}
	}
}
