module github.com/hedzr/store/providers/flags

go 1.23.0

toolchain go1.23.3

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/codecs/json => ../../codecs/json

// replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

require github.com/hedzr/store v1.3.13

require (
	github.com/hedzr/evendeep v1.3.13 // indirect
	github.com/hedzr/is v0.7.13 // indirect
	github.com/hedzr/logg v0.8.13 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
