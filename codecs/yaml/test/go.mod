module github.com/hedzr/store/codecs/yaml/test

go 1.21

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/yaml => ../

require (
	github.com/hedzr/store v0.0.0-00010101000000-000000000000
	github.com/hedzr/store/codecs/yaml v0.0.0-00010101000000-000000000000
	github.com/hedzr/store/providers/file v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.8.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hedzr/evendeep v1.0.0 // indirect
	github.com/hedzr/is v0.5.13 // indirect
	github.com/hedzr/logg v0.5.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
