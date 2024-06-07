module github.com/hedzr/store/examples/testconsul

go 1.21

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

//replace github.com/hedzr/store => ../../
//
//replace github.com/hedzr/store/codecs/json => ../../codecs/json
//
//replace github.com/hedzr/store/codecs/yaml => ../../codecs/yaml
//
//replace github.com/hedzr/store/providers/consul => ../../providers/consul
//
//replace github.com/hedzr/store/providers/env => ../../providers/env
//
//replace github.com/hedzr/store/providers/etcd => ../../providers/etcd
//
//replace github.com/hedzr/store/providers/file => ../../providers/file
//
//replace github.com/hedzr/store/providers/fs => ../../providers/fs
//
//replace github.com/hedzr/store/providers/maps => ../../providers/maps

require (
	github.com/hashicorp/consul/api v1.29.1
	github.com/hedzr/logg v0.5.21
	github.com/hedzr/store v1.0.7
	github.com/hedzr/store/providers/consul v1.0.7
)

require (
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.6.3 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/golang-lru v1.0.2 // indirect
	github.com/hashicorp/serf v0.10.1 // indirect
	github.com/hedzr/evendeep v1.1.11 // indirect
	github.com/hedzr/is v0.5.20 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/exp v0.0.0-20240604190554-fc45aab8b7f8 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/term v0.21.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
)
