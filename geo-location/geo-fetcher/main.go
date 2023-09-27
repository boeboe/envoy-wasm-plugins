package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"geo-fetcher/utils"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/tidwall/gjson"
)

const (
	geoDBKey             = "geolocation_db"
	geoDBUpdateQueue     = "geolocation_update_queue"
	sharedDataPadByte    = byte(0) // Default padding byte is 0
	sharedDataTargetSize = 8       // Desired length or its multiple for the byte slice
	geoDBTaggerVMID      = "geo-tagger"
)

// vmContext is the main context for the VM.
type vmContext struct {
	types.DefaultVMContext
}

func main() {
	proxywasm.SetVMContext(&vmContext{})
}

// geoFetchConfig represents the configuration for fetching the GeoDB.
type geoFetchConfig struct {
	GeoDBURLPath    string `json:"geo_db_url_path"`
	PollingInterval uint32 `json:"polling_interval"`
}

// pluginContext represents the context for the plugin.
type pluginContext struct {
	types.DefaultPluginContext
	config  geoFetchConfig
	queueID uint32
}

// NewPluginContext creates a new plugin context.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

// OnPluginStart initializes the plugin.
func (ctx *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	// Fetch the plugin configuration
	data, err := proxywasm.GetPluginConfiguration()
	if err != nil && err != types.ErrorStatusNotFound {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfo("successfully read plugin configuration")

	// Parse the configuration data
	config, err := parseGeoServiceConfiguration(data)
	if err != nil {
		proxywasm.LogCriticalf("error parsing plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.config = config
	proxywasm.LogInfof("successfully parsed plugin configuration: %+v", config)

	// Initialize the shared data with an empty value
	err = utils.SetSharedDataSafe(geoDBKey, []byte{}, 0)
	if err != nil {
		proxywasm.LogCriticalf("failed to initialize shared data: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("successfully initialize shared data with key: %v", geoDBKey)

	// Register the shared queue here
	queueID, err := proxywasm.RegisterSharedQueue(geoDBUpdateQueue)
	if err != nil {
		proxywasm.LogCriticalf("failed to register shared queue: %v", err)
		return types.OnPluginStartStatusFailed
	}
	ctx.queueID = queueID
	proxywasm.LogInfof("successfully register shared queue: %v", geoDBUpdateQueue)

	// Set the tick period for periodic tasks
	if err := proxywasm.SetTickPeriodMilliSeconds(config.PollingInterval); err != nil {
		proxywasm.LogCriticalf("failed to set tick period: %v", err)
		return types.OnPluginStartStatusFailed
	}
	proxywasm.LogInfof("successfully set tick period milliseconds: %d", config.PollingInterval)

	return types.OnPluginStartStatusOK
}

// parseGeoServiceConfiguration parses the configuration for the Geo service.
func parseGeoServiceConfiguration(data []byte) (geoFetchConfig, error) {
	if len(data) == 0 {
		return geoFetchConfig{}, nil
	}
	if !gjson.ValidBytes(data) {
		return geoFetchConfig{}, fmt.Errorf("the plugin configuration is not a valid json: %q", string(data))
	}
	jsonData := gjson.ParseBytes(data)
	return geoFetchConfig{
		PollingInterval: uint32(jsonData.Get("polling_interval").Uint()),
		GeoDBURLPath:    jsonData.Get("geo_db_url_path").Str,
	}, nil
}

// OnTick is called periodically based on the tick period set.
func (ctx *pluginContext) OnTick() {
	// Fetch the GeoDB data on every tick
	ctx.fetchGeoDB()
}

// fetchGeoDB fetches the GeoDB data from the provided URL.
func (ctx *pluginContext) fetchGeoDB() {
	headers := [][2]string{
		{":authority", "download.db-ip.com"},
		{":method", "GET"},
		{":path", ctx.config.GeoDBURLPath},
		{":scheme", "https"},
	}

	callback := func(numHeaders, bodySize, numTrailers int) {
		body, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
		if err != nil {
			proxywasm.LogErrorf("failed to get response body: %v", err)
			return
		}
		proxywasm.LogInfof("fetchGeoDB response body size: %v", len(body))

		headers, err := proxywasm.GetHttpCallResponseHeaders()
		if err != nil {
			proxywasm.LogErrorf("failed to get response headers: %v", err)
			return
		}
		proxywasm.LogInfof("fetchGeoDB response headers: %+v", headers)

		// Decompress the gzipped data
		gz, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			proxywasm.LogErrorf("failed to create gzip reader: %v", err)
			return
		}
		defer gz.Close()

		extractedData, err := io.ReadAll(gz)
		if err != nil {
			proxywasm.LogErrorf("failed to read gzipped data: %v", err)
		}
		proxywasm.LogInfof("successfully read gzipped data of size: %v", len(extractedData))

		if err := ctx.storeInSharedMemory(extractedData); err != nil {
			proxywasm.LogErrorf("failed to store data in shared memory: %v", err)
			return
		}
		proxywasm.LogInfof("successfully stored data of size %v in shared memory", len(extractedData))
		ctx.notifyHTTPFilter()
	}

	if _, err := proxywasm.DispatchHttpCall("db-ip", headers, nil, nil, 5000, callback); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
	}
}

// storeInSharedMemory stores the fetched data in shared memory.
func (ctx *pluginContext) storeInSharedMemory(data []byte) error {
	err := utils.SetSharedDataSafe(geoDBKey, data, 0)
	if err != nil {
		return fmt.Errorf("failed to set shared data: %v", err)
	}
	return nil
}

// notifyHTTPFilter notifies the HTTP filter about the updated data.
func (ctx *pluginContext) notifyHTTPFilter() {
	if err := proxywasm.EnqueueSharedQueue(ctx.queueID, []byte("update_available")); err != nil {
		proxywasm.LogErrorf("failed to enqueue shared queue: %v", err)
	}
}
