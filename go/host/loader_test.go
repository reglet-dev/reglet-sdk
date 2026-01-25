package host_test

import (
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/reglet-dev/reglet-sdk/go/host"
	"github.com/reglet-dev/reglet-sdk/go/host/registry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// LoaderIntegrationSuite tests the Loader with full integration.
type LoaderIntegrationSuite struct {
	suite.Suite
	registry *registry.Registry
	loader   *host.Loader
}

func (s *LoaderIntegrationSuite) SetupTest() {
	// Create and configure registry
	reg := registry.NewRegistry(registry.WithStrictMode(false))
	err := reg.Register("network", entities.NetworkCapability{})
	s.Require().NoError(err)
	err = reg.Register("fs", entities.FileSystemCapability{})
	s.Require().NoError(err)
	err = reg.Register("env", entities.EnvironmentCapability{})
	s.Require().NoError(err)
	err = reg.Register("exec", entities.ExecCapability{})
	s.Require().NoError(err)
	err = reg.Register("kv", entities.KeyValueCapability{})
	s.Require().NoError(err)

	s.registry = reg.(*registry.Registry)
	s.loader = host.NewLoader(host.WithRegistry(reg))
}

func (s *LoaderIntegrationSuite) TestValidManifest() {
	yaml := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: ["example.com"]
        ports: ["80", "443"]
  fs:
    rules:
      - read: ["/data/**"]
        write: ["/tmp/*"]
`
	manifest, err := s.loader.LoadManifest([]byte(yaml), nil)
	s.Require().NoError(err)
	s.Equal("test-plugin", manifest.Name)
	s.NotNil(manifest.Capabilities.Network)
	s.Len(manifest.Capabilities.Network.Rules, 1)
	s.Equal([]string{"example.com"}, manifest.Capabilities.Network.Rules[0].Hosts)
	s.NotNil(manifest.Capabilities.FS)
	s.Len(manifest.Capabilities.FS.Rules, 1)
}

func (s *LoaderIntegrationSuite) TestManifestWithMultipleRules() {
	yaml := `
name: "multi-rule-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: ["api.internal"]
        ports: ["80"]
      - hosts: ["*.external.com"]
        ports: ["443"]
  kv:
    rules:
      - keys: ["config/*"]
        op: "read"
      - keys: ["cache/*"]
        op: "read-write"
`
	manifest, err := s.loader.LoadManifest([]byte(yaml), nil)
	s.Require().NoError(err)
	s.Len(manifest.Capabilities.Network.Rules, 2)
	s.Equal("api.internal", manifest.Capabilities.Network.Rules[0].Hosts[0])
	s.Equal("*.external.com", manifest.Capabilities.Network.Rules[1].Hosts[0])
	s.Len(manifest.Capabilities.KV.Rules, 2)
	s.Equal("read", manifest.Capabilities.KV.Rules[0].Operation)
	s.Equal("read-write", manifest.Capabilities.KV.Rules[1].Operation)
}

func (s *LoaderIntegrationSuite) TestInvalidYAML() {
	yaml := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: 123  # Should be string array
`
	_, err := s.loader.LoadManifest([]byte(yaml), nil)
	s.Require().Error(err)
	s.Contains(err.Error(), "cannot unmarshal")
}

func (s *LoaderIntegrationSuite) TestMissingSchemaRegistration() {
	// Create loader with empty registry
	emptyReg := registry.NewRegistry()
	loaderEmpty := host.NewLoader(host.WithRegistry(emptyReg))

	yaml := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: ["example.com"]
        ports: ["443"]
`
	_, err := loaderEmpty.LoadManifest([]byte(yaml), nil)
	s.Require().Error(err)
	s.Contains(err.Error(), "no schema registered for capability network")
}

func (s *LoaderIntegrationSuite) TestEnvCapability() {
	yaml := `
name: "env-plugin"
version: "1.0.0"
capabilities:
  env:
    vars: ["APP_*", "DEBUG"]
`
	manifest, err := s.loader.LoadManifest([]byte(yaml), nil)
	s.Require().NoError(err)
	s.NotNil(manifest.Capabilities.Env)
	s.ElementsMatch([]string{"APP_*", "DEBUG"}, manifest.Capabilities.Env.Variables)
}

func (s *LoaderIntegrationSuite) TestExecCapability() {
	yaml := `
name: "exec-plugin"
version: "1.0.0"
capabilities:
  exec:
    commands: ["/usr/bin/ls", "/usr/bin/cat"]
`
	manifest, err := s.loader.LoadManifest([]byte(yaml), nil)
	s.Require().NoError(err)
	s.NotNil(manifest.Capabilities.Exec)
	s.ElementsMatch([]string{"/usr/bin/ls", "/usr/bin/cat"}, manifest.Capabilities.Exec.Commands)
}

func TestLoaderIntegrationSuite(t *testing.T) {
	suite.Run(t, new(LoaderIntegrationSuite))
}

// Additional standalone tests for backwards compatibility
func TestLoader_Integration(t *testing.T) {
	// 1. Setup Registry
	reg := registry.NewRegistry(registry.WithStrictMode(false))
	err := reg.Register("network", entities.NetworkCapability{})
	require.NoError(t, err)
	err = reg.Register("fs", entities.FileSystemCapability{})
	require.NoError(t, err)

	// 2. Setup Loader
	loader := host.NewLoader(
		host.WithRegistry(reg),
	)

	t.Run("Valid Manifest", func(t *testing.T) {
		yaml := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: ["example.com"]
        ports: ["80"]
`
		manifest, err := loader.LoadManifest([]byte(yaml), nil)
		require.NoError(t, err)
		assert.Equal(t, "test-plugin", manifest.Name)
		assert.NotNil(t, manifest.Capabilities.Network)
		assert.Len(t, manifest.Capabilities.Network.Rules, 1)
		assert.Equal(t, []string{"example.com"}, manifest.Capabilities.Network.Rules[0].Hosts)
	})

	t.Run("Invalid Schema", func(t *testing.T) {
		yaml := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: 123  # Should be string array
`
		_, err := loader.LoadManifest([]byte(yaml), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot unmarshal")
	})

	t.Run("Missing Capability Registration", func(t *testing.T) {
		emptyReg := registry.NewRegistry()
		loaderEmpty := host.NewLoader(host.WithRegistry(emptyReg))

		yaml2 := `
name: "test-plugin"
version: "1.0.0"
capabilities:
  network:
    rules:
      - hosts: ["example.com"]
        ports: ["443"]
`
		_, err := loaderEmpty.LoadManifest([]byte(yaml2), nil)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no schema registered for capability network")
	})
}
