module github.com/hedzr/store

go 1.25.0

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

// replace github.com/hedzr/go-errors/v2 => ../libs.errors

// replace github.com/hedzr/evendeep => ../libs.diff

// replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/go-utils/v2 => ./

require (
	github.com/hedzr/evendeep v1.4.0
	github.com/hedzr/is v0.9.0
	github.com/hedzr/logg v0.9.0
	gopkg.in/hedzr/errors.v3 v3.3.5
)

require (
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
)
