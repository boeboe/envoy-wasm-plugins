package utils

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

const blockSize = 8

// SetSharedDataSafe safely encodes the input data by prefixing it with its original length (in 8 bytes)
// and then uses the SetSharedData function from the lower layer to store the data.
func SetSharedDataSafe(key string, data []byte, cas uint32) error {
	proxywasm.LogInfo("Encoding data for safe shared storage...")

	encodedData := encodeSharedData(data)
	err := proxywasm.SetSharedData(key, encodedData, cas)
	if err != nil {
		proxywasm.LogError(fmt.Sprintf("Failed to set shared data for key %s: %v", key, err))
	}
	return err
}

// GetSharedDataSafe uses the GetSharedData function from the lower layer to retrieve the data,
// and then decodes it to return the original data.
func GetSharedDataSafe(key string) ([]byte, uint32, error) {
	proxywasm.LogInfo(fmt.Sprintf("Retrieving data for key %s...", key))

	rawData, cas, err := proxywasm.GetSharedData(key)
	if err != nil {
		proxywasm.LogError(fmt.Sprintf("Failed to get shared data for key %s: %v", key, err))
		return nil, cas, err
	}

	decodedData, err := decodeSharedData(rawData)
	if err != nil {
		proxywasm.LogError(fmt.Sprintf("Failed to decode data for key %s: %v", key, err))
	}
	return decodedData, cas, err
}

func encodeSharedData(data []byte) []byte {
	totalLength := len(data) + blockSize
	padLength := (blockSize - (totalLength % blockSize)) % blockSize
	totalLength += padLength

	result := make([]byte, totalLength)

	binary.LittleEndian.PutUint64(result, uint64(len(data)))
	copy(result[blockSize:], data)

	return result
}

func decodeSharedData(data []byte) ([]byte, error) {
	if len(data) < blockSize {
		errMsg := "Invalid data length: received data is shorter than expected."
		proxywasm.LogError(errMsg)
		return nil, errors.New(errMsg)
	}

	originalLength := binary.LittleEndian.Uint64(data)
	if int(originalLength)+blockSize > len(data) {
		errMsg := fmt.Sprintf("Inconsistent data length: expected %d bytes but found more.", originalLength)
		proxywasm.LogError(errMsg)
		return nil, errors.New(errMsg)
	}

	return data[blockSize : blockSize+originalLength], nil
}
