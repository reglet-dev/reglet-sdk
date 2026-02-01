package entities

// Capability represents a permission or capability required by an SDK operation.
// Capabilities follow the format "category:resource" (e.g., "network:outbound:443").
type Capability struct {
	// Category is the capability category (e.g., "network", "fs", "exec").
	Category string `json:"kind" yaml:"kind"`

	// Resource is the specific resource within the category.
	Resource string `json:"pattern" yaml:"pattern"`

	// Action is the permitted action (e.g., "read", "write", "connect").
	Action string `json:"action,omitempty" yaml:"action,omitempty"`
}

// NewCapability creates a new Capability with the given category and resource.
func NewCapability(category, resource string) Capability {
	return Capability{
		Category: category,
		Resource: resource,
	}
}

// WithAction returns a copy of the Capability with the action set.
func (c Capability) WithAction(action string) Capability {
	c.Action = action
	return c
}

// String returns the capability in "category:resource" format.
func (c Capability) String() string {
	if c.Action != "" {
		return c.Category + ":" + c.Resource + ":" + c.Action
	}
	return c.Category + ":" + c.Resource
}

// Common capabilities
var (
	// CapabilityHTTP represents a general HTTP network capability.
	// In practice, plugins should request specific domains, but this is a broad capability.
	CapabilityHTTP = NewCapability("network", "http")

	// CapabilityDNS represents a DNS network capability.
	CapabilityDNS = NewCapability("network", "dns")

	// CapabilityFile represents a general file system capability.
	CapabilityFile = NewCapability("fs", "**")

	// CapabilityTCP represents a TCP network capability.
	CapabilityTCP = NewCapability("network", "tcp")

	// CapabilitySMTP represents an SMTP network capability.
	CapabilitySMTP = NewCapability("network", "smtp")

	// CapabilityExec represents a general execution capability.
	CapabilityExec = NewCapability("exec", "**")
)
