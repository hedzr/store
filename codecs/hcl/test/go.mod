module github.com/hedzr/store/codecs/hcl/test

go 1.21

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/hcl => ../

require (
	github.com/hedzr/store v1.0.8
	github.com/hedzr/store/codecs/hcl v1.0.8
	github.com/hedzr/store/providers/file v1.0.8
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hedzr/evendeep v1.1.11 // indirect
	github.com/hedzr/is v0.5.21 // indirect
	github.com/hedzr/logg v0.5.22 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
