module github.com/hedzr/store/codecs/all

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/store => ./../..

require (
	github.com/hedzr/store v1.3.45
	github.com/hedzr/store/codecs/gob v1.3.45
	github.com/hedzr/store/codecs/hcl v1.3.45
	github.com/hedzr/store/codecs/hjson v1.3.45
	github.com/hedzr/store/codecs/json v1.3.45
	github.com/hedzr/store/codecs/nestext v1.3.45
	github.com/hedzr/store/codecs/toml v1.3.45
	github.com/hedzr/store/codecs/yaml v1.3.45
)

require (
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hedzr/evendeep v1.3.43 // indirect
	github.com/hedzr/is v0.8.45 // indirect
	github.com/hedzr/logg v0.8.45 // indirect
	github.com/hjson/hjson-go/v4 v4.5.0 // indirect
	github.com/npillmayer/nestext v0.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/term v0.33.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
