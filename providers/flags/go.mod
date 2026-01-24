module github.com/hedzr/store/providers/flags

go 1.24.0

toolchain go1.24.5

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/codecs/json => ../../codecs/json

// replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

require github.com/hedzr/store v1.3.67

require (
	github.com/hedzr/evendeep v1.3.67 // indirect
	github.com/hedzr/is v0.8.67 // indirect
	github.com/hedzr/logg v0.8.67 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/term v0.39.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
