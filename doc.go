// Package store provides an extensible, high-performance configuration management
// library, specially optimized for hierarchical data.
//
// The [Store] interface gives these APIs.
//
// The `hedzr/store` (https://github.com/hedzr/store) accesses tree data with a dotted key path, which means you
// may point to a specified tree node and access it, monitor it or
// remove it.
//
//	conf := store.New()
//	conf.Set("app.debug", false)
//	conf.Set("app.verbose", true)
//	conf.Set("app.dump", 3)
//	conf.Set("app.logging.file", "/tmp/1.log")
//	conf.Set("app.server.start", 5)
//
//	ss := conf.WithPrefix("app.logging")
//	ss.Set("rotate", 6)
//	ss.Set("words", []any{"a", 1, false})
//	ss.Set("keys", map[any]any{"a": 3.13, 1.73: "zz", false: true})
//
//	conf.Set("app.bool", "[on,off,   true]")
//	conf.SetComment("app.bool", "a bool slice", "remarks here")
//	conf.SetTag("app.bool", []any{"on", "off", true})
//
//	states.Env().SetNoColorMode(true) // to disable ansi escape sequences in dump output
//	fmt.Println(conf.Dump())
//
//	data, found := conf.Get("app.logging.rotate")
//	println(data, found)
//	data := conf.MustGet("app.logging.rotate")
//	println(data)
//
// The `store` provides advanced APIs to extract typed data from node.
//
//	iData := conf.MustInt("app.logging.rotate")
//	debugMode := conf.MustBool("app.debug")
//	...
//
// The searching tool is also used to locate whether a key exists or not:
//
//	found := conf.Has("app.logging.rotate")
//	node, isBranch, isPartialMatched, found := conf.Locate("app.logging.rotate")
//	t.Logf("%v | %s | %v |     | %v, %v, found: %v", node.Data(), node.Comment(), node.Tag(), isBranch, isPartialMatched, found)
//
// The `store` provides many providers and codecs.
// A provider represents an external data source, such as file, environment, consul, etc.
// And a codec represents the data format, just like yaml, json, toml, etc.
//
// So an app can [Store.Load] the external yaml files like the following way:
//
//	func TestStoreS_Load(t *testing.T) {
//	    conf := newBasicStore(WithWatchEnable(true))
//	    defer conf.Close()
//	    ctx := context.Background()
//
//	    parser := yaml.New()
//	    _, err := conf.Load(ctx,
//	        store.WithStorePrefix("app.yaml"),
//	        store.WithCodec(parser),
//	        store.WithProvider(file.New("../../../testdata/2.yaml")),
//
//	        store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
//	    )
//
//	    assert.Equal(t, `-s`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.0"))
//	    assert.Equal(t, `-w`, s.MustGet("app.yaml.app.bgo.build.projects.000-default-group.items.001-bgo.ldflags.1"))
//
//	    m := map[string]any{
//	        "m1.s1": "cool",
//	        "m1.s2": 9,
//	        "key2": map[any]any{
//	            9: 1,
//	            8: false,
//	        },
//	        "slice": []map[any]any{
//	            {7.981: true, "cool": "maps"},
//	            {"hello": "world"},
//	        },
//	    }
//	    _, err := conf.Load(ctx,
//	        WithProvider(maps.New(m, ".")),
//	        WithStoreFlattenSlice(true),
//	        WithStorePrefix("app.maps"),
//	        WithPosition(""),
//	    )
//	    if ErrorIsNotFound(err) {
//	        t.Fail()
//	    }
//	    if err != nil {
//	        t.Fatalf("err: %v", err)
//	    }
//
//	    t.Logf("\nPath of 'conf' (delimeter=%v, prefix=%v)\n%v\n",
//	        conf.Delimiter(),
//	        conf.Prefix(),
//	        conf.Dump())
//
//	    assertEqual(t, false, conf.MustBool("app.maps.key2.8"))
//	    assertEqual(t, 1, conf.MustInt("app.maps.key2.9", -1))
//	    assertEqual(t, "cool", conf.MustString("app.maps.m1.s1"))
//	    assertEqual(t, 9, conf.MustInt("app.maps.m1.s2", -1))
//	}
//
// For more information, browse these public sites:
//
// - https://pkg.go.dev/github.com/hedzr/store
//
// - https://github.com/hedzr/store
package store

const Version = "v1.2.3" // Version of libs.store
