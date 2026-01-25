package entities

import "strings"

// RiskLevel represents the security risk level of a capability or grant set.
type RiskLevel int

const (
	RiskLevelLow    RiskLevel = iota // Specific, narrow permissions
	RiskLevelMedium                  // Network access, sensitive reads
	RiskLevelHigh                    // Broad permissions, arbitrary execution
)

// String returns the human-readable name of the risk level.
func (r RiskLevel) String() string {
	switch r {
	case RiskLevelLow:
		return "Low"
	case RiskLevelMedium:
		return "Medium"
	case RiskLevelHigh:
		return "High"
	default:
		return "Unknown"
	}
}

// Dangerous patterns (security domain knowledge)
var (
	// Broad filesystem patterns that grant excessive access
	BroadFilesystemPatterns = []string{
		"**", "/**", "read:**", "write:**", "read:/", "write:/",
		"read:/etc/**", "write:/etc/**",
		"read:/root/**", "write:/root/**",
		"read:/home/**", "write:/home/**",
	}

	// Shell interpreters that allow arbitrary command execution
	DangerousShells = []string{
		"bash", "sh", "zsh", "fish",
		"/bin/bash", "/bin/sh", "/bin/zsh",
	}

	// Script interpreters (matches base + versioned variants)
	DangerousInterpreters = []string{
		"python", "perl", "ruby", "node", "nodejs",
		"php", "lua", "awk", "gawk", "tclsh",
	}

	// Broad environment variable patterns
	BroadEnvPatterns = []string{"*", "AWS_*", "AZURE_*", "GCP_*"}
)

// riskAssessorConfig holds configuration for the RiskAssessor.
type riskAssessorConfig struct {
	customBroadPatterns map[string][]string
}

func defaultRiskAssessorConfig() riskAssessorConfig {
	return riskAssessorConfig{
		customBroadPatterns: make(map[string][]string),
	}
}

// RiskAssessorOption configures a RiskAssessor instance.
type RiskAssessorOption func(*riskAssessorConfig)

// WithCustomBroadPatterns adds additional patterns considered "broad" for a kind.
func WithCustomBroadPatterns(kind string, patterns []string) RiskAssessorOption {
	return func(c *riskAssessorConfig) {
		c.customBroadPatterns[kind] = append(c.customBroadPatterns[kind], patterns...)
	}
}

// RiskAssessor evaluates the security risk of capabilities.
type RiskAssessor struct {
	config riskAssessorConfig
}

// NewRiskAssessor creates a new RiskAssessor with the given options.
func NewRiskAssessor(opts ...RiskAssessorOption) *RiskAssessor {
	cfg := defaultRiskAssessorConfig()
	for _, opt := range opts {
		opt(&cfg)
	}
	return &RiskAssessor{config: cfg}
}

// AssessGrantSet evaluates the overall risk level of a GrantSet.
func (r *RiskAssessor) AssessGrantSet(g *GrantSet) RiskLevel {
	if g == nil {
		return RiskLevelLow
	}

	highest := RiskLevelLow

	// Exec is always High risk
	if g.Exec != nil && len(g.Exec.Commands) > 0 {
		for _, cmd := range g.Exec.Commands {
			if cmd == "*" || cmd == "**" {
				return RiskLevelHigh
			}
			if matchesAny(cmd, DangerousShells) || matchesInterpreter(cmd) {
				return RiskLevelHigh
			}
		}
		// Any exec = medium risk
		if highest < RiskLevelMedium {
			highest = RiskLevelMedium
		}
	}

	// Network
	if g.Network != nil {
		for _, rule := range g.Network.Rules {
			for _, h := range rule.Hosts {
				if h == "*" {
					return RiskLevelHigh
				}
			}
		}
		// Any network access = at least medium
		if len(g.Network.Rules) > 0 && highest < RiskLevelMedium {
			highest = RiskLevelMedium
		}
	}

	// Filesystem
	if g.FS != nil {
		allPatterns := r.getBroadFilesystemPatterns()
		for _, rule := range g.FS.Rules {
			for _, p := range rule.Read {
				// Check for recursive glob patterns
				if strings.Contains(p, "**") || matchesAny(p, allPatterns) {
					return RiskLevelHigh
				}
			}
			for _, p := range rule.Write {
				// Check for recursive glob patterns
				if strings.Contains(p, "**") || matchesAny(p, allPatterns) {
					return RiskLevelHigh
				}
				// Any write is at least Medium
				if highest < RiskLevelMedium {
					highest = RiskLevelMedium
				}
			}
			// Sensitive reads
			for _, p := range rule.Read {
				if strings.HasPrefix(p, "/etc/") && highest < RiskLevelMedium {
					highest = RiskLevelMedium
				}
			}
		}
	}

	// Environment
	if g.Env != nil {
		allPatterns := r.getBroadEnvPatterns()
		for _, v := range g.Env.Variables {
			if matchesAny(v, allPatterns) {
				return RiskLevelHigh
			}
		}
	}

	// KeyValue
	if g.KV != nil {
		for _, rule := range g.KV.Rules {
			if rule.Operation == "write" || rule.Operation == "read-write" {
				if highest < RiskLevelMedium {
					highest = RiskLevelMedium
				}
			}
		}
	}

	return highest
}

