module geo-tagger

go 1.19

require github.com/tetratelabs/proxy-wasm-go-sdk v0.22.0

require (
	geo-tagger/utils v0.0.0
	github.com/stretchr/testify v1.8.4 // indirect
)

replace geo-tagger/utils => ./utils
