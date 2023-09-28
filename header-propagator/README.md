# Header Propagation WASM Plugin

This repository contains a plugin designed to enhance the functionality of Envoy by adding custom header propagation.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Makefile Commands](#makefile-commands)
- [Releases](#releases)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- Custom header propagation for both inbound and outbound requests and responses.
- Dynamic configuration through JSON.
- Detailed logging for debugging and monitoring.

## Getting Started

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

3. Build the Docker container:

```bash
make release-docker
```

4. Deploy the demo applications and WASM filters:

```bash
make deploy
```

## Usage

Use the provided curl commands for manual verification

```bash
make curl
```

## Makefile Commands

`make help`: Display the help menu.
`make all`: Rebuild, rerelease, and redeploy.
`make compile`: Compile the WASM binary.
`make test`: Run GoLang tests.
`make clean`: Remove output artifacts.
`make release-docker`: Release the Docker container on DockerHub.
`make release-github`: Release the WASM binary on GitHub.
`make deploy`: Deploy the demo applications and WASM filters.
`make dump-config`: Dump ingress and sidecar Envoy configs.
`make dump-logs`: Dump ingress and sidecar logs.
`make enable-full-debug`: Enable debug:all on ingress and sidecar.
`make enable-full-info`: Enable info:all on ingress and sidecar.
`make enable-wasm-debug`: Enable debug:wasm on ingress and sidecar.
`make curl`: Print some sample curl commands for manual verification.

## Releases

Docker Container: [DockerHub](https://hub.docker.com/r/boeboe/envoy-wasm-plugins/tags)
WASM Binary: [GitHub Releases](https://github.com/boeboe/envoy-wasm-plugins/releases)

## Contributing

Contributions are welcome! PR's are encourared.

## License

This project is licensed under the Apache License. See the [LICENSE](./LICENSE) file for details.

## Contact

For any questions or feedback, please open an issue on the [GitHub repository](https://github.com/boeboe/envoy-wasm-plugins/issues).