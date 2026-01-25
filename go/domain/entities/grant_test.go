package entities_test

import (
	"testing"

	"github.com/reglet-dev/reglet-sdk/go/domain/entities"
	"github.com/stretchr/testify/assert"
)

func TestGrantSet_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		grantSet *entities.GrantSet
		want     bool
	}{
		{
			name:     "Nil GrantSet",
			grantSet: nil,
			want:     true,
		},
		{
			name:     "Empty GrantSet",
			grantSet: &entities.GrantSet{},
			want:     true,
		},
		{
			name: "GrantSet with empty capabilities",
			grantSet: &entities.GrantSet{
				Network: &entities.NetworkCapability{},
			},
			want: true,
		},
		{
			name: "GrantSet with Network capability",
			grantSet: &entities.GrantSet{
				Network: &entities.NetworkCapability{
					Rules: []entities.NetworkRule{
						{Hosts: []string{"example.com"}, Ports: []string{"80"}},
					},
				},
			},
			want: false,
		},
		{
			name: "GrantSet with FS capability",
			grantSet: &entities.GrantSet{
				FS: &entities.FileSystemCapability{
					Rules: []entities.FileSystemRule{
						{Read: []string{"/data/**"}},
					},
				},
			},
			want: false,
		},
		{
			name: "GrantSet with Env capability",
			grantSet: &entities.GrantSet{
				Env: &entities.EnvironmentCapability{
					Variables: []string{"DEBUG"},
				},
			},
			want: false,
		},
		{
			name: "GrantSet with Exec capability",
			grantSet: &entities.GrantSet{
				Exec: &entities.ExecCapability{
					Commands: []string{"/usr/bin/ls"},
				},
			},
			want: false,
		},
		{
			name: "GrantSet with KV capability",
			grantSet: &entities.GrantSet{
				KV: &entities.KeyValueCapability{
					Rules: []entities.KeyValueRule{
						{Keys: []string{"config/*"}, Operation: "read"},
					},
				},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.grantSet.IsEmpty())
		})
	}
}

func TestGrantSet_Merge(t *testing.T) {
	t.Run("Merge Nil", func(t *testing.T) {
		g := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"h1"}, Ports: []string{"80"}},
				},
			},
		}
		g.Merge(nil)
		assert.Len(t, g.Network.Rules, 1)
		assert.Equal(t, []string{"h1"}, g.Network.Rules[0].Hosts)
	})

	t.Run("Merge Empty", func(t *testing.T) {
		g := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"h1"}, Ports: []string{"80"}},
				},
			},
		}
		g.Merge(&entities.GrantSet{})
		assert.Len(t, g.Network.Rules, 1)
	})

	t.Run("Merge Network Rules", func(t *testing.T) {
		g1 := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"h1"}, Ports: []string{"80"}},
				},
			},
		}
		g2 := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"h2"}, Ports: []string{"443"}},
				},
			},
		}
		g1.Merge(g2)
		assert.Len(t, g1.Network.Rules, 2)
		assert.Equal(t, "h1", g1.Network.Rules[0].Hosts[0])
		assert.Equal(t, "h2", g1.Network.Rules[1].Hosts[0])
	})

	t.Run("Merge FS Rules", func(t *testing.T) {
		g1 := &entities.GrantSet{
			FS: &entities.FileSystemCapability{
				Rules: []entities.FileSystemRule{
					{Read: []string{"/data/**"}},
				},
			},
		}
		g2 := &entities.GrantSet{
			FS: &entities.FileSystemCapability{
				Rules: []entities.FileSystemRule{
					{Write: []string{"/tmp/**"}},
				},
			},
		}
		g1.Merge(g2)
		assert.Len(t, g1.FS.Rules, 2)
	})

	t.Run("Merge Env", func(t *testing.T) {
		g1 := &entities.GrantSet{
			Env: &entities.EnvironmentCapability{
				Variables: []string{"FOO"},
			},
		}
		g2 := &entities.GrantSet{
			Env: &entities.EnvironmentCapability{
				Variables: []string{"BAR"},
			},
		}
		g1.Merge(g2)
		assert.ElementsMatch(t, []string{"FOO", "BAR"}, g1.Env.Variables)
	})

	t.Run("Merge Exec", func(t *testing.T) {
		g1 := &entities.GrantSet{
			Exec: &entities.ExecCapability{
				Commands: []string{"/usr/bin/ls"},
			},
		}
		g2 := &entities.GrantSet{
			Exec: &entities.ExecCapability{
				Commands: []string{"/usr/bin/cat"},
			},
		}
		g1.Merge(g2)
		assert.ElementsMatch(t, []string{"/usr/bin/ls", "/usr/bin/cat"}, g1.Exec.Commands)
	})

	t.Run("Merge KV Rules", func(t *testing.T) {
		g1 := &entities.GrantSet{
			KV: &entities.KeyValueCapability{
				Rules: []entities.KeyValueRule{
					{Keys: []string{"config/*"}, Operation: "read"},
				},
			},
		}
		g2 := &entities.GrantSet{
			KV: &entities.KeyValueCapability{
				Rules: []entities.KeyValueRule{
					{Keys: []string{"cache/*"}, Operation: "write"},
				},
			},
		}
		g1.Merge(g2)
		assert.Len(t, g1.KV.Rules, 2)
	})

	t.Run("Merge into nil capability", func(t *testing.T) {
		g1 := &entities.GrantSet{}
		g2 := &entities.GrantSet{
			Network: &entities.NetworkCapability{
				Rules: []entities.NetworkRule{
					{Hosts: []string{"h2"}, Ports: []string{"443"}},
				},
			},
		}
		g1.Merge(g2)
		assert.NotNil(t, g1.Network)
		assert.Len(t, g1.Network.Rules, 1)
	})
}
