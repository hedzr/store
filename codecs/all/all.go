// Package all imports all known codecs
package all

import (
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/gob"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/codecs/hjson"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/codecs/nestext"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/codecs/yaml"
)

var suffixCodecMap = map[string]func() store.Codec{
	"toml":       func() store.Codec { return toml.New() },
	"yaml":       func() store.Codec { return yaml.New() },
	"yml":        func() store.Codec { return yaml.New() },
	"gob":        func() store.Codec { return gob.New() },
	"json":       func() store.Codec { return json.New() },
	"hjson":      func() store.Codec { return hjson.New() },
	"tf":         func() store.Codec { return hcl.New() },
	"hcl":        func() store.Codec { return hcl.New() },
	"nestedtext": func() store.Codec { return nestext.New() },
	"txt":        func() store.Codec { return nestext.New() },
	"conf":       func() store.Codec { return nestext.New() },
	"":           func() store.Codec { return nestext.New() },
}

func Register(ext string, getter func() store.Codec) {
	suffixCodecMap[ext] = getter
}

func Deregister(ext string) {
	delete(suffixCodecMap, ext)
}

func ExtCodecMap() map[string]func() store.Codec { return suffixCodecMap }

func Codec(ext string) (getter func() store.Codec, exists bool) {
	getter, exists = suffixCodecMap[ext]
	return
}

func MustCodec(ext string) (getter func() store.Codec) {
	return suffixCodecMap[ext]
}
