package test_test

import (
	"context"
	"encoding/gob"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	gobcodec "github.com/hedzr/store/codecs/gob"
	"github.com/hedzr/store/providers/file"
)

const filename = "../../../testdata/101.gob"

func TestNew(t *testing.T) {
	s := store.New()
	parser := gobcodec.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.gob"),
		store.WithCodec(parser),
		store.WithProvider(file.New(filename)),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, []string{"a", "1", "false"},
		s.MustGet("app.gob.app.logging.words"),
	)
	//	 `expecting store.Get("app.hcl.service.1.http.0.web_proxy.0.process.0.main.0.command.0") return '/usr/local/bin/awesome-app'`)
}

func TestSaveGob(t *testing.T) {
	// if _, yes, _ := Exists(filename); yes {
	// 	return
	// }

	conf := newBasicStore()
	m := conf.MustM("")
	t.Logf("conf -> map: \n%v", spew.Sdump(m))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	enc := gob.NewEncoder(f)
	err = enc.Encode(m)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("encoding to 101.gob OK.")
}

func TestLoadGob(t *testing.T) {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0o644)
	if err != nil {
		t.Fatal(err)
	}
	dec := gob.NewDecoder(f)

	var m map[string]any
	err = dec.Decode(&m)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("decoding to 101.gob OK. map is:\n%v", spew.Sdump(m))
}

func newBasicStore(opts ...store.Opt) store.Store {
	conf := store.New(opts...)
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	// conf.Set("app.logging.rotate", 6)
	// conf.Set("app.logging.words", []string{"a", "1", "false"})

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	return conf
}

// Exists returns the existence of an directory or file.
// See the short version FileExists.
func Exists(filepath string) (os.FileInfo, bool, error) {
	if fi, err := os.Stat(os.ExpandEnv(filepath)); err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, true, err
	} else {
		return fi, true, nil
	}
}
