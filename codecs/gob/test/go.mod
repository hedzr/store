module github.com/hedzr/store/codecs/gob/test

go 1.25.0

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/gob => ../

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc
	github.com/hedzr/store v1.4.0
	github.com/hedzr/store/codecs/gob v1.4.0
	github.com/hedzr/store/providers/file v1.4.0
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/hedzr/evendeep v1.4.0 // indirect
	github.com/hedzr/is v0.9.1 // indirect
	github.com/hedzr/logg v0.9.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/term v0.41.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
