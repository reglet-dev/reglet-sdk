# Example Host Runtime

This example demonstrates how to build a WASM host runtime that integrates the Reglet SDK's host functions and executes plugins.

## Overview

The host runtime:
- Initializes a `wazero` WASM environment.
- Registers standard SDK bundles (`NetworkBundle`).
- Defines and registers a **custom host function** (`tls_check`).
- Loads a compiled plugin and executes its `Check` method.

Key integration points:
- `hostfuncs.NewRegistry`: Creating the collection of available host functions.
- `executor.go`: Boilerplate for mapping `HandlerRegistry` to wazero host modules.
- `custom_bundle.go`: Implementing a domain-specific host capability (TLS inspection).

## Running

1. Build the example plugin:
   ```bash
   cd ../plugin
   GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.wasm .
   ```

2. Run the host:
   ```bash
   go run .
   ```

Or use the root `Makefile`:
```bash
make example
```
