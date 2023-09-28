# Header Propagation WASM Plugin

This repository contains a plugin designed to enhance the functionality of Envoy by adding custom header propagation.

## Table of Contents

- [Features](#features)
- [Local Development](#development)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Configuration](#configuration)
- [Makefile Commands](#makefile-commands)
- [Releases](#releases)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- Custom header propagation for both inbound and outbound requests and responses.
- Dynamic configuration through JSON.
- Detailed logging for debugging and monitoring.

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
cd envoy-wasm-plugins
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

# Single request...
curl -v -H "x-b3-sampled: 1" -H "x-request-id: 31a9dd6e-2ac1-4058-96bd-30a9f7938714" -H "x-tetrate-swimlaneid: lane-a" --resolve "propagate.tetrate.io:80:172.18.0.101" "http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers"
curl -v -H "x-b3-sampled: 1" -H "x-request-id: 31a9dd6e-2ac1-4058-96bd-30a9f7938714" -H "x-tetrate-swimlaneid: lane-a" --resolve "propagate.tetrate.io:80:172.18.0.101" "http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/response-headers?x-request-id=31a9dd6e-2ac1-4058-96bd-30a9f7938714"

# Execute in a loop...
while true; do curl -v -H "x-b3-sampled: 1" -H "x-request-id: 31a9dd6e-2ac1-4058-96bd-30a9f7938714" -H "x-tetrate-swimlaneid: lane-a" --resolve "propagate.tetrate.io:80:172.18.0.101" "http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers" ; sleep 1 ; done
while true; do curl -v -H "x-b3-sampled: 1" -H "x-request-id: 31a9dd6e-2ac1-4058-96bd-30a9f7938714" -H "x-tetrate-swimlaneid: lane-a" --resolve "propagate.tetrate.io:80:172.18.0.101" "http://propagate.tetrate.io/proxy/app-b.ns-b/proxy/httpbin.ns-httpbin/headers/response-headers?x-request-id=31a9dd6e-2ac1-4058-96bd-30a9f7938714" ; sleep 1 ; done
```

## Configuration:

The `pluginConfig` section provides specific configurations for the wasm plugin, determining its behavior and processing of headers.

| Parameter             | Description                                                                                   | Example               |
|-----------------------|-----------------------------------------------------------------------------------------------|-----------------------|
| `correlationHeader`   | Specifies the name of the header that the plugin will use for correlation purposes            | `x-request-id`        |
| `propagationHeader`   | This subsection provides details about the header used for propagation                        |                       |
| `propagationHeader.default`             | The default value to be used for the propagation header if it's not present in the request    | `lane-a`              |
| `propagationHeader.name`                | The name of the header that will be used for propagation purposes                             | `x-tetrate-swimlaneid`|
| `requestPropagation`  | A boolean flag that determines if the plugin should handle propagation for incoming requests  | `true` or `false`     |
| `responsePropagation` | A boolean flag that determines if the plugin should handle propagation for outgoing responses | `true` or `false`     |


To configure the `header-propagator` wasm plugin within istio, apply the following `WasmPlugin` configuration:

```yaml
apiVersion: extensions.istio.io/v1alpha1
kind: WasmPlugin
metadata:
  name: header-propagator
  namespace: istio-system
spec:
  imagePullPolicy: Always
  pluginConfig:
    correlationHeader: x-request-id
    propagationHeader:
      default: lane-a
      name: x-tetrate-swimlaneid
    requestPropagation: true
    responsePropagation: true
  pluginName: header-propagator
  selector:
    matchLabels:
      header-propagation.tetrate.io/enabled: 'true'
  url: oci://docker.io/boeboe/envoy-wasm-plugins:header-propagator-0.1
```

Make sure the appropriate label selector is configured on your pods or deployments.

## Makefile Commands

- `make help` : Display the help menu.
- `make all` : Rebuild, rerelease, and redeploy.
- `make compile` : Compile the WASM binary.
- `make test` : Run GoLang tests.
- `make clean` : Remove output artifacts.
- `make release-docker` : Release the Docker container on DockerHub.
- `make release-github` : Release the WASM binary on GitHub.
- `make deploy` : Deploy the demo applications and WASM filter.
- `reboot-pods` : Force reboot workload pods.
- `make dump-config` : Dump ingress and sidecar Envoy configs.
- `make dump-logs` : Dump ingress and sidecar logs.
- `make enable-full-debug` : Enable debug:all on ingress and sidecar.
- `make enable-full-info` : Enable info:all on ingress and sidecar.
- `make enable-http-debug` : Enable debug:http on ingress and sidecar.
- `make enable-wasm-debug` : Enable debug:wasm on ingress and sidecar.
- `make curl` : Print some sample curl commands for manual verification.

## Releases

- Docker Container: [DockerHub](https://hub.docker.com/r/boeboe/envoy-wasm-plugins/tags)
- WASM Binary: [GitHub Releases](https://github.com/boeboe/envoy-wasm-plugins/releases)

## Contributing

Contributions are welcome! PR's are encouraged.

## License

This project is licensed under the Apache License. See the [LICENSE](./LICENSE) file for details.

## Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/boeboe/envoy-wasm-plugins/issues).