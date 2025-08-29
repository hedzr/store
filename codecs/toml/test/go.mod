module github.com/hedzr/store/codecs/toml/test

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/toml => ../

require (
	github.com/hedzr/store v1.3.53
	github.com/hedzr/store/codecs/toml v1.3.53
	github.com/hedzr/store/providers/file v1.3.53
	github.com/stretchr/testify v1.11.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.9.0 // indirect
	github.com/hedzr/evendeep v1.3.53 // indirect
	github.com/hedzr/is v0.8.53 // indirect
	github.com/hedzr/logg v0.8.53 // indirect
	github.com/pelletier/go-toml/v2 v2.2.4 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
