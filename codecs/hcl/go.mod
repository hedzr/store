module github.com/hedzr/store/codecs/hcl

go 1.21

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hashicorp/hcl v1.0.0
	github.com/hedzr/store v1.0.8
)

require (
	github.com/hedzr/evendeep v1.1.11 // indirect
	github.com/hedzr/is v0.5.21 // indirect
	github.com/hedzr/logg v0.5.22 // indirect
	golang.org/x/crypto v0.25.0 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
	golang.org/x/term v0.22.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
)
