module github.com/hedzr/store/codecs/gob/test

go 1.22.7

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/gob => ../

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hedzr/store v1.1.0
	github.com/hedzr/store/codecs/gob v1.1.0
	github.com/hedzr/store/providers/file v1.1.0
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/hedzr/evendeep v1.2.3 // indirect
	github.com/hedzr/is v0.6.0 // indirect
	github.com/hedzr/logg v0.7.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.28.0 // indirect
	golang.org/x/net v0.30.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/term v0.25.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
