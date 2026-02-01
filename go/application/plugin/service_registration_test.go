package plugin_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/plugin"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestService is a struct with methods to be registered as operations.
type TestService struct {
	plugin.Service `name:"test_service" desc:"Test Service"`
	EchoOp         plugin.Op `desc:"Echoes the message back" method:"Echo"`
	AddOp          plugin.Op `desc:"Adds two numbers" method:"Add"`
}

type EchoRequest struct {
	Message string `json:"message"`
}

// Echo is a valid operation.
// Handler signature must match: func(ctx, *Request) (*Result, error)
func (s *TestService) Echo(ctx context.Context, req *plugin.Request) (*entities.Result, error) {
	var body EchoRequest
	if len(req.Raw) > 0 {
		if err := json.Unmarshal(req.Raw, &body); err != nil {
			return nil, err
		}
	}

	// ResultSuccess expects Data as map[string]any
	resData := map[string]any{"reply": body.Message}
	res := entities.ResultSuccess("echoed", resData)
	return &res, nil
}

// Add is another valid operation.
func (s *TestService) Add(ctx context.Context, req *plugin.Request) (*entities.Result, error) {
	// Simplified for testing, ignoring input
	// entities.ResultSuccess takes map[string]any
	resData := map[string]any{"sum": 42}
	res := entities.ResultSuccess("added", resData)
	return &res, nil
}

func TestServiceRegistration(t *testing.T) {
	// 1. Define plugin
	def := plugin.DefinePlugin(plugin.PluginDef{
		Name:    "test-plugin",
		Version: "1.0.0",
	})

	// 2. Register the service
	svc := &TestService{}
	err := plugin.RegisterService(def, svc)
	require.NoError(t, err)

	// 3. Generate Manifest
	manifest := def.Manifest()
	assert.NotNil(t, manifest)

	// Verify Service Manifest
	svcManifest, ok := manifest.Services["test_service"]
	assert.True(t, ok, "service 'test_service' not found")

	// Check Operations
	foundEcho := false
	for _, op := range svcManifest.Operations {
		if op.Name == "echo_op" { // Converted from field name 'EchoOp'
			foundEcho = true
			assert.Equal(t, "Echoes the message back", op.Description)
		}
	}
	assert.True(t, foundEcho, "echo_op operation not found")

	// 5. Verify Handler Registration
	// Field EchoOp -> name echo_op
	handlerEcho, ok := def.GetHandler("test_service", "echo_op")
	assert.True(t, ok)
	assert.NotNil(t, handlerEcho)

	handlerAdd, ok := def.GetHandler("test_service", "add_op")
	assert.True(t, ok)
	assert.NotNil(t, handlerAdd)

	// 6. Execute Handler (Echo)
	reqBody := []byte(`{"message": "hello"}`)
	req := &plugin.Request{
		Raw: reqBody,
	}

	res, err := handlerEcho(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, entities.ResultStatusSuccess, res.Status)
	assert.Equal(t, "echoed", res.Message)

	// Verify output data
	// Result.Data is map[string]any
	assert.Equal(t, "hello", res.Data["reply"])
}
