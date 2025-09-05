# GitVault

[![Build Status](https://github.com/gitvault/gitvault/workflows/CI/badge.svg)](https://github.com/gitvault/gitvault/actions)
[![Join the chat at https://img.shields.io/discord/322538954119184384.svg](https://img.shields.io/discord/322538954119184384.svg)](https://discord.gg/NsatcWJ)
[![Go Report Card](https://goreportcard.com/badge/code.gitea.io/gitea)](https://goreportcard.com/report/code.gitea.io/gitea)
[![Go Version](https://img.shields.io/github/go-mod/go-version/gitvault/gitvault)](https://github.com/gitvault/gitvault)
[![Release](https://img.shields.io/github/release/gitvault/gitvault.svg)](https://github.com/gitvault/gitvault/releases/latest)
[![LICENSE](https://img.shields.io/github/license/gitvault/gitvault.svg)](https://github.com/gitvault/gitvault/blob/main/LICENSE)
[![GitHub release](https://img.shields.io/github/downloads/gitvault/gitvault/total.svg)](https://github.com/gitvault/gitvault/releases/latest)
[![Docker Pulls](https://img.shields.io/docker/pulls/gitvault/gitvault.svg)](https://hub.docker.com/r/gitvault/gitvault)
[![GoDoc](https://godoc.org/code.gitea.io/gitea?status.svg)](https://godoc.org/code.gitea.io/gitea)

GitVault is a secure, private repository storage platform with integrated PhantomKit runtime capabilities. It provides a cloud-like storage solution for developers who want to securely store and import code projects and snippets.

## What is GitVault?

GitVault is a fork of Gitea that has been transformed into a specialized code storage and runtime platform. It combines the robust Git infrastructure of Gitea with PhantomKit's secure code execution capabilities.

### Key Features

- **Secure Code Storage**: Store proprietary source code in private repositories
- **PhantomKit Integration**: Import and execute code with language-specific loaders
- **Runtime Isolation**: V8 isolates and WASM support for secure code execution
- **MinIO Storage**: Scalable object storage backend
- **CLI Tools**: `phantom init`, `phantom load`, and other developer-friendly commands
- **Ephemeral Caching**: Smart caching with TTL for performance optimization

### PhantomKit Usage

```javascript
// Generated loader file for your project
import { PhantomKit } from 'phantomkit';

const phantom = new PhantomKit(process.env.PKIT_KEY);
await phantom.load('myScript');
phantom.loader
phantom.config
phantom.lock
```

## Quick Start

### Self-Hosted Installation

```bash
# Clone the repository
git clone https://github.com/gitvault/gitvault.git
cd gitvault

# Build the binary
make build

# Run GitVault
./gitvault web
```

### Docker Installation

```bash
docker run -d --name=gitvault -p 3000:3000 gitvault/gitvault:latest
```

### PhantomKit CLI

```bash
# Install PhantomKit CLI
npm install -g @gitvault/phantomkit

# Initialize a new project
phantom init my-project

# Load and execute code
phantom load myScript
```

## Architecture

GitVault extends Gitea's architecture with:

- **Vault Module**: Secure code storage and retrieval
- **Runtime Module**: V8 isolates and WASM runtime management
- **PhantomKit Module**: Core functionality for code execution
- **Storage Integration**: MinIO backend for scalable storage

## Development

### Prerequisites

- Go 1.21+
- Node.js 20+
- Git
- Make

### Building

```bash
make build
make test
make lint
```

### Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

## License

Portions of this software are derived from Gitea (https://gitea.io),
licensed under the MIT License. The MIT license applies only to those portions.

All other code, including but not limited to PhantomKit, GitVault SDK, CLI, and extensions, are proprietary to GitVault Technologies and licensed under the GitVault Proprietary License.

You may not copy, modify, distribute, sublicense, or reverse engineer
the proprietary portions of this software without explicit written permission.

For questions or commercial licensing, contact legal@gitvault.io.

## Acknowledgments

- Built on top of [Gitea](https://gitea.io/) - the community managed Git service
- PhantomKit runtime capabilities for secure code execution
- MinIO for scalable object storage

## Support

- [Documentation](https://docs.gitvault.io)
- [Issues](https://github.com/gitvault/gitvault/issues)
- [Discord](https://discord.gg/NsatcWJ)
- [Discussions](https://github.com/gitvault/gitvault/discussions)

---

**GitVault** - Secure Code Storage with PhantomKit Integration
