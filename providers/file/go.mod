module github.com/hedzr/store/providers/file

go 1.22.7

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/evendeep => ../. ./../libs.diff

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/codecs/json => ../../codecs/json

// replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/hedzr/store v1.0.15
)

require (
	github.com/hedzr/evendeep v1.1.15 // indirect
	github.com/hedzr/is v0.5.23 // indirect
	github.com/hedzr/logg v0.5.23 // indirect
	golang.org/x/crypto v0.27.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/term v0.24.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
)
