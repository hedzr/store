module github.com/hedzr/store/codecs/toml

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
	github.com/hedzr/store v1.2.11
	github.com/pelletier/go-toml/v2 v2.2.3
)

require (
	github.com/hedzr/evendeep v1.2.12 // indirect
	github.com/hedzr/is v0.6.8 // indirect
	github.com/hedzr/logg v0.7.20 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/term v0.29.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
