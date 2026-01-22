package main

import (
	"context"
	"fmt"
	"time"

	"github.com/reglet-dev/reglet-sdk/go/application/schema"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// ExamplePlugin implements the SDK's Plugin interface.
type ExamplePlugin struct{}

func (p *ExamplePlugin) Describe(ctx context.Context) (entities.Metadata, error) {
	return entities.Metadata{
		Name:        "Example Compliance Plugin",
		Description: "Verifies TLS certificate expiry using a custom host function",
		Version:     "1.0.0",
		Capabilities: []entities.Capability{
			entities.NewCapability("custom", "tls_check"),
		},
	}, nil
}

func (p *ExamplePlugin) Schema(ctx context.Context) ([]byte, error) {
	return schema.GenerateSchema(&PluginConfig{})
}

func (p *ExamplePlugin) Check(ctx context.Context, cfgMap map[string]any) (entities.Result, error) {
	// 1. Load and validate config
	cfg, err := LoadConfig(cfgMap)
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("config", err.Error())), nil
	}

	// 2. Perform TLS Check via Host Function
	resp, err := PerformTLSCheck(ctx, cfg.TargetHost, cfg.TargetPort)
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("network", err.Error())), nil
	}

	if resp.Error != nil {
		return entities.ResultError(entities.NewErrorDetail("network", resp.Error.Message).WithCode(resp.Error.Code)), nil
	}

	// 3. Evaluate
	notAfter, err := time.Parse(time.RFC3339, resp.NotAfter)
	if err != nil {
		return entities.ResultError(entities.NewErrorDetail("logic", "Failed to parse expiry date")), nil
	}

	daysRemaining := int(time.Until(notAfter).Hours() / 24)
	compliant := daysRemaining >= cfg.MinDays

	data := map[string]any{
		"target":         fmt.Sprintf("%s:%d", cfg.TargetHost, cfg.TargetPort),
		"not_after":      resp.NotAfter,
		"days_remaining": daysRemaining,
		"issuer":         resp.Issuer,
		"compliant":      compliant,
	}

	message := fmt.Sprintf("Certificate for %s expires in %d days", cfg.TargetHost, daysRemaining)

	if compliant {
		return entities.ResultSuccess(message, data), nil
	}

	return entities.ResultFailure(message, data), nil
}
