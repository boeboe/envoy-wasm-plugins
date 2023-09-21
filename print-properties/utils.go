package main

import (
	"encoding/binary"
	"math"
	"strconv"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// convert array of tuplets to map
func tupletArrayToMap(inArray [][2]string) map[string]string {
	outMap := make(map[string]string)
	for _, tuple := range inArray {
		outMap[tuple[0]] = tuple[1]
	}
	return outMap
}

// safely get string value from map
func getStringValueFromMap(inMap map[string]string, key string) string {
	if outString, ok := inMap[key]; ok {
		return outString
	}
	proxywasm.LogWarnf("map '%v' does not contain key '%v'\n", inMap, key)
	return ""
}

// safely get int value from map
func getIntValueFromMap(inMap map[string]string, key string) int {
	if outString, ok := inMap[key]; ok {
		return int(float64fromByteArray([]byte(outString)))
	}
	proxywasm.LogWarnf("map '%v' does not contain key '%v'\n", inMap, key)
	return 0
}

// safely get bool value from map
func getBoolValueFromMap(inMap map[string]string, key string) bool {
	if outString, ok := inMap[key]; ok {
		outBool, err := strconv.ParseBool(outString)
		if err != nil {
			proxywasm.LogErrorf("error converting string '%v' to bool: %v", outString, err)
			return false
		}
		return outBool
	}
	proxywasm.LogWarnf("map '%v' does not contain key '%v'\n", inMap, key)
	return false
}

// convert byte array to float64
func float64fromByteArray(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
