module github.com/hedzr/store/codecs/hjson/test

go 1.22.7

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/hjson => ../

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/hedzr/store v1.2.8
	github.com/hedzr/store/codecs/hjson v1.2.8
	github.com/hedzr/store/providers/file v1.2.8
	github.com/hjson/hjson-go/v4 v4.4.0
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/hedzr/evendeep v1.2.10 // indirect
	github.com/hedzr/is v0.6.6 // indirect
	github.com/hedzr/logg v0.7.18 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
