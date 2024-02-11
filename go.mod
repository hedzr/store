module github.com/hedzr/store

go 1.21

// replace gopkg.in/hedzr/errors.v3 => ../../24/libs.errors

//replace github.com/hedzr/go-errors/v2 => ../libs.errors

//replace github.com/hedzr/evendeep => ../libs.diff

//replace github.com/hedzr/env => ../libs.env

// replace github.com/hedzr/is => ../libs.is

// replace github.com/hedzr/logg => ../libs.logg

// replace github.com/hedzr/go-utils/v2 => ./

require (
	github.com/hedzr/evendeep v1.0.0
	github.com/hedzr/is v0.5.13
	github.com/hedzr/logg v0.5.7
	gopkg.in/hedzr/errors.v3 v3.3.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	golang.org/x/crypto v0.19.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
	golang.org/x/term v0.17.0 // indirect
)
