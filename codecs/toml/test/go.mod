module github.com/hedzr/store/codecs/toml/test

go 1.21.0

toolchain go1.22.7

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/toml => ../

require (
	github.com/hedzr/store v1.0.15
	github.com/hedzr/store/codecs/toml v1.0.15
	github.com/hedzr/store/providers/file v1.0.15
	github.com/stretchr/testify v1.9.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hedzr/evendeep v1.1.15 // indirect
	github.com/hedzr/is v0.5.23 // indirect
	github.com/hedzr/logg v0.5.23 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/term v0.24.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
