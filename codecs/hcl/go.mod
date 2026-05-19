module github.com/hedzr/store/codecs/hcl

go 1.26

// replace gopkg.in/hedzr/errors.v3 => ../../../../24/libs.errors

// replace github.com/hedzr/evendeep => ../../../libs.diff

// replace github.com/hedzr/go-errors/v2 => ../../../libs.errors

// replace github.com/hedzr/env => ../../../libs.env

// replace github.com/hedzr/is => ../../../libs.is

// replace github.com/hedzr/logg => ../../../libs.logg

replace github.com/hedzr/store => ../..

// replace github.com/hedzr/store/providers/file => ../../providers/file

require (
	github.com/hashicorp/hcl/v2 v2.24.0
	github.com/hedzr/store v1.4.0
	github.com/zclconf/go-cty v1.18.1
)

require (
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/hedzr/evendeep v1.4.0 // indirect
	github.com/hedzr/is v0.9.1 // indirect
	github.com/hedzr/logg v0.9.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	golang.org/x/mod v0.34.0 // indirect
	golang.org/x/net v0.52.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/term v0.41.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	golang.org/x/tools v0.43.0 // indirect
	gopkg.in/hedzr/errors.v3 v3.3.5 // indirect
)
