package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/reglet-dev/reglet-sdk/go/host"
	"github.com/reglet-dev/reglet-sdk/go/hostfuncs"
)

func main() {
	// 1. Initialize Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	// 2. Define plugin path
	pluginPath := filepath.Join("..", "plugin", "plugin.wasm")
	if len(os.Args) > 1 {
		pluginPath = os.Args[1]
	}

	slog.Info("Host: Starting plugin execution", "path", pluginPath)

	// 3. Read WASM file
	wasmBytes, err := os.ReadFile(pluginPath)
	if err != nil {
		slog.Error("Failed to read plugin file", "error", err)
		os.Exit(1)
	}

	// 4. Create Host Function Registry

	registry, err := hostfuncs.NewRegistry(
		hostfuncs.WithBundle(hostfuncs.NetworkBundle()),
		hostfuncs.WithBundle(CustomTLSBundle()),
	)
	if err != nil {
		slog.Error("Failed to create registry", "error", err)
		os.Exit(1)
	}

	// 5. Initialize WASM Runtime via SDK Executor
	ctx := context.Background()
	executor, err := host.NewExecutor(ctx, host.WithHostFunctions(registry))
	if err != nil {
		slog.Error("Failed to create executor", "error", err)
		os.Exit(1)
	}
	defer executor.Close(ctx)

	// 6. Load Plugin
	plugin, err := executor.LoadPlugin(ctx, wasmBytes)
	if err != nil {
		slog.Error("Failed to load plugin", "error", err)
		os.Exit(1)
	}

	// 7. Execute Check
	config := map[string]any{
		"target_host": "example.com",
		"target_port": 443,
		"min_days":    30,
	}

	slog.Info("Host: Executing Check", "config", config)
	result, err := plugin.Check(ctx, config)
	if err != nil {
		slog.Error("Check execution failed", "error", err)
		os.Exit(1)
	}

	// 8. Print Result
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	slog.Info("Host: Plugin Result Status", "status", result.Status)
	fmt.Printf("Host: Result: \n%s\n", string(resultJSON))
}
