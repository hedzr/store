module github.com/hedzr/store/examples/testflags

go 1.22.7

// toolchain go1.22.0

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../../

replace github.com/hedzr/store/codecs/json => ../../codecs/json

replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml

replace github.com/hedzr/store/providers/consul => ../../providers/consul

replace github.com/hedzr/store/providers/env => ../../providers/env

replace github.com/hedzr/store/providers/etcd => ../../providers/etcd

replace github.com/hedzr/store/providers/file => ../../providers/file

replace github.com/hedzr/store/providers/flags => ../../providers/flags

replace github.com/hedzr/store/providers/fs => ../../providers/fs

replace github.com/hedzr/store/providers/maps => ../../providers/maps

require (
	github.com/hedzr/logg v0.7.17
	github.com/hedzr/store v1.2.8
	github.com/hedzr/store/providers/flags v1.2.8
)

require (
	github.com/hedzr/evendeep v1.2.9 // indirect
	github.com/hedzr/is v0.6.5 // indirect
	golang.org/x/crypto v0.32.0 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
	golang.org/x/term v0.28.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
