package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"

	"github.com/tidwall/gjson"
)

const geoDBKey = "geolocation_db"
const geoDBUpdateQueue = "geolocation_update_queue"

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
	err = proxywasm.SetSharedData(geoDBKey, padToCorrectLength([]byte{}), 0)
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

const (
	sharedDataPadByte    = byte(0) // Default padding byte is 0
	sharedDataTargetSize = 8       // Desired length or its multiple for the byte slice
)

// PadToCorrectLength prepares a byte slice such that its length is a multiple of the desired
// targetSize (in this case 8). This is done by padding the byte slice with a specified padByte.
func padToCorrectLength(data []byte) []byte {
	if len(data) == 0 {
		return make([]byte, sharedDataTargetSize)
	}
	padding := len(data) % sharedDataTargetSize
	if padding != 0 {
		padding = sharedDataTargetSize - padding
	}
	paddedData := make([]byte, len(data)+padding)
	copy(paddedData, data)
	for i := len(data); i < len(paddedData); i++ {
		paddedData[i] = sharedDataPadByte
	}
	return paddedData
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
	data, err := ctx.fetchGeoDB()
	if err != nil {
		proxywasm.LogErrorf("failed to fetch GeoDB: %v", err)
		return
	}
	proxywasm.LogInfo("successfully fetched GeoDB")
	if err := ctx.storeInSharedMemory(data); err != nil {
		proxywasm.LogErrorf("failed to store data in shared memory: %v", err)
		return
	}
	proxywasm.LogInfof("successfully stored data of size %v in shared memory", len(data))
	ctx.notifyHTTPFilter()
}

// fetchGeoDB fetches the GeoDB data from the provided URL.
func (ctx *pluginContext) fetchGeoDB() ([]byte, error) {
	headers := [][2]string{
		{":authority", "download.db-ip.com"},
		{":method", "GET"},
		{":path", ctx.config.GeoDBURLPath},
		{":scheme", "https"},
	}

	var fetchedData []byte
	callback := func(numHeaders, bodySize, numTrailers int) {
		body, err := proxywasm.GetHttpCallResponseBody(0, bodySize)
		if err != nil {
			proxywasm.LogErrorf("failed to get response body: %v", err)
			return
		}

		// Decompress the gzipped data
		gz, err := gzip.NewReader(bytes.NewReader(body))
		if err != nil {
			proxywasm.LogErrorf("failed to create gzip reader: %v", err)
			return
		}
		defer gz.Close()

		fetchedData, err = io.ReadAll(gz)
		if err != nil {
			proxywasm.LogErrorf("failed to read gzipped data: %v", err)
		}
		proxywasm.LogInfof("successfully read gzipped data of size: %v", len(fetchedData))
	}

	if _, err := proxywasm.DispatchHttpCall("db-ip", headers, nil, nil, 5000, callback); err != nil {
		proxywasm.LogCriticalf("dispatch httpcall failed: %v", err)
		return nil, err
	}

	return fetchedData, nil
}

// storeInSharedMemory stores the fetched data in shared memory.
func (ctx *pluginContext) storeInSharedMemory(data []byte) error {
	err := proxywasm.SetSharedData(geoDBKey, padToCorrectLength(data), 0)
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
