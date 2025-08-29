module github.com/hedzr/store

go 1.24.0

toolchain go1.24.5

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/go-errors/v2 => ../libs.errors

// replace github.com/hedzr/evendeep => ../libs.diff

// replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/go-utils/v2 => ./

require (
	github.com/hedzr/evendeep v1.3.53
	github.com/hedzr/is v0.8.53
	github.com/hedzr/logg v0.8.53
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	golang.org/x/net v0.43.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/term v0.34.0 // indirect
)
