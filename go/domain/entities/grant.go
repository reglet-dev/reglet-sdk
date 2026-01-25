package entities

// GrantSet is a structured collection of rules representing all capabilities granted to a plugin.
type GrantSet struct {
	Network *NetworkCapability     `json:"network,omitempty" yaml:"network,omitempty"`
	FS      *FileSystemCapability  `json:"fs,omitempty" yaml:"fs,omitempty"`
	Env     *EnvironmentCapability `json:"env,omitempty" yaml:"env,omitempty"`
	Exec    *ExecCapability        `json:"exec,omitempty" yaml:"exec,omitempty"`
	KV      *KeyValueCapability    `json:"kv,omitempty" yaml:"kv,omitempty"`
}

// IsEmpty returns true if no capabilities are present.
func (g *GrantSet) IsEmpty() bool {
	if g == nil {
		return true
	}
	if g.Network != nil && len(g.Network.Rules) > 0 {
		return false
	}
	if g.FS != nil && len(g.FS.Rules) > 0 {
		return false
	}
	if g.Env != nil && len(g.Env.Variables) > 0 {
		return false
	}
	if g.Exec != nil && len(g.Exec.Commands) > 0 {
		return false
	}
	if g.KV != nil && len(g.KV.Rules) > 0 {
		return false
	}
	return true
}

// Merge unions two grant sets.
func (g *GrantSet) Merge(other *GrantSet) {
	if other == nil {
		return
	}

	// Merge Network
	if other.Network != nil && len(other.Network.Rules) > 0 {
		if g.Network == nil {
			g.Network = &NetworkCapability{}
		}
		g.Network.Rules = append(g.Network.Rules, other.Network.Rules...)
	}

	// Merge FS
	if other.FS != nil && len(other.FS.Rules) > 0 {
		if g.FS == nil {
			g.FS = &FileSystemCapability{}
		}
		g.FS.Rules = append(g.FS.Rules, other.FS.Rules...)
	}

	// Merge Env
	if other.Env != nil && len(other.Env.Variables) > 0 {
		if g.Env == nil {
			g.Env = &EnvironmentCapability{}
		}
		g.Env.Variables = append(g.Env.Variables, other.Env.Variables...)
	}

	// Merge Exec
	if other.Exec != nil && len(other.Exec.Commands) > 0 {
		if g.Exec == nil {
			g.Exec = &ExecCapability{}
		}
		g.Exec.Commands = append(g.Exec.Commands, other.Exec.Commands...)
	}

	// Merge KV
	if other.KV != nil && len(other.KV.Rules) > 0 {
		if g.KV == nil {
			g.KV = &KeyValueCapability{}
		}
		g.KV.Rules = append(g.KV.Rules, other.KV.Rules...)
	}
}