// getBroadFilesystemPatterns returns base + custom broad filesystem patterns.
func (r *RiskAssessor) getBroadFilesystemPatterns() []string {
	patterns := make([]string, len(BroadFilesystemPatterns))
	copy(patterns, BroadFilesystemPatterns)
	if custom, ok := r.config.customBroadPatterns["fs"]; ok {
		patterns = append(patterns, custom...)
	}
	return patterns
}

// getBroadEnvPatterns returns base + custom broad env patterns.
func (r *RiskAssessor) getBroadEnvPatterns() []string {
	patterns := make([]string, len(BroadEnvPatterns))
	copy(patterns, BroadEnvPatterns)
	if custom, ok := r.config.customBroadPatterns["env"]; ok {
		patterns = append(patterns, custom...)
	}
	return patterns
}

// DescribeRisks returns a list of human-readable risk descriptions.
func (r *RiskAssessor) DescribeRisks(g *GrantSet) []string {
	var risks []string

	if g == nil {
		return risks
	}

	if g.Exec != nil && len(g.Exec.Commands) > 0 {
		risks = append(risks, "Executes external commands (High Risk)")
	}

	if g.Network != nil {
		for _, rule := range g.Network.Rules {
			for _, h := range rule.Hosts {
				if h == "*" {
					risks = append(risks, "Accesses any network host (High Risk)")
					break
				}
			}
		}
	}

	if g.FS != nil {
		hasHighFS := false
		for _, rule := range g.FS.Rules {
			for _, p := range rule.Read {
				if strings.Contains(p, "**") {
					risks = append(risks, "Recursive read access to filesystem (High Risk)")
					hasHighFS = true
					break
				}
			}
			if hasHighFS {
				break
			}
		}
		if !hasHighFS {
			for _, rule := range g.FS.Rules {
				for _, p := range rule.Write {
					if strings.Contains(p, "**") {
						risks = append(risks, "Recursive write access to filesystem (High Risk)")
						hasHighFS = true
						break
					}
				}
				if hasHighFS {
					break
				}
			}
		}

		// General write check
		for _, rule := range g.FS.Rules {
			if len(rule.Write) > 0 {
				risks = append(risks, "Write access to filesystem")
				break
			}
		}
	}

	if g.Env != nil {
		for _, v := range g.Env.Variables {
			if v == "*" {
				risks = append(risks, "Accesses all environment variables (High Risk)")
				break
			}
		}
	}

	if g.KV != nil {
		for _, rule := range g.KV.Rules {
			if rule.Operation == "write" || rule.Operation == "read-write" {
				risks = append(risks, "Write access to Key-Value store")
				break
			}
		}
	}

	return risks
}

// matchesAny checks if value matches any pattern in the list.
func matchesAny(value string, patterns []string) bool {
	for _, p := range patterns {
		if value == p {
			return true
		}
	}
	return false
}

// matchesInterpreter checks if cmd matches a dangerous interpreter (base or versioned).
func matchesInterpreter(cmd string) bool {
	for _, interp := range DangerousInterpreters {
		if cmd == interp || strings.HasPrefix(cmd, interp) {
			return true
		}
	}
	return false
}
