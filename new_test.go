package store

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/hedzr/store/internal/cvt"
)

func TestRandomStringPure(t *testing.T) {
	// t.Log(stringtool.RandomStringPure(8))
	t.Log("End.")
}

func TestWithinLoading(t *testing.T) {
	conf := newBasicStore()
	defer conf.Close()

	t.Logf("\nPath of 'conf' (delimeter=%v, prefix=%v)\n%v\n",
		conf.Delimiter(),
		conf.Prefix(),
		conf.Dump())

	conf.WithinLoading(func() {
		t.Log("")
	})

	assertEqual(t, 5, conf.MustInt("app.server.start"))
	assertEqual(t, 6, conf.MustInt("app.logging.rotate", -1))

	// load with the explicit empty provider and codec

	ctx := context.TODO()
	_, err := conf.Load(ctx,
		WithCodec(nil),
		WithProvider(nil),
		WithStoreFlattenSlice(false),
		WithStorePrefix(""),
		WithPosition(""),
	)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	conf.Remove("app.logging")

	assertEqual(t, 5, conf.MustInt("app.server.start"))
	assertEqual(t, -1, conf.MustInt("app.logging.rotate", -1))

	t.Logf("\nPath of 'conf' (delimeter=%v, prefix=%v)\n%v\n",
		conf.Delimiter(),
		conf.Prefix(),
		conf.Dump())
}

func TestStoreS_join(t *testing.T) {
	// join ...

	conf := newBasicStore()
	assertEqual(t, "", conf.join())
	assertEqual(t, "", conf.join(""))
	assertEqual(t, "x", conf.join("x", ""))
	assertEqual(t, "x.y", conf.join("x", "y"))
	assertEqual(t, "x.y", conf.join("x", "y", ""))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z"))
	assertEqual(t, "x.y.z", conf.join("x", "y", "z", ""))
	assertEqual(t, "x.y.z", conf.join("", "x", "y", "z", ""))
}

func TestOp_Marshal(t *testing.T) {
	var op Op
	for k, v := range opStrings {
		err := (&op).UnmarshalText([]byte(v))
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		if op != k {
			t.Fatalf("unmarshaltext failed, expecting '%v' but got '%v'", k, op)
		}

		b, err := op.MarshalText()
		if err != nil {
			t.Fatalf("failed: %v", err)
		}
		t.Logf("marshal text: %v", string(b))

		_ = op.Marshal()
	}

	err := (&op).UnmarshalText([]byte("???"))
	if err == nil {
		t.Fail()
	}

	op = Op(65535)
	v, _ := op.MarshalText()
	if !strings.HasPrefix(string(v), "Op(") {
		t.Fail()
	}
}

func TestStoreS_Load(t *testing.T) {
	conf := newBasicStore(WithWatchEnable(true))
	defer conf.Close()

	// load with a small copy of maps provider, it's inline here for testing only

	m := map[string]any{
		"m1.s1": "cool", // a key containing delimiter should be splitted as children automatically
		"m1.s2": 9,
		"key2": map[any]any{
			9: 1,
			8: false,
		},
		"slice": []map[any]any{
			{7.981: true, "st": true, "cool": "maps"},
			{"hello": "world"},
		},
	}
	ctx := context.TODO()
	_, err := conf.Load(ctx,
		WithProvider(newMaps(m, ".")),
		WithStorePrefix("app.maps"),
		WithPosition(""),
		WithStoreFlattenSlice(true), // this option allows 'slice' to split as child nodes
	)
	if ErrorIsNotFound(err) {
		t.Fail()
	}
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	t.Logf("\nPath of 'conf' (delimeter=%v, prefix=%v)\n%v\n",
		conf.Delimiter(),
		conf.Prefix(),
		conf.Dump())

	assertEqual(t, false, conf.MustBool("app.maps.key2.8"))
	assertEqual(t, 1, conf.MustInt("app.maps.key2.9", -1))
	assertEqual(t, "cool", conf.MustString("app.maps.m1.s1"))
	assertEqual(t, 9, conf.MustInt("app.maps.m1.s2", -1))
	// assertEqual(t, true, conf.MustBool("app.maps.slice.0.7.981", false))
	assertEqual(t, "true", conf.MustString("app.maps.slice.0.st"))
}

//

//

//

func newMaps(m map[string]any, delimiter string, opts ...mapsOpt) *pvdr { //nolint:revive
	s := &pvdr{delimiter: delimiter}
	for _, opt := range opts {
		opt(s)
	}

	cp := cvt.Normalize(cvt.Copy(m), nil)
	if s.delimiter != "" {
		cp = cvt.Deflate(cp, s.delimiter)
	}

	// s.m = cp

	s.m = make(map[string]ValPkg)
	for k, v := range cp {
		s.m[k] = ValPkg{
			Value:   v,
			Desc:    "",
			Comment: "",
			Tag:     nil,
		}
	}
	return s
}

type mapsOpt func(s *pvdr)

type pvdr struct {
	m         map[string]ValPkg
	delimiter string
	prefix    string
	codec     Codec
	watching  int32
}

func mapsWithCodec(codec Codec) mapsOpt { //nolint:unused
	return func(s *pvdr) {
		s.codec = codec
	}
}

func mapsWithPosition(prefix string) mapsOpt { //nolint:unused
	return func(s *pvdr) {
		s.prefix = prefix
	}
}

func mapsWithDelimiter(d string) mapsOpt { //nolint:unused
	return func(s *pvdr) {
		s.delimiter = d
	}
}

func (s *pvdr) GetCodec() (codec Codec)    { return s.codec }
func (s *pvdr) GetPosition() (pos string)  { return s.prefix }
func (s *pvdr) WithCodec(codec Codec)      { s.codec = codec }
func (s *pvdr) WithPosition(prefix string) { s.prefix = prefix }

func (s *pvdr) Count() int {
	return 0
}

func (s *pvdr) Has(key string) bool { //nolint:revive
	return false
}

func (s *pvdr) Next() (key string, eol bool) {
	eol = true
	return
}

func (s *pvdr) Keys() (keys []string, err error) {
	err = ErrNotImplemented
	return
}

func (s *pvdr) Value(key string) (value any, ok bool) { //nolint:revive
	ok = false
	return
}

func (s *pvdr) MustValue(key string) (value any) { //nolint:revive
	return
}

func (s *pvdr) Reader() (r Reader, err error) { //nolint:revive
	err = ErrNotImplemented
	return
}

// Read returns the loaded map[string]interface{}.
func (s *pvdr) Read() (data map[string]ValPkg, err error) {
	return s.m, nil
}

// ReadBytes is not supported by the confmap provider.
func (s *pvdr) ReadBytes() (data []byte, err error) {
	err = ErrNotImplemented
	return
}

func (s *pvdr) Write(data []byte) (err error) { //nolint:revive
	err = ErrNotImplemented
	return
}

func (s *pvdr) Close() {
	atomic.CompareAndSwapInt32(&s.watching, 1, 0)
}

func (s *pvdr) Watch(ctx context.Context, cb func(event any, err error)) (err error) {
	if !atomic.CompareAndSwapInt32(&s.watching, 0, 1) {
		return
	}

	// do some stuff here
	_, _ = ctx, cb
	return
}
