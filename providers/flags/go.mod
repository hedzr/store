module github.com/hedzr/store/providers/flags

go 1.25.0

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/codecs/json => ../../codecs/json

// replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

require github.com/hedzr/store v1.4.0

require (
	github.com/hedzr/evendeep v1.4.0 // indirect
	github.com/hedzr/is v0.9.0 // indirect
	github.com/hedzr/logg v0.9.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
