package entities_test

import (
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestRiskAssessor_AssessGrantSet(t *testing.T) {
	assessor := entities.NewRiskAssessor()

	t.Run("Empty grant set is Low risk", func(t *testing.T) {
		g := &entities.GrantSet{}
		assert.Equal(t, entities.RiskLevelLow, assessor.AssessGrantSet(g))
	})

	t.Run("Specific read access is Low risk", func(t *testing.T) {
		g := &entities.GrantSet{
			FS: &entities.FileSystemCapability{
				Rules: []entities.FileSystemRule{
					{Read: []string{"/tmp/file.txt"}},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelLow, assessor.AssessGrantSet(g))
	})

	t.Run("Filesystem write is Medium risk", func(t *testing.T) {
		g := &entities.GrantSet{
			FS: &entities.FileSystemCapability{
				Rules: []entities.FileSystemRule{
					{Write: []string{"/tmp/file.txt"}},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelMedium, assessor.AssessGrantSet(g))
	})

	t.Run("Recursive filesystem access is High risk", func(t *testing.T) {
		g := &entities.GrantSet{
			FS: &entities.FileSystemCapability{
				Rules: []entities.FileSystemRule{
					{Read: []string{"/data/**"}},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelHigh, assessor.AssessGrantSet(g))
	})

	t.Run("Exec with safe command is Medium risk", func(t *testing.T) {
		g := &entities.GrantSet{
			Exec: &entities.ExecCapability{
				Commands: []string{"/usr/bin/ls"},
			},
		}
		assert.Equal(t, entities.RiskLevelMedium, assessor.AssessGrantSet(g))
	})

	t.Run("Exec with shell is High risk", func(t *testing.T) {
		g := &entities.GrantSet{
			Exec: &entities.ExecCapability{
				Commands: []string{"/bin/bash"},
			},
		}
		assert.Equal(t, entities.RiskLevelHigh, assessor.AssessGrantSet(g))
	})

	t.Run("All Network is High risk", func(t *testing.T) {
		g := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"*"}, Ports: []string{"443"}},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelHigh, assessor.AssessGrantSet(g))
	})

	t.Run("Specific Network is Medium risk", func(t *testing.T) {
		g := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"example.com"}, Ports: []string{"443"}},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelMedium, assessor.AssessGrantSet(g))
	})

	t.Run("All Env is High risk", func(t *testing.T) {
		g := &entities.GrantSet{
			Env: &entities.EnvironmentCapability{
				Variables: []string{"*"},
			},
		}
		assert.Equal(t, entities.RiskLevelHigh, assessor.AssessGrantSet(g))
	})

	t.Run("KV Write is Medium risk", func(t *testing.T) {
		g := &entities.GrantSet{
			KV: &entities.KeyValueCapability{
				Rules: []entities.KeyValueRule{
					{Keys: []string{"config/*"}, Operation: "write"},
				},
			},
		}
		assert.Equal(t, entities.RiskLevelMedium, assessor.AssessGrantSet(g))
	})
}

func TestRiskAssessor_DescribeRisks(t *testing.T) {
	assessor := entities.NewRiskAssessor()

	g := &entities.GrantSet{
		Exec: &entities.ExecCapability{
			Commands: []string{"ls"},
		},
		Network: &entities.NetworkCapability{
			Rules: []entities.NetworkRule{
				{Hosts: []string{"*"}, Ports: []string{"443"}},
			},
		},
		FS: &entities.FileSystemCapability{
			Rules: []entities.FileSystemRule{
				{Write: []string{"/tmp/**"}},
			},
		},
	}

	risks := assessor.DescribeRisks(g)
	assert.Contains(t, risks, "Executes external commands (High Risk)")
	assert.Contains(t, risks, "Accesses any network host (High Risk)")
	assert.Contains(t, risks, "Recursive write access to filesystem (High Risk)")
	assert.Contains(t, risks, "Write access to filesystem")
}

func TestRiskAssessor_WithCustomBroadPatterns(t *testing.T) {
	// Test that custom broad patterns work
	assessor := entities.NewRiskAssessor(
		entities.WithCustomBroadPatterns("fs", []string{"/custom/**"}),
	)

	g := &entities.GrantSet{
		FS: &entities.FileSystemCapability{
			Rules: []entities.FileSystemRule{
				{Read: []string{"/custom/**"}},
			},
		},
	}

	assert.Equal(t, entities.RiskLevelHigh, assessor.AssessGrantSet(g))
}

func TestRiskLevel_String(t *testing.T) {
	assert.Equal(t, "Low", entities.RiskLevelLow.String())
	assert.Equal(t, "Medium", entities.RiskLevelMedium.String())
	assert.Equal(t, "High", entities.RiskLevelHigh.String())
}
