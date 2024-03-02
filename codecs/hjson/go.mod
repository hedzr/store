module github.com/hedzr/store/codecs/hjson

go 1.21

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

//replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hedzr/store v0.0.0-00010101000000-000000000000
	github.com/hjson/hjson-go/v4 v4.4.0
)

require (
	github.com/hedzr/evendeep v1.1.3 // indirect
	github.com/hedzr/is v0.5.15 // indirect
	github.com/hedzr/logg v0.5.11 // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
