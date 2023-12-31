module geo-fetcher

go 1.19

require github.com/tetratelabs/proxy-wasm-go-sdk v0.22.0

require (
	geo-fetcher/utils v0.0.0
	github.com/tidwall/gjson v1.17.0
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

replace geo-fetcher/utils => ./utils
