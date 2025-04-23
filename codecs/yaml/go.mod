module github.com/hedzr/store/codecs/yaml

go 1.23.0

toolchain go1.23.3

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

//replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hedzr/store v1.3.15
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/hedzr/evendeep v1.3.15 // indirect
	github.com/hedzr/is v0.7.15 // indirect
	github.com/hedzr/logg v0.8.15 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/term v0.31.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
