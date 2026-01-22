module examples/host-runtime

go 1.25.5

replace github.com/reglet-dev/reglet-sdk => ../../go

replace github.com/reglet-dev/reglet-sdk/go => ../../go

require (
	github.com/reglet-dev/reglet-sdk/go v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tetratelabs/wazero v1.11.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
