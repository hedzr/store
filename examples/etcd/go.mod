module github.com/hedzr/store/examples/testetcd

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
	github.com/hedzr/logg v0.5.20
	github.com/hedzr/store v1.0.5
	github.com/hedzr/store/providers/etcd v1.0.5
)

require (
	github.com/coreos/go-semver v0.3.1 // indirect
	github.com/coreos/go-systemd/v22 v22.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/hedzr/evendeep v1.1.10 // indirect
	github.com/hedzr/is v0.5.19 // indirect
	go.etcd.io/etcd/api/v3 v3.5.13 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.13 // indirect
	go.etcd.io/etcd/client/v3 v3.5.13 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/term v0.19.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240415180920-8c6c420018be // indirect
	google.golang.org/grpc v1.63.2 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.2 // indirect
)
