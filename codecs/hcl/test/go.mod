module github.com/hedzr/store/codecs/hcl/test

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/hcl => ../

require (
	github.com/hedzr/store v1.3.20
	github.com/hedzr/store/codecs/hcl v1.3.20
	github.com/hedzr/store/providers/file v1.3.20
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hedzr/evendeep v1.3.20 // indirect
	github.com/hedzr/is v0.7.20 // indirect
	github.com/hedzr/logg v0.8.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
