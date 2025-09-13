module github.com/hedzr/store/codecs/hcl/test

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/hcl => ../

require (
	github.com/hedzr/store v1.3.60
	github.com/hedzr/store/codecs/hcl v1.3.60
	github.com/hedzr/store/providers/file v1.3.60
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.24.0 // indirect
	github.com/hedzr/evendeep v1.3.60 // indirect
	github.com/hedzr/is v0.8.60 // indirect
	github.com/hedzr/logg v0.8.60 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zclconf/go-cty v1.17.0 // indirect
	golang.org/x/mod v0.28.0 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
	golang.org/x/text v0.29.0 // indirect
	golang.org/x/tools v0.37.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
