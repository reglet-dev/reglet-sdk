package entities

// Capability represents a permission or capability required by an SDK operation.
// Capabilities follow the format "category:resource" (e.g., "network:outbound:443").
type Capability struct {
	// Category is the capability category (e.g., "network", "fs", "exec").
	Category string `json:"kind"`

	// Resource is the specific resource within the category.
	Resource string `json:"pattern"`

	// Action is the permitted action (e.g., "read", "write", "connect").
	Action string `json:"action,omitempty"`
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

// NetworkCapability creates a network capability.
func NetworkCapability(resource string) Capability {
	return NewCapability("network", resource)
}

// FileSystemCapability creates a filesystem capability.
func FileSystemCapability(resource string) Capability {
	return NewCapability("fs", resource)
}

// ExecCapability creates an exec capability.
func ExecCapability(command string) Capability {
	return NewCapability("exec", command)
}
