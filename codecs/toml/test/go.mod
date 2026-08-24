module github.com/hedzr/store/codecs/toml/test

go 1.26

replace github.com/hedzr/store => ../../..

replace github.com/hedzr/store/providers/file => ../../../providers/file

replace github.com/hedzr/store/codecs/toml => ../

require (
	github.com/hedzr/store v1.4.3
	github.com/hedzr/store/codecs/toml v1.4.3
	github.com/hedzr/store/providers/file v1.4.3
	github.com/stretchr/testify v1.12.1
)

require (
	github.com/fsnotify/fsnotify v1.10.1 // indirect
	github.com/hedzr/evendeep v1.4.3 // indirect
	github.com/hedzr/is v0.9.5 // indirect
	github.com/hedzr/logg v0.9.5 // indirect
	github.com/pelletier/go-toml/v2 v2.4.3 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/net v0.57.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
	golang.org/x/term v0.45.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
