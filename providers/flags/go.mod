module github.com/hedzr/store/providers/flags

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/codecs/json => ../../codecs/json

// replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

require github.com/hedzr/store v1.3.60

require (
	github.com/hedzr/evendeep v1.3.60 // indirect
	github.com/hedzr/is v0.8.60 // indirect
	github.com/hedzr/logg v0.8.60 // indirect
	golang.org/x/net v0.44.0 // indirect
	golang.org/x/sys v0.36.0 // indirect
	golang.org/x/term v0.35.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
