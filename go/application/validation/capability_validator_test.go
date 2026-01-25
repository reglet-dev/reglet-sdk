package validation_test

import (
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/application/validation"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRegistry struct {
	schemas map[string]string
}

func (m *mockRegistry) Register(name string, capability interface{}) error { return nil }
func (m *mockRegistry) GetSchema(name string) (string, bool) {
	s, ok := m.schemas[name]
	return s, ok
}
func (m *mockRegistry) List() []string { return nil }

func TestCapabilityValidator_Validate(t *testing.T) {
	registry := &mockRegistry{
		schemas: map[string]string{
			"network": `{"type": "object", "properties": {"rules": {"type": "array"}}}`,
			"fs":      `{"type": "object", "required": ["rules"], "properties": {"rules": {"type": "array"}}}`,
		},
	}
	validator := validation.NewCapabilityValidator(registry)

	t.Run("Valid Manifest", func(t *testing.T) {
		manifest := &entities.PluginManifest{
			Capabilities: &entities.GrantSet{
				Network: &entities.NetworkCapability{
					Rules: []entities.NetworkRule{
						{Hosts: []string{"example.com"}, Ports: []string{"443"}},
					},
				},
			},
		}
		res, err := validator.Validate(manifest)
		require.NoError(t, err)
		assert.True(t, res.Valid)
		assert.Empty(t, res.Errors)
	})

	t.Run("Invalid Capability Schema", func(t *testing.T) {
		// missing required 'rules' in fs
		manifest := &entities.PluginManifest{
			Capabilities: &entities.GrantSet{
				FS: &entities.FileSystemCapability{},
			},
		}
		res, err := validator.Validate(manifest)
		require.NoError(t, err)
		assert.False(t, res.Valid)
		assert.NotEmpty(t, res.Errors)
		assert.Equal(t, "fs", res.Errors[0].Field)
	})

	t.Run("Unknown Capability", func(t *testing.T) {
		// Env not in registry
		manifest := &entities.PluginManifest{
			Capabilities: &entities.GrantSet{
				Env: &entities.EnvironmentCapability{},
			},
		}
		res, err := validator.Validate(manifest)
		require.NoError(t, err)
		assert.False(t, res.Valid)
		assert.Contains(t, res.Errors[0].Message, "no schema registered for capability env")
	})
}
