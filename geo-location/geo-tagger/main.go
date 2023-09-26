package main

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

const geoDBKey = "geolocation_db"
const geoDBUpdateQueue = "geolocation_update_queue"

type httpFilter struct {
	proxywasm.DefaultHttpContext
	geoDB *geoip2.Reader
}

func main() {
	proxywasm.SetHttpContext(new(httpFilter))
}

func (ctx *httpFilter) OnHttpRequestHeaders(numHeaders int, endOfStream bool) proxywasm.Action {
	ipStr := ctx.GetHttpRequestHeader("X-Forwarded-For")
	if ipStr == "" {
		return proxywasm.ActionContinue
	}

	ip := net.ParseIP(ipStr)
	city, err := ctx.geoDB.City(ip)
	if err != nil {
		// Handle error
		return proxywasm.ActionContinue
	}

	// Use the geolocation data (e.g., city.Country.IsoCode, city.City.Names["en"])
	// ...

	return proxywasm.ActionContinue
}

func (ctx *httpFilter) OnQueueReady(queueID uint32) {
	data, err := proxywasm.DequeueSharedQueue(queueID)
	if err != nil {
		return
	}

	if string(data) == "update_available" {
		ctx.updateGeoDBFromSharedMemory()
	}
}

func (ctx *httpFilter) updateGeoDBFromSharedMemory() {
	data, _, err := proxywasm.GetSharedData(geoDBKey)
	if err != nil {
		return
	}

	// Update the local geolocation database with the fetched data
	ctx.geoDB, err = geoip2.FromBytes(data)
	if err != nil {
		// Handle error
	}
}
