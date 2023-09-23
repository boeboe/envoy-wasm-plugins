package main

import (
	"encoding/binary"
	"math"
	"unsafe"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
)

// get string property
func getPropertyString(path []string) (string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// get uint64 property
func getPropertyUint64(path []string) (uint64, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeToFUint64(b), nil
}

// get float64 property
func getPropertyFloat64(path []string) (float64, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return 0, err
	}

	return deserializeToFloat64(b), nil
}

// // get complex property object as a map of byte slices
// // to be used when dealing with mixed type properties
// func getPropertyByteSliceMap(path []string) (map[string][]byte, error) {
// 	b, err := proxywasm.GetProperty(path)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return deserializeToByteMap(b), nil
// }

// get complex property object as a map of string
// to be used when dealing with string only type properties
func getPropertyStringMap(path []string) (map[string]string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeToStringMap(b), nil
}

// get complex property object as a string slice
func getPropertyStringSlice(path []string) ([]string, error) {
	b, err := proxywasm.GetProperty(path)
	if err != nil {
		return nil, err
	}

	return deserializeToStringSlice(b), nil
}

// deserialize byte slice to string slice
func deserializeToStringSlice(bs []byte) []string {
	numStrings := int(binary.LittleEndian.Uint32(bs[:4]))
	ret := make([]string, numStrings)
	idx := 4
	dataIdx := 4 + 8*numStrings
	for i := 0; i < numStrings; i++ {
		strLen := int(binary.LittleEndian.Uint64(bs[idx : idx+8]))
		idx += 8
		ret[i] = string(bs[dataIdx : dataIdx+strLen])
		dataIdx += strLen + 2
	}
	return ret
}

// deserialize byte array to float64
func deserializeToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

// deserialize byte array to float64
func deserializeToFUint64(bytes []byte) uint64 {
	return binary.LittleEndian.Uint64(bytes)
}

// // deserialize byte slice to key value map, used for mixed type maps
// //   - keys are always string
// //   - value are raw byte strings that need further parsing
// func deserializeToByteMap(bs []byte) map[string][]byte {
// 	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
// 	var sizeIndex = 4
// 	var dataIndex = 4 + 4*2*int(numHeaders)
// 	ret := make(map[string][]byte)
// 	for i := 0; i < int(numHeaders); i++ {
// 		keySize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
// 		sizeIndex += 4
// 		keyPtr := bs[dataIndex : dataIndex+keySize]
// 		key := *(*string)(unsafe.Pointer(&keyPtr))
// 		dataIndex += keySize + 1

// 		valueSize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
// 		sizeIndex += 4
// 		valuePtr := bs[dataIndex : dataIndex+valueSize]
// 		value := *(*[]byte)(unsafe.Pointer(&valuePtr))
// 		dataIndex += valueSize + 1
// 		ret[key] = value
// 	}
// 	return ret
// }

// deserialize byte array to key value map, used for string only type maps
//   - keys are always string
//   - value are always string
func deserializeToStringMap(bs []byte) map[string]string {
	numHeaders := binary.LittleEndian.Uint32(bs[0:4])
	var sizeIndex = 4
	var dataIndex = 4 + 4*2*int(numHeaders)
	ret := make(map[string]string)
	for i := 0; i < int(numHeaders); i++ {
		keySize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		keyPtr := bs[dataIndex : dataIndex+keySize]
		key := *(*string)(unsafe.Pointer(&keyPtr))
		dataIndex += keySize + 1

		valueSize := int(binary.LittleEndian.Uint32(bs[sizeIndex : sizeIndex+4]))
		sizeIndex += 4
		valuePtr := bs[dataIndex : dataIndex+valueSize]
		value := *(*string)(unsafe.Pointer(&valuePtr))
		dataIndex += valueSize + 1
		ret[key] = value
	}
	return ret
}
