package tests

import (
	"context"
	fs2 "io/fs"
	"os"
	"testing"
	"testing/fstest"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/json"
	"github.com/hedzr/store/providers/fs"
)

func TestStore_fs_Load(t *testing.T) {
	s := newBasicStore()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.maps"),
		store.WithCodec(json.New()),
		store.WithProvider(fs.New(newFs(), "4.json")),
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	// assert.EqualTrue(t, s.MustGet("app.maps.cool.station.8s") == true,
	//    `expecting store.Get("app.maps.cool.station.8s") return 'true'`)
	// assert.EqualTrue(t, s.MustGet("app.maps.cool.station.flush.interval") == 5*time.Hour,
	//    `expecting store.Get("app.maps.cool.station.flush.interval") return '5h0m0s'`)
}

func newFs() fs2.FS {
	// root := "/usr/local/opt/go/bin"
	// fileSystem := os.DirFS(root)
	// return fileSystem

	testFS := os.DirFS("../testdata")
	err := fstest.TestFS(testFS, "4.json")
	if err != nil {
		return nil
	}
	return testFS
}
