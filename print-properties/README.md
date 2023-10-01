# Print Properties WASM Plugin

This repository contains a plugin designed to explore the WASM functionality of Envoy by printing various [WASM properties/attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes) based on the configuration provided.

## Table of Contents

- [Features](#features)
- [Local Development](#local-development)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Configuration](#configuration)
- [Makefile Commands](#makefile-commands)
- [Releases](#releases)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- Print [WASM properties/attributes](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes) based on plugin configuration.

## Local Development

### Prerequisites

- [TinyGo](https://tinygo.org) (for compiling the WASM binary)
- Docker (for building and pushing the container)
- `kubectl` (for deploying to a Kubernetes cluster)
- An Envoy or Istio setup that supports WASM plugins.

### Installation

1. Clone the repository:

```bash
git clone https://github.com/boeboe/envoy-wasm-plugins.git
cd envoy-wasm-plugins/print-properties
```

2. Compile the WASM binary:

```bash
make compile
```

3. Build and release the docker container or the wasm binary:

```bash
make release-docker
make release-github
```

4. Deploy the demo applications and WASM filters:

```bash
make deploy
```

5. Use the provided `curl` commands for manual verification

```bash
make curl

# Call to ingress that forwards to http.org on the internet
curl -v -H "X-B3-Sampled: 1" --resolve "wasm.httpbin.org:80:172.18.0.101" "http://wasm.httpbin.org/headers"

# Call to ingress that forwards to tetrate.io on the internet
curl -v -H "X-B3-Sampled: 1" --resolve "wasm.tetrate.io:80:172.18.0.101" "http://wasm.tetrate.io/"
```

## Configuration

The plugin can be configured using a JSON configuration. The configuration specifies which properties should be printed for various Plugin and HTTP events. The available Plugin and HTTP events include:

| Parameter | Description | Type |
|-----------|-------------|------|
| `onPluginStart` | Called for all plugin contexts (after OnVmStart if this is the VM context) | `propertiesPrinting` struct |
| `onHttpRequestHeaders` | Called when request headers arrive | `propertiesPrinting` struct |
| `onHttpRequestBody` | Called when a request body *frame* arrives | `propertiesPrinting` struct |
| `onHttpRequestTrailers` | Called when request trailers arrive | `propertiesPrinting` struct |
| `onHttpResponseHeaders` | Called when response headers arrive | `propertiesPrinting` struct |
| `onHttpResponseBody` | Called when a response body *frame* arrives | `propertiesPrinting` struct |
| `onHttpResponseTrailers` | Called when response trailers arrive | `propertiesPrinting` struct |
| `onHttpStreamDone` | Called before the host deletes this context. You can retrieve the HTTP request/response information (such as headers, etc.) during this call. This can be used to implement logging features | propertiesPrinting struct |

The `propertiesPrinting` struct contains the available properties that can be printed:

| Parameter | Description | Type | Documentation Link |
|-----------|-------------|------|--------------------|
| `printWasmProperties` | Print Wasm Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#wasm-attributes)                  |
| `printNodeMetadataProperties` | Print Node Metadata Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/core/v3/base.proto#envoy-v3-api-msg-config-core-v3-node)   |
| `printNodeProxyConfigProperties` | Print Node Proxy Config Properties | bool | [docs](https://istio.io/latest/docs/reference/config/istio.mesh.v1alpha1/#ProxyConfig)                                       |
| `printXdsProperties` | Print XDS Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#configuration-attributes)         |
| `printUpstreamProperties` | Print Upstream Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#upstream-attributes)              |
| `printConnectionProperties` | Print Connection Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#connection-attributes)            |
| `printResponseProperties` | Print Response Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#response-attributes)              |
| `printRequestProperties` | Print Request Properties | bool | [docs](https://www.envoyproxy.io/docs/envoy/latest/intro/arch_overview/advanced/attributes#request-attributes)               |

To configure the `print-properties` wasm plugin within istio, apply the following [`WasmPlugin`](https://istio.io/latest/docs/reference/config/proxy_extensions/wasm-plugin/) configuration:

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: print-properties
  namespace: istio-ingress
spec:
  imagePullPolicy: Always
  phase: AUTHN
  pluginConfig:
    onPluginStart:
      printWasmProperties: true
      printNodeMetadataProperties: true
      printNodeProxyConfigProperties: true
      printXdsProperties: true
    onHttpRequestHeaders:
      printUpstreamProperties: true
      printRequestProperties: true
    onHttpResponseHeaders:
      printConnectionProperties: true
      printResponseProperties: true
  pluginName: print-properties
  selector:
    matchLabels:
      istio: ingress
  url: oci://docker.io/boeboe/envoy-wasm-plugins:print-properties-0.1
  vmConfig:
    env:
      - name: POD_NAME
        valueFrom: HOST
      - name: TEST_ENV_VAR
        value: testEnvValue
```

## Makefile Commands

The `makefile` has a self describing help target, set as default target.

```bash
$ make

help                           This help
all                            Rebuild, rerelease and redeploy
compile                        Compile the wasm binary
test                           Run golang tests
clean                          Remove output artifats
release-docker                 Release docker container on dockerhub
release-github                 Release wasm binary on github
deploy                         Deploy the demo applications and wasm filters
reboot-pods                    Force reboot workload pods
dump-config                    Dump ingress and sidecar envoy configs
dump-logs                      Dump ingress and sidecar logs
enable-full-debug              Enable debug:all on ingress and sidecar
enable-full-info               Enable info:all on ingress and sidecar
enable-http-debug              Enable debug:http on ingress and sidecar
enable-wasm-debug              Enable debug:wasm on ingress and sidecar
curl                           Print some sample curl commands for manual verification
```

## Releases

- Docker Container: [DockerHub](https://hub.docker.com/r/boeboe/envoy-wasm-plugins/tags)
- WASM Binary: [GitHub Releases](https://github.com/boeboe/envoy-wasm-plugins/releases)

## Contributing

Contributions are welcome! PR's are encouraged.

## License

This project is licensed under the Apache License. See the [LICENSE](./LICENSE) file for details.

## Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/boeboe/envoy-wasm-plugins/issues).