package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"

	"github.com/hedzr/evendeep"
	"github.com/hedzr/store"
	"github.com/hedzr/store/codecs/hcl"
	"github.com/hedzr/store/codecs/toml"
	"github.com/hedzr/store/internal/times"
	"github.com/hedzr/store/providers/file"
	"github.com/hedzr/store/radix"
)

func ExampleNew() {
	conf := store.New()
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []any{"a", 1, false})
	ss.Set("keys", map[any]any{"a": 3.13, 1.73: "zz", false: true})

	radix.StatesEnvSetColorMode(true) // to disable ansi escape sequences in dump output
	_, _ = fmt.Println(conf.Dump())

	// Output:
	//   app.                          <B>
	//     d                           <B>
	//       ebug                      <L> app.debug => false
	//       ump                       <L> app.dump => 3
	//     verbose                     <L> app.verbose => true
	//     logging.                    <B>
	//       file                      <L> app.logging.file => /tmp/1.log
	//       rotate                    <L> app.logging.rotate => 6
	//       words                     <L> app.logging.words => [a 1 false]
	//       keys                      <L> app.logging.keys => map[a:3.13 1.73:zz false:true]
	//     server.start                <L> app.server.start => 5
}

func ExamplestoreS_Dump() {
	conf := store.New()
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []any{"a", 1, false})
	ss.Set("keys", map[any]any{"a": 3.13, 1.73: "zz", false: true})

	conf.Set("app.bool", "[on,off,   true]")
	conf.SetComment("app.bool", "a bool slice", "remarks here")
	conf.SetTag("app.bool", []any{"on", "off", true})

	radix.StatesEnvSetColorMode(true) // to disable ansi escape sequences in dump output
	_, _ = fmt.Println(conf.Dump())

	// Output:
	//   app.                          <B>
	//     d                           <B>
	//       ebug                      <L> app.debug => false
	//       ump                       <L> app.dump => 3
	//     verbose                     <L> app.verbose => true
	//     logging.                    <B>
	//       file                      <L> app.logging.file => /tmp/1.log
	//       rotate                    <L> app.logging.rotate => 6
	//       words                     <L> app.logging.words => [a 1 false]
	//       keys                      <L> app.logging.keys => map[a:3.13 1.73:zz false:true]
	//     server.start                <L> app.server.start => 5
	//     bool                        <L> app.bool => [on,off,   true] // remarks here | tag = [on off true] ~ a bool slice
}

func TestStoreS_Dump2(t *testing.T) {
	conf := store.New()
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []any{"a", 1, false})
	ss.Set("keys", map[any]any{"a": 3.13, 1.73: "zz", false: true})

	conf.Set("app.bool", "[on,off,   true]")
	conf.SetComment("app.bool", "a bool slice", "remarks here")
	conf.SetTag("app.bool", []any{"on", "off", true})

	// radix.StatesEnvSetColorMode(true) // to disable ansi escape sequences in dump output
	t.Log("\n", conf.Dump())
}

func TestStoreS_Get(t *testing.T) {
	conf := store.New()
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []any{"a", 1, false})
	ss.Set("keys", map[any]any{"a": 3.13, 1.73: "zz", false: true})

	conf.Set("app.bool", "[on,off,   true]")
	conf.SetComment("app.bool", "a bool slice", "remarks here")
	conf.SetTag("app.bool", []any{"on", "off", true})

	// data, found := conf.Get("app.logging.rotate")
	// println(data, found)
	// data = conf.MustGet("app.logging.rotate")
	// println(data)
	// fmt.Println(conf.MustInt("app.dump"))
	// fmt.Println(conf.MustString("app.dump"))
	// fmt.Println(conf.MustBool("app.dump")) // convert 3 to bool will get true since hedzr/evendeep v1.1.0

	data, found := conf.Get("app.logging.rotate")
	assert.Equal(t, 6, data)
	assert.Equal(t, true, found)
	data = conf.MustGet("app.logging.rotate")
	assert.Equal(t, 6, data)

	assert.Equal(t, 3, conf.MustInt("app.dump"))
	assert.Equal(t, "3", conf.MustString("app.dump"))
	assert.Equal(t, true, conf.MustBool("app.dump"))

	assert.Equal(t, []bool{false, true, false}, conf.MustBoolSlice("app.logging.words"))
	assert.Equal(t, []string{"a", "1", "false"}, conf.MustStringSlice("app.logging.words"))
	assert.Equal(t, []int{0, 1, 0}, conf.MustIntSlice("app.logging.words"))
	assert.Equal(t, []int32{0, 1, 0}, conf.MustInt32Slice("app.logging.words"))
	assert.Equal(t, []uint32{0, 1, 0}, conf.MustUint32Slice("app.logging.words"))
	assert.Equal(t, map[string]any{"words": []any{"a", 1, false}}, conf.MustM("app.logging.words"))

	assert.Equal(t, map[string]string{"a": "3.13", "1.73": "zz", "false": "true"}, conf.MustStringMap("app.logging.keys"))
	assert.Equal(t, map[any]any{"a": 3.13, 1.73: "zz", false: true}, conf.MustGet("app.logging.keys"))
	assert.Equal(t, map[string]any{"keys": map[any]any{"a": 3.13, 1.73: "zz", false: true}},
		conf.MustM("app.logging.keys"))
	assert.Equal(t, map[string]any{"app.logging.keys": map[any]any{"a": 3.13, 1.73: "zz", false: true}},
		conf.MustR("app.logging.keys"))

	t.Logf("\n%v", conf.Dump())
}

