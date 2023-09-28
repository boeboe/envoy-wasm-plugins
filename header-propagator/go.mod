module header-propagator

go 1.19

require github.com/tetratelabs/proxy-wasm-go-sdk v0.22.0

require github.com/tidwall/gjson v1.17.0

require (
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
)

require header-propagator/properties v0.0.0
require header-propagator/utils v0.0.0

replace header-propagator/properties => ./properties
replace header-propagator/utils => ./utils
