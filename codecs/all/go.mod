module github.com/hedzr/store/codecs/all

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ./../..

require (
	github.com/hedzr/store v1.3.53
	github.com/hedzr/store/codecs/gob v1.3.53
	github.com/hedzr/store/codecs/hcl v1.3.53
	github.com/hedzr/store/codecs/hjson v1.3.53
	github.com/hedzr/store/codecs/json v1.3.53
	github.com/hedzr/store/codecs/nestext v1.3.53
	github.com/hedzr/store/codecs/toml v1.3.53
	github.com/hedzr/store/codecs/yaml v1.3.53
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hedzr/evendeep v1.3.53 // indirect
	github.com/hedzr/is v0.8.53 // indirect
	github.com/hedzr/logg v0.8.53 // indirect
	github.com/hjson/hjson-go/v4 v4.5.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/npillmayer/nestext v0.1.3 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/zclconf/go-cty v1.16.4 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