func TestStoreResultTypedGetters(t *testing.T) {
	s := store.New()
	parser := hcl.New() // hcl.WithFlattenSlices(true))
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.hcl"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.hcl")),

		store.WithStoreFlattenSlice(true),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, false, s.MustGet("app.hcl.debug").(bool))

	assert.Equal(t, `12345`, s.MustString("app.hcl.server.0.getters.0.mf"))
	assert.Equal(t, `-12.7831`, s.MustString("app.hcl.server.0.getters.0.ff"))
	assert.Equal(t, `true`, s.MustString("app.hcl.server.0.getters.0.bf"))
	assert.Equal(t, `1s52ms`, s.MustString("app.hcl.server.0.getters.0.tf"))
	assert.Equal(t, `2023-1-2`, s.MustGet("app.hcl.server.0.getters.0.time"))

	assert.Equal(t, int64(-12), s.MustInt64("app.hcl.server.0.getters.0.ff"))
}

func TestTypedGetters(t *testing.T) {
	s := store.New()
	parser := toml.New()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.toml"),
		store.WithCodec(parser),
		store.WithProvider(file.New("../testdata/getters.toml")),
	); err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	t.Logf("\n%-32sData\n%v\n", "Path", s.Dump())

	assert.Equal(t, int64(5000), s.MustGet("app.toml.database.connection_max"))

	assert.Equal(t, times.MustSmartParseTime("1979-05-27 07:32:00 -0800"), s.MustGet("app.toml.owner.dob"))
	assert.Equal(t, time.Minute+37*time.Second+512*time.Millisecond, times.MustParseDuration(s.MustString("app.toml.owner.duration")))
}

func TestDecompoundMap(t *testing.T) {
	conf := newBasicStore()

	conf.Set("app.map", false)
	err := conf.Merge("app.map", map[string]any{
		"k1": 1,
		"k2": false,
		"m3": map[string]any{
			"bobo": "joe",
		},
	})
	if err != nil {
		t.Fatalf("Merge failed: %v", err)
	}

	assert.Equal(t, int(1), conf.MustGet("app.map.k1"))
	assert.Equal(t, false, conf.MustGet("app.map.k2"))
	assert.Equal(t, "joe", conf.MustGet("app.map.m3.bobo"))
}

func TestTagAndComment(t *testing.T) {
	conf := newBasicStore()

	found := conf.Has("app.logging.rotate")
	println(found)
	node, isBranch, isPartialMatched, found := conf.Locate("app.logging.rotate")
	t.Logf("%v    | %v, %v, found: %v", node.Data(), isBranch, isPartialMatched, found)

	conf.Set("debug", false)
	conf.SetComment("debug", "a flag to identify app debug mode", "remarks here")
	conf.SetTag("debug", map[string]any{
		"handler": func() {},
	})

	node, _, _, found = conf.Locate("debug")
	if found {
		t.Log(node.Tag(), node.Description(), node.Comment())
	}

	println(conf.Dump())
}

//

//

//

func assertEqual(t testing.TB, expect, actual any, msg ...any) { //nolint:govet //it's a printf/println dual interface
	if evendeep.DeepEqual(expect, actual) {
		return
	}

	var mesg string
	if len(msg) > 0 {
		if format, ok := msg[0].(string); ok {
			mesg = fmt.Sprintf(format, msg[1:]...)
		} else {
			mesg = fmt.Sprint(msg...)
		}
	}

	t.Fatalf("assertEqual failed: %v\n    expect: %v\n    actual: %v\n", mesg, spew.Sdump(expect), spew.Sdump(actual))
}
