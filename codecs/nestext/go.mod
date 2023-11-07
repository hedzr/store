module github.com/hedzr/store/codecs/nestext

go 1.21

replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

replace github.com/hedzr/go-diff/v2 => ../../../libs.diff

//replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

replace github.com/hedzr/env => ../../../libs.env

replace github.com/hedzr/is => ../../../libs.is

replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hedzr/env v0.0.0-00010101000000-000000000000
	github.com/hedzr/store v0.0.0-00010101000000-000000000000
	github.com/hedzr/store/providers/file v0.0.0-00010101000000-000000000000
	github.com/npillmayer/nestext v0.1.3
)

require (
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/hedzr/go-diff/v2 v2.0.0-00010101000000-000000000000 // indirect
	github.com/hedzr/is v0.5.1 // indirect
	github.com/hedzr/logg v0.0.0-00010101000000-000000000000 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.3 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.0 // indirect
)
