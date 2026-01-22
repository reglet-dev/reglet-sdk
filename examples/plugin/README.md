# Example Compliance Plugin

This example demonstrates how to build a Reglet compliance plugin in Go that uses a custom host function.

## Overview

The plugin performs a **TLS Certificate Expiry Check**. It verifies that the TLS certificate of a target host is valid and has a minimum number of days remaining before expiry.

Key SDK features demonstrated:
- Implementing the `Plugin` interface (`Describe`, `Schema`, `Check`).
- Using `config` helpers for structured configuration.
- Calling a **custom host function** (`tls_check`) defined by the host runtime.
- Returning structured `Evidence` (Result).

## Structure

- `main.go`: Plugin entrypoint and registration.
- `plugin.go`: Main logic for the compliance check.
- `config.go`: Configuration parsing.
- `tls_adapter.go`: Low-level bridge to the host's `tls_check` function.

## Building

To compile the plugin to WASM:

```bash
GOOS=wasip1 GOARCH=wasm go build -buildmode=c-shared -o plugin.wasm .
```

*Note: `-buildmode=c-shared` is required for Go wasip1 plugins to export functions correctly as a library (reactor module).*

## Running

This plugin is designed to be run by the `examples/host-runtime`. See the root `README.md` or `Makefile` for end-to-end execution.
