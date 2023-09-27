module github.com/mabels/wueste

go 1.19

require (
	github.com/google/uuid v1.3.0
	github.com/iancoleman/orderedmap v0.3.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	github.com/traefik/yaegi v0.15.1
)

replace github.com/iancoleman/orderedmap => github.com/mabels/orderedmap v0.0.0-20230926124100-82392f2f89fe

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
