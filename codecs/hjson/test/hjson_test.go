package test_test

import (
	"context"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	hjsonapi "github.com/hjson/hjson-go/v4"
	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hjson"
	"github.com/hedzr/store/providers/file"
)

func TestNew(t *testing.T) {
	s := store.New()
	parser := hjson.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.hjson"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../../../testdata/6.hjson")),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, `r.Header.Get("From")`, s.MustGet("app.hjson.messages.0.placeholders.0.expr").(string))
	assert.Equal(t, `r.Header.Get("User-Agent")`, s.MustGet("app.hjson.messages.1.placeholders.0.expr").(string))
}

func TestHjson(t *testing.T) {
	const filename = "../../../testdata/6.hjson"

	b, err := os.ReadFile(filename)
	if err != nil {
		t.Fatal(err)
	}

	var node hjsonapi.Node
	if err = hjsonapi.UnmarshalWithOptions(b, &node, hjsonapi.DefaultDecoderOptions()); err != nil {
		// panic(err)
		t.Fatal(err)
	}

	t.Logf("%+v", node.Value)

	nd := node.NK("map")
	// nd := node.NK("messages")
	t.Logf("%+v", nd.Value)

	// switch nn := nd.Value.(type) {
	// case string:
	// case *hjsonapi.OrderedMap:
	// case []any:
	// 	for i := 0; i < len(nn); i++ {
	// 		if nc, ok := nn[i].(*hjsonapi.Node); ok {
	// 			t.Logf("  %d. %v", i, nc.Value)
	// 		} else {
	// 			t.Fatalf("  %d. node[%d] is not a *Node", i, i)
	// 		}
	// 	}
	// }

	parser := hjson.New()
	if ce, ok := parser.(store.CodecEx); ok {
		var data map[string]store.ValPkg
		if data, err = ce.UnmarshalEx(b); err != nil {
			t.Fatal(err)
		}
		// t.Logf("valpkg: %+v", data)
		t.Log("-------- ValPkg:\n", spew.Sdump(data))
	}

	// b, err = hjsonapi.Marshal(node)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Logf("hjson content: \n\n%v\n", string(b))
}
