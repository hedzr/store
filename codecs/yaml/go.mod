module github.com/hedzr/store/codecs/yaml

go 1.22.7

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

//replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hedzr/store v1.2.2
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/hedzr/evendeep v1.2.7 // indirect
	github.com/hedzr/is v0.6.2 // indirect
	github.com/hedzr/logg v0.7.15 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/net v0.31.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
	golang.org/x/term v0.26.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
