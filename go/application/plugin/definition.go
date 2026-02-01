package plugin

import (
	"encoding/json"
	"sync"

	"github.com/reglet-dev/reglet-sdk/go/application/schema"
	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
)

// PluginDef defines plugin identity and configuration.
type PluginDef struct {
	Name         string
	Version      string
	Description  string
	Config       interface{} // Struct for schema generation
	Capabilities []entities.Capability
}

// PluginDefinition holds the parsed plugin definition and registered services.
type PluginDefinition struct {
	def          PluginDef
	configSchema json.RawMessage
	services     map[string]*serviceEntry
	mu           sync.RWMutex
}

// serviceEntry holds a registered service.
type serviceEntry struct {
	name        string
	description string
	operations  map[string]*operationEntry
}

// operationEntry holds a registered operation.
type operationEntry struct {
	name        string
	description string
	handler     HandlerFunc
}

// DefinePlugin creates a new plugin definition.
// Call this once at package level in your plugin.
func DefinePlugin(def PluginDef) *PluginDefinition {
	var configSchema []byte
	var err error
	if def.Config != nil {
		configSchema, err = schema.GenerateSchema(def.Config)
		if err != nil {
			panic("failed to generate config schema: " + err.Error())
		}
	} else {
		// Empty schema or default
		configSchema = []byte("{}")
	}

	return &PluginDefinition{
		def:          def,
		configSchema: configSchema,
		services:     make(map[string]*serviceEntry),
	}
}

// Manifest returns the complete plugin manifest.
func (p *PluginDefinition) Manifest() *entities.Manifest {
	p.mu.RLock()
	defer p.mu.RUnlock()

	services := make(map[string]entities.ServiceManifest)
	for name, svc := range p.services {
		ops := make([]entities.OperationManifest, 0, len(svc.operations))
		for _, op := range svc.operations {
			ops = append(ops, entities.OperationManifest{
				Name:        op.name,
				Description: op.description,
			})
		}
		services[name] = entities.ServiceManifest{
			Name:        svc.name,
			Description: svc.description,
			Operations:  ops,
		}
	}

	return &entities.Manifest{
		Name:         p.def.Name,
		Version:      p.def.Version,
		Description:  p.def.Description,
		SDKVersion:   Version, // From sdk version.go
		Capabilities: p.def.Capabilities,
		ConfigSchema: p.configSchema,
		Services:     services,
	}
}

// RegisterHandler registers a handler for a service/operation.
// Called internally by RegisterService.
func (p *PluginDefinition) RegisterHandler(serviceName, serviceDesc, opName, opDesc string, handler HandlerFunc) {
	p.mu.Lock()
	defer p.mu.Unlock()

	svc, ok := p.services[serviceName]
	if !ok {
		svc = &serviceEntry{
			name:        serviceName,
			description: serviceDesc,
			operations:  make(map[string]*operationEntry),
		}
		p.services[serviceName] = svc
	}

	svc.operations[opName] = &operationEntry{
		name:        opName,
		description: opDesc,
		handler:     handler,
	}
}

// GetHandler returns a handler for the given service/operation.
func (p *PluginDefinition) GetHandler(serviceName, opName string) (HandlerFunc, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	svc, ok := p.services[serviceName]
	if !ok {
		return nil, false
	}

	op, ok := svc.operations[opName]
	if !ok {
		return nil, false
	}

	return op.handler, true
}
