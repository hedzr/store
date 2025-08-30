module github.com/hedzr/store/codecs/hjson/test

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/hjson => ../

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hedzr/store v1.3.55
	github.com/hedzr/store/codecs/hjson v1.3.55
	github.com/hedzr/store/providers/file v1.3.55
	github.com/hjson/hjson-go/v4 v4.5.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/hedzr/evendeep v1.3.55 // indirect
	github.com/hedzr/is v0.8.55 // indirect
	github.com/hedzr/logg v0.8.55 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
