module github.com/hedzr/store/codecs/json/test

go 1.22.7

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/json => ../

require (
	github.com/hedzr/store v1.2.11
	github.com/hedzr/store/codecs/json v1.2.11
	github.com/hedzr/store/providers/file v1.2.11
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/hedzr/evendeep v1.2.12 // indirect
	github.com/hedzr/is v0.6.8 // indirect
	github.com/hedzr/logg v0.7.20 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
