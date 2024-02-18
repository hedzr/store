package store

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/hedzr/logg/slog"
)

func TestNewLogEntry(t *testing.T) {
	le := NewLogEntry(time.Now(), slog.DebugLevel, "", getpc(0, 0, 0))
	le2 := le.Clone()

	le.AddAttrs(
		String("string", "hello"),
		Int64("int64", -112233),
		Int("int", -112233),
		Uint64("uint64", 112233),
		Uint("uint", 112233),
		Float64("float", 2.1),
		Bool("bool", true),
		Time("time", time.Now()),
		Duration("duration", time.Second),
		Group("group",
			Uint("uint", 112233),
			Float64("float", 2.1),
			Bool("bool", true),
			Time("time", time.Now()),
			Any("any", "x"),
		),
		Any("any", "y"),
	)
	t.Logf("%v / %v / %v", le, le2, le.NumAttrs())

	a1 := Any("any1", "y")
	a2 := Any("string1", "y")
	if a1.Equal(a2) {
		t.Fail()
	}
	t.Logf("a1: %v, empty: %v, kind: %v", a1.String(), a1.isEmpty(), a1.Kind())

	t.Logf("countAttrs: %v", countAttrs([]any{
		String("string", "hello"),
		Int64("int64", -112233),
		"ss",
		Group("empty-group"),
	}))

	le.Add(a1, a2)
	le.Add("z1", 1, "z2", 2)
	le.Add("z1")

	le.Attrs(func(item Item) bool {
		t.Log(item)
		return true
	})

	t.Logf("source: %v", le.source())
	t.Logf("source: %v", le2.source())
	t.Logf("group: %v", le.source().group())
}

func getpc(skip int, extra ...int) (pc uintptr) { //nolint:unused
	var pcs [1]uintptr
	for _, ee := range extra {
		if ee > 0 {
			skip += ee //nolint:revive
		}
	}
	runtime.Callers(skip+1, pcs[:])
	pc = pcs[0]
	return
}

func TestValue_Any(t *testing.T) {
	data := []any{
		int64(1),
		int32(1),
		int16(1),
		int8(1),
		int(1),
		uint64(1),
		uint32(1),
		uint16(1),
		uint8(1),
		uint(1),
		float32(1),
		float64(1),
		uintptr(1),
		true,
		time.Second,
		time.Now(),
		StringValue(""),
		[]Item{},
		KindAny,
		struct{}{},
	}

	for _, vv := range data {
		v := AnyValue(vv)
		t.Log(v.Any())
	}

	rv := AnyValue(fmt.Errorf("StoreValue panicked\n%s", stack(1, 2)))
	t.Log(rv)

	v := IntValue(1)
	t.Log(v.Int64(), rv.Equal(v), v.Equal(v))
	v = Int64Value(1)
	t.Log(v.Int64(), rv.Equal(v), v.Equal(v))
	v = Uint64Value(1)
	t.Log(v.Uint64(), rv.Equal(v), v.Equal(v))
	v = BoolValue(true)
	t.Log(v.Bool(), rv.Equal(v), v.Equal(v))
	v = DurationValue(time.Second)
	t.Log(v.Duration(), rv.Equal(v), v.Equal(v))
	v = TimeValue(time.Now())
	t.Log(v.Time(), rv.Equal(v), v.Equal(v))
	v = TimeValue(time.Time{})
	t.Log(v.Time(), rv.Equal(v), v.Equal(v))
	v = StringValue("true")
	t.Log(v.String(), rv.Equal(v), v.Equal(v))
	v = Float64Value(1.1)
	t.Log(v.Float64(), rv.Equal(v), v.Equal(v))
	v = AnyValue(struct{}{})
	t.Log(v.Any(), rv.Equal(v), v.Equal(v))
	v = GroupValue()
	v.Resolve()
	t.Log(v.Group(), rv.Equal(v), v.Equal(v))
	v = GroupValue(
		String("string", "hello"),
		Int64("int64", -112233),
		Int("int", -112233),
		Uint64("uint64", 112233),
		Uint("uint", 112233),
		Float64("float", 2.1),
		Bool("bool", true),
		Time("time", time.Now()),
		Duration("duration", time.Second),
	)
	v.Resolve()
	t.Log(v.Group(), rv.Equal(v), v.Equal(v))

	it := Any("any", AnyValue(struct{}{}))
	t.Log(it.String(), it.Kind())

	it = Any("any", AnyValue(nil))
	t.Log(it.String(), it.Kind())

	val := Value{}
	val.num, val.any = 0, nil
	it = Item{
		"",
		val,
	}
	t.Log(it.String(), it.Kind())
}
