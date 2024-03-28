package radix

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hedzr/store/internal/times"
)

func TestTrieS_GetR(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	m, err := trie.GetR("app")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('app') returns:\n\n%v\n\n", m)
	// t.Logf("GetR('app') returns:")
	// spew.Default.Println(m)

	var m1 map[string]any
	m1, err = trie.GetR("")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('') returns:\n\n%v\n\n", m1)
	// t.Logf("GetR('') returns:")
	// spew.Default.Println(m)

	m, err = trie.GetR("app")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('app') returns:\n\n%v\n\n", m)

	if !reflect.DeepEqual(m, m1) {
		t.Fatalf("expecting m == m1 but failed.\n  m : %v\n  m1: %v", m, m1)
	}

	m, err = trie.GetR("app.logging")
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR('app.logging') returns:\n\n%v\n\n", m)

	m, err = trie.GetR("unexist", nil)
	if err != nil {
		t.Fatalf("GetR failed: %v", err)
	}
	t.Logf("GetR(unexist'') returns:\n\n%v\n\n", m)

	m = trie.MustR("unexist", nil)
	t.Logf("MustR(unexist'') returns:\n\n%v\n\n", m)
}

func TestTrieS_GetM(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	m2, err := trie.GetM("")
	t.Logf("map: %v | err: %v", m2, err)

	m, err := trie.GetM("app.logging")
	if err != nil {
		t.Fatalf("GetM failed: %v", err)
	}
	t.Logf("GetM('app.logging') returns:\n\n%v\n\n", m) // no prefix

	m, err = trie.GetM("", WithFilter[any](func(node Node[any]) bool {
		return strings.HasPrefix(node.Key(), "app.")
	}))
	if err != nil {
		t.Fatalf("GetM failed: %v", err)
	}
	t.Logf("GetM('app') returns:\n\n%v\n\n", m)

	if !reflect.DeepEqual(m, m2) {
		t.Fatalf("expecting m == m2, but:\n  m : %v\n  m2: %v\n", m, m2)
	}

	m, err = trie.GetM("", WithFilter[any](func(node Node[any]) bool {
		return strings.HasPrefix(node.Key(), "app.logging.")
	}), WithoutFlattenKeys[any](true))
	m1 := trie.MustM("app.logging", WithKeepPrefix[any](false))
	m2 = trie.MustM("app.logging", WithKeepPrefix[any](true))
	if !reflect.DeepEqual(m, m2) {
		t.Fatalf("expecting m == m2, but:\n  m : %v\n  m1: %v\n  m2: %v\n", m, m1, m2)
	}

	m = trie.MustM("app.logging")
	if err != nil {
		t.Fatalf("MustM failed: %v", err)
	}
	t.Logf("MustM('app.logging') returns:\n\n%v\n\n", m)

	// t.Logf("GetR('app') returns:")
	// spew.Default.Println(m)
}

func TestTrieS_GetSectionFrom(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(true)
	t.Logf("\nPath\n%v\n", ret)
	// assertTrue(t, trie.Search("apple"), `expecting trie.Search("apple") return true`)     // 返回 True
	// assertFalse(t, trie.Search("app"), `expecting trie.Search("app") return false`)       // 返回 False
	// assertTrue(t, trie.StartsWith("app"), `expecting trie.StartsWith("app") return true`) // 返回 True
	// trie.Insert("app")
	// assertTrue(t, trie.Search("app"), `expecting trie.Search("app") return true`) // 返回 True

	type loggingS struct {
		File   uint
		Rotate uint64
		Words  []any
	}

	type serverS struct {
		Start int
		Sites int
	}

	type appS struct {
		Debug   int
		Dump    int
		Verbose int64
		Logging loggingS
		Server  serverS
	}

	type cfgS struct {
		App appS
	}

	// m := trie.MustM("")
	// m1 := trie.splitCompactKeys(m)
	// t.Logf("MustM('') returns m1: \n%v", spew.Sdump(m1)) // github.com/davecgh/go-spew/spew

	// app.                          <B>
	//   d                           <B>
	//     ebug                      <L> app.debug => 1
	//     ump                       <L> app.dump => 3
	//   verbose                     <L> app.verbose => 2
	//   logging.                    <B>
	//     file                      <L> app.logging.file => 4
	//     rotate                    <L> app.logging.rotate => 6
	//     words                     <L> app.logging.words => [a 1 false]
	//   server.s                    <B>
	//     tart                      <L> app.server.start => 5
	//     ites                      <L> app.server.sites => 1

	var ss cfgS
	err := trie.GetSectionFrom("", &ss)
	t.Logf("cfgS: %v | err: %v", ss, err)

	if !reflect.DeepEqual(ss.App.Logging.Words, []any{"a", 1, false}) {
		t.Fail()
	}
}

func TestTrieS_GetString(t *testing.T) {
	trie := newTrieTree()
	ss, _ := trie.GetString("app.logging.words")
	t.Logf("app.logging.words: %v", ss)
	if ss != "[a 1 false]" {
		t.Fail()
	}

	ss = trie.MustString("app.logging.words")
	t.Logf("app.logging.words: %v", ss)
	ss = trie.MustString("app.logging.words111", "")
	t.Logf("app.logging.words: %v", ss)

	a := trie.MustStringSlice("app.logging.words")
	t.Logf("app.logging.words: %v", a)
	a = trie.MustStringSlice("app.logging.words111", "")
	t.Logf("app.logging.words: %v", a)

	b, e := trie.GetStringSlice("app.logging.words111", "")
	if e == nil {
		t.Fail()
	}
	t.Logf("app.logging.words: %v", b)
}

func TestTrieS_GetStringMap(t *testing.T) {
	trie := newTrieTree()
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)
	ss, _ := trie.GetStringMap("app")
	t.Logf("app.logging.words: %v", ss)
	// if ss != "[a 1 false]" {
	// 	t.Fail()
	// }

	trie.Set("app.map", map[string]string{"hello": "world"})
	ss = trie.MustStringMap("app.map")
	t.Logf("app.map: %v", ss)

	trie.Set("app.map", map[string]any{"hello": "world"})
	ss = trie.MustStringMap("app.map")
	t.Logf("app.map: %v", ss)

	ss = trie.MustStringMap("app.map.absent", map[string]string{"ok": "bye"})
	t.Logf("app.map.absent: %v", ss)
}

func TestTrieS_GetInt(t *testing.T) { //nolint:revive
	trie := newTrieTree()
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	trie.Set("app.int", int64(-123))

	ss64, err := trie.GetInt64("app.int")
	t.Logf("app.int64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustInt64("app.int")
	t.Logf("app.int64: %v", ss64)

	ss64, err = trie.GetInt64("app.int-absent", 9)
	t.Logf("app.int64: %v, err: %v", ss64, err)
	ss64 = trie.MustInt64("app.int-absent", 9)
	t.Logf("app.int64-absent: %v", ss64)

	ss32, err := trie.GetInt32("app.int")
	t.Logf("app.int32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustInt32("app.int")
	t.Logf("app.int32: %v", ss32)

	ss32, err = trie.GetInt32("app.int-absent", 9)
	t.Logf("app.int32: %v, err: %v", ss32, err)
	ss32 = trie.MustInt32("app.int-absent", 9)
	t.Logf("app.int32-absent: %v", ss32)

	ss16, err := trie.GetInt16("app.int")
	t.Logf("app.int16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustInt16("app.int")
	t.Logf("app.int16: %v", ss16)

	ss16, err = trie.GetInt16("app.int-absent", 9)
	t.Logf("app.int16: %v, err: %v", ss16, err)
	ss16 = trie.MustInt16("app.int-absent", 9)
	t.Logf("app.int16-absent: %v", ss16)

	ss8, err := trie.GetInt8("app.int")
	t.Logf("app.int8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustInt8("app.int")
	t.Logf("app.int8: %v", ss8)

	ss8, err = trie.GetInt8("app.int-absent", 9)
	t.Logf("app.int8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustInt8("app.int-absent", 9)
	t.Logf("app.int8-absent: %v", ss8)

	ss, err := trie.GetInt("app.int")
	t.Logf("app.int: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustInt("app.int")
	t.Logf("app.int: %v", ss)

	ss, err = trie.GetInt("app.int-absent", 9)
	t.Logf("app.int-absent: %v, err: %v", ss, err)
	ss = trie.MustInt("app.int-absent", 9)
	t.Logf("app.int-absent: %v", ss)
}

func TestTrieS_GetUint(t *testing.T) { //nolint:revive
	trie := newTrieTree()
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	trie.Set("app.uint", uint64(123))

	ss64, err := trie.GetUint64("app.uint")
	t.Logf("app.uint64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustUint64("app.uint")
	t.Logf("app.uint64: %v", ss64)

	ss64, err = trie.GetUint64("app.uint-absent", 9)
	t.Logf("app.uint64-absent: %v, err: %v", ss64, err)
	ss64 = trie.MustUint64("app.uint-absent", 9)
	t.Logf("app.uint64-absent: %v", ss64)

	ss32, err := trie.GetUint32("app.uint")
	t.Logf("app.uint32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustUint32("app.uint")
	t.Logf("app.uint32: %v", ss32)

	ss32, err = trie.GetUint32("app.uint-absent", 9)
	t.Logf("app.uint32-absent: %v, err: %v", ss32, err)
	ss32 = trie.MustUint32("app.uint-absent", 9)
	t.Logf("app.uint32-absent: %v", ss32)

	ss16, err := trie.GetUint16("app.uint")
	t.Logf("app.uint16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustUint16("app.uint")
	t.Logf("app.uint16: %v", ss16)

	ss16, err = trie.GetUint16("app.uint-absent", 9)
	t.Logf("app.uint16-absent: %v, err: %v", ss16, err)
	ss16 = trie.MustUint16("app.uint-absent", 9)
	t.Logf("app.uint16-absent: %v", ss16)

	ss8, err := trie.GetUint8("app.uint")
	t.Logf("app.uint8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustUint8("app.uint")
	t.Logf("app.uint8: %v", ss8)

	ss8, err = trie.GetUint8("app.uint-absent", 9)
	t.Logf("app.uint8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustUint8("app.uint-absent", 9)
	t.Logf("app.uint8-absent: %v", ss8)

	ss, err := trie.GetUint("app.uint")
	t.Logf("app.uint: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustUint("app.uint")
	t.Logf("app.uint: %v", ss)

	ss, err = trie.GetUint("app.uint-absent", 9)
	t.Logf("app.uint-absent: %v, err: %v", ss, err)
	ss = trie.MustUint("app.uint-absent", 9)
	t.Logf("app.uint-absent: %v", ss)
}

func TestTrieS_GetIntSlice(t *testing.T) { //nolint:revive
	trie := newTrieTree()
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	trie.Set("app.int", []int64{-123, 73})

	ss64, err := trie.GetInt64Slice("app.int")
	t.Logf("app.int64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustInt64Slice("app.int")
	t.Logf("app.int64: %v", ss64)

	ss64, err = trie.GetInt64Slice("app.int-absent", 9)
	t.Logf("app.int64: %v, err: %v", ss64, err)
	ss64 = trie.MustInt64Slice("app.int-absent", 9)
	t.Logf("app.int64-absent: %v", ss64)

	ss32, err := trie.GetInt32Slice("app.int")
	t.Logf("app.int32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustInt32Slice("app.int")
	t.Logf("app.int32: %v", ss32)

	ss32, err = trie.GetInt32Slice("app.int-absent", 9)
	t.Logf("app.int32: %v, err: %v", ss32, err)
	ss32 = trie.MustInt32Slice("app.int-absent", 9)
	t.Logf("app.int32-absent: %v", ss32)

	ss16, err := trie.GetInt16Slice("app.int")
	t.Logf("app.int16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustInt16Slice("app.int")
	t.Logf("app.int16: %v", ss16)

	ss16, err = trie.GetInt16Slice("app.int-absent", 9)
	t.Logf("app.int16: %v, err: %v", ss16, err)
	ss16 = trie.MustInt16Slice("app.int-absent", 9)
	t.Logf("app.int16-absent: %v", ss16)

	ss8, err := trie.GetInt8Slice("app.int")
	t.Logf("app.int8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustInt8Slice("app.int")
	t.Logf("app.int8: %v", ss8)

	ss8, err = trie.GetInt8Slice("app.int-absent", 9)
	t.Logf("app.int8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustInt8Slice("app.int-absent", 9)
	t.Logf("app.int8-absent: %v", ss8)

	ss, err := trie.GetIntSlice("app.int")
	t.Logf("app.int: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustIntSlice("app.int")
	t.Logf("app.int: %v", ss)

	ss, err = trie.GetIntSlice("app.int-absent", 9)
	t.Logf("app.int-absent: %v, err: %v", ss, err)
	ss = trie.MustIntSlice("app.int-absent", 9)
	t.Logf("app.int-absent: %v", ss)
}

func TestTrieS_GetUintSlice(t *testing.T) { //nolint:revive
	trie := newTrieTree()
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	trie.Set("app.uint", []uint64{123})

	ss64, err := trie.GetUint64Slice("app.uint")
	t.Logf("app.uint64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustUint64Slice("app.uint")
	t.Logf("app.uint64: %v", ss64)

	ss64, err = trie.GetUint64Slice("app.uint-absent", 9)
	t.Logf("app.uint64-absent: %v, err: %v", ss64, err)
	ss64 = trie.MustUint64Slice("app.uint-absent", 9)
	t.Logf("app.uint64-absent: %v", ss64)

	ss32, err := trie.GetUint32Slice("app.uint")
	t.Logf("app.uint32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustUint32Slice("app.uint")
	t.Logf("app.uint32: %v", ss32)

	ss32, err = trie.GetUint32Slice("app.uint-absent", 9)
	t.Logf("app.uint32-absent: %v, err: %v", ss32, err)
	ss32 = trie.MustUint32Slice("app.uint-absent", 9)
	t.Logf("app.uint32-absent: %v", ss32)

	ss16, err := trie.GetUint16Slice("app.uint")
	t.Logf("app.uint16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustUint16Slice("app.uint")
	t.Logf("app.uint16: %v", ss16)

	ss16, err = trie.GetUint16Slice("app.uint-absent", 9)
	t.Logf("app.uint16-absent: %v, err: %v", ss16, err)
	ss16 = trie.MustUint16Slice("app.uint-absent", 9)
	t.Logf("app.uint16-absent: %v", ss16)

	ss8, err := trie.GetUint8Slice("app.uint")
	t.Logf("app.uint8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustUint8Slice("app.uint")
	t.Logf("app.uint8: %v", ss8)

	ss8, err = trie.GetUint8Slice("app.uint-absent", 9)
	t.Logf("app.uint8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustUint8Slice("app.uint-absent", 9)
	t.Logf("app.uint8-absent: %v", ss8)

	ss, err := trie.GetUintSlice("app.uint")
	t.Logf("app.uint: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustUintSlice("app.uint")
	t.Logf("app.uint: %v", ss)

	ss, err = trie.GetUintSlice("app.uint-absent", 9)
	t.Logf("app.uint-absent: %v, err: %v", ss, err)
	ss = trie.MustUintSlice("app.uint-absent", 9)
	t.Logf("app.uint-absent: %v", ss)
}

func TestTrieS_GetIntMap(t *testing.T) { //nolint:revive
	trie := newTrieTree()

	trie.Set("app.int", map[string]int64{"neg": -123, "pos": 73})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetInt64Map("app.int")
	t.Logf("app.int64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustInt64Map("app.int")
	t.Logf("app.int64: %v", ss64)

	ss64, err = trie.GetInt64Map("app.int-absent", map[string]int64{"t": 9})
	t.Logf("app.int64: %v, err: %v", ss64, err)
	ss64 = trie.MustInt64Map("app.int-absent", map[string]int64{"t": 9})
	t.Logf("app.int64-absent: %v", ss64)

	ss32, err := trie.GetInt32Map("app.int")
	t.Logf("app.int32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustInt32Map("app.int")
	t.Logf("app.int32: %v", ss32)

	ss32, err = trie.GetInt32Map("app.int-absent", map[string]int32{"t": 9})
	t.Logf("app.int32: %v, err: %v", ss32, err)
	ss32 = trie.MustInt32Map("app.int-absent", map[string]int32{"t": 9})
	t.Logf("app.int32-absent: %v", ss32)

	ss16, err := trie.GetInt16Map("app.int")
	t.Logf("app.int16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustInt16Map("app.int")
	t.Logf("app.int16: %v", ss16)

	ss16, err = trie.GetInt16Map("app.int-absent", map[string]int16{"t": 9})
	t.Logf("app.int16: %v, err: %v", ss16, err)
	ss16 = trie.MustInt16Map("app.int-absent", map[string]int16{"t": 9})
	t.Logf("app.int16-absent: %v", ss16)

	ss8, err := trie.GetInt8Map("app.int")
	t.Logf("app.int8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustInt8Map("app.int")
	t.Logf("app.int8: %v", ss8)

	ss8, err = trie.GetInt8Map("app.int-absent", map[string]int8{"t": 9})
	t.Logf("app.int8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustInt8Map("app.int-absent", map[string]int8{"t": 9})
	t.Logf("app.int8-absent: %v", ss8)

	ss, err := trie.GetIntMap("app.int")
	t.Logf("app.int: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustIntMap("app.int")
	t.Logf("app.int: %v", ss)

	ss, err = trie.GetIntMap("app.int-absent", map[string]int{"t": 9})
	t.Logf("app.int-absent: %v, err: %v", ss, err)
	ss = trie.MustIntMap("app.int-absent", map[string]int{"t": 9})
	t.Logf("app.int-absent: %v", ss)
}

func TestTrieS_GetUintMap(t *testing.T) { //nolint:revive
	trie := newTrieTree()

	trie.Set("app.uint", map[string]uint64{"non-neg": 2, "pos": 73})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetUint64Map("app.uint")
	t.Logf("app.uint64: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustUint64Map("app.uint")
	t.Logf("app.uint64: %v", ss64)

	ss64, err = trie.GetUint64Map("app.uint-absent", map[string]uint64{"t": 9})
	t.Logf("app.uint64-absent: %v, err: %v", ss64, err)
	ss64 = trie.MustUint64Map("app.uint-absent", map[string]uint64{"t": 9})
	t.Logf("app.uint64-absent: %v", ss64)

	ss32, err := trie.GetUint32Map("app.uint")
	t.Logf("app.uint32: %v", ss32)
	if err != nil {
		t.Fail()
	}
	ss32 = trie.MustUint32Map("app.uint")
	t.Logf("app.uint32: %v", ss32)

	ss32, err = trie.GetUint32Map("app.uint-absent", map[string]uint32{"t": 9})
	t.Logf("app.uint32-absent: %v, err: %v", ss32, err)
	ss32 = trie.MustUint32Map("app.uint-absent", map[string]uint32{"t": 9})
	t.Logf("app.uint32-absent: %v", ss32)

	ss16, err := trie.GetUint16Map("app.uint")
	t.Logf("app.uint16: %v", ss16)
	if err != nil {
		t.Fail()
	}
	ss16 = trie.MustUint16Map("app.uint")
	t.Logf("app.uint16: %v", ss16)

	ss16, err = trie.GetUint16Map("app.uint-absent", map[string]uint16{"t": 9})
	t.Logf("app.uint16-absent: %v, err: %v", ss16, err)
	ss16 = trie.MustUint16Map("app.uint-absent", map[string]uint16{"t": 9})
	t.Logf("app.uint16-absent: %v", ss16)

	ss8, err := trie.GetUint8Map("app.uint")
	t.Logf("app.uint8: %v", ss8)
	if err != nil {
		t.Fail()
	}
	ss8 = trie.MustUint8Map("app.uint")
	t.Logf("app.uint8: %v", ss8)

	ss8, err = trie.GetUint8Map("app.uint-absent", map[string]uint8{"t": 9})
	t.Logf("app.uint8-absent: %v, err: %v", ss8, err)
	ss8 = trie.MustUint8Map("app.uint-absent", map[string]uint8{"t": 9})
	t.Logf("app.uint8-absent: %v", ss8)

	ss, err := trie.GetUintMap("app.uint")
	t.Logf("app.uint: %v", ss)
	if err != nil {
		t.Fail()
	}
	ss = trie.MustUintMap("app.uint")
	t.Logf("app.uint: %v", ss)

	ss, err = trie.GetUintMap("app.uint-absent", map[string]uint{"t": 9})
	t.Logf("app.uint-absent: %v, err: %v", ss, err)
	ss = trie.MustUintMap("app.uint-absent", map[string]uint{"t": 9})
	t.Logf("app.uint-absent: %v", ss)
}

func TestTrieS_GetKibiBytes(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.kibi", "3gb") // 3G = 3*1024*1024*1024 = 3,221,225,472
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetKibiBytes("app.kibi")
	t.Logf("app.kibi: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustKibiBytes("app.kibi")
	t.Logf("app.kibi: %v", ss64)
	ss := trie.MustKibiBytes("app.kibi-absent", 102)
	t.Logf("app.kibi-absent: %v", ss)
	ss, err = trie.GetKibiBytes("app.kibi-absent", 102)
	t.Logf("app.kibi-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetKiloBytes(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.kilo", "3gb") // 3G = 3*1000*1000*1000 = 3,000,000,000
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetKiloBytes("app.kilo")
	t.Logf("app.kilo: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustKiloBytes("app.kilo")
	t.Logf("app.kilo: %v", ss64)
	ss := trie.MustKiloBytes("app.kilo-absent", 102)
	t.Logf("app.kilo-absent: %v", ss)
	ss, err = trie.GetKiloBytes("app.kilo-absent", 102)
	t.Logf("app.kilo-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat64(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", "3.14159")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat64("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat64("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat64("app.float-absent", 102)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat64("app.float-absent", 102)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat32(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", "3.14159")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat32("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat32("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat32("app.float-absent", 102)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat32("app.float-absent", 102)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat64Slice(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", []string{"3.14159"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat64Slice("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat64Slice("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat64Slice("app.float-absent", 102)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat64Slice("app.float-absent", 102)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat32Slice(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", []string{"3.14159"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat32Slice("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat32Slice("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat32Slice("app.float-absent", 102)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat32Slice("app.float-absent", 102)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat64Map(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", map[string]string{"e": "3.14159"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat64Map("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat64Map("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat64Map("app.float-absent", map[string]float64{"a": 102.1})
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat64Map("app.float-absent", map[string]float64{"a": 102.1})
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetFloat32Map(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", map[string]string{"e": "3.14159"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetFloat32Map("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustFloat32Map("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustFloat32Map("app.float-absent", map[string]float32{"a": 102.1})
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetFloat32Map("app.float-absent", map[string]float32{"a": 102.1})
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex128(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", "3.14159+1.3579i")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex128("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex128("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex128("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex128("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex64(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", "3.14159+1.3579i")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex64("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex64("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex64("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex64("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex128Slice(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", []string{"3.14159+1.3579i"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex128Slice("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex128Slice("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex128Slice("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex128Slice("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex64Slice(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", []string{"3.14159+1.3579i"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex64Slice("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex64Slice("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex64Slice("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex64Slice("app.float-absent", 102.1+56.35i)
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex128Map(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", map[string]string{"e": "3.14159+1.3579i"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex128Map("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex128Map("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex128Map("app.float-absent", map[string]complex128{"a": 102.1 + 56.35i})
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex128Map("app.float-absent", map[string]complex128{"a": 102.1 + 56.35i})
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetComplex64Map(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.float", map[string]string{"e": "3.14159+1.3579i"})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetComplex64Map("app.float")
	t.Logf("app.float: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustComplex64Map("app.float")
	t.Logf("app.float: %v", ss64)
	ss := trie.MustComplex64Map("app.float-absent", map[string]complex64{"a": 102.1 + 56.35i})
	t.Logf("app.float-absent: %v", ss)
	ss, err = trie.GetComplex64Map("app.float-absent", map[string]complex64{"a": 102.1 + 56.35i})
	t.Logf("app.float-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetBool(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.bool", "on")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetBool("app.bool")
	t.Logf("app.bool: %v", ss64)
	if err != nil {
		t.Fail()
	}
	ss64 = trie.MustBool("app.bool")
	t.Logf("app.bool: %v", ss64)
	ss := trie.MustBool("app.bool-absent", true)
	t.Logf("app.bool-absent: %v", ss)
	ss, err = trie.GetBool("app.bool-absent", true)
	t.Logf("app.bool-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetBoolSlice(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.bool", "[on,off,   true]")
	trie.SetComment("app.bool", "a bool slice", "remarks here")
	trie.SetTag("app.bool", []any{"on", "off", true})
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetBoolSlice("app.bool")
	t.Logf("app.bool: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, []bool{true, false, true}) {
		t.Fail()
	}
	ss64 = trie.MustBoolSlice("app.bool")
	t.Logf("app.bool: %v", ss64)
	ss := trie.MustBoolSlice("app.bool-absent", true, false)
	t.Logf("app.bool-absent: %v", ss)
	ss, err = trie.GetBoolSlice("app.bool-absent", true, false)
	t.Logf("app.bool-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetBoolMap(t *testing.T) {
	trie := newTrieTree()

	trie.Set("app.bool", "{a:on,  b:n,c:true}")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetBoolMap("app.bool")
	t.Logf("app.bool: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, map[string]bool{"a": true, "b": false, "c": true}) {
		t.Fail()
	}
	ss64 = trie.MustBoolMap("app.bool")
	t.Logf("app.bool: %v", ss64)
	ss := trie.MustBoolMap("app.bool-absent", map[string]bool{"a": true})
	t.Logf("app.bool-absent: %v", ss)
	ss, err = trie.GetBoolMap("app.bool-absent", map[string]bool{"a": true})
	t.Logf("app.bool-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetDuration(t *testing.T) {
	trie := newTrieTree()

	const defval = (0*24+1)*time.Hour + 8*time.Millisecond + 719*time.Microsecond + 76*time.Nanosecond

	trie.Set("app.duration", "1h8ms719us76ns")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetDuration("app.duration")
	t.Logf("app.duration: %v", times.SmartDurationString(ss64))
	if err != nil || !reflect.DeepEqual(ss64, defval) {
		t.Fail()
	}
	ss64 = trie.MustDuration("app.duration")
	t.Logf("app.duration: %v", times.SmartDurationString(ss64))
	ss := trie.MustDuration("app.duration-absent", defval)
	t.Logf("app.duration-absent: %v", times.SmartDurationString(ss))
	ss, err = trie.GetDuration("app.duration-absent", defval)
	t.Logf("app.duration-absent: %v, err: %v", times.SmartDurationString(ss), err)
}

func TestTrieS_GetDurationSlice(t *testing.T) {
	trie := newTrieTree()

	const defval = (0*24+1)*time.Hour + 8*time.Millisecond + 719*time.Microsecond + 76*time.Nanosecond

	trie.Set("app.duration", "[1h8ms719us76ns,0s]")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetDurationSlice("app.duration")
	t.Logf("app.duration: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, []time.Duration{defval, 0}) {
		t.Fail()
	}
	ss64 = trie.MustDurationSlice("app.duration")
	t.Logf("app.duration: %v", ss64)
	ss := trie.MustDurationSlice("app.duration-absent", defval, time.Microsecond)
	t.Logf("app.duration-absent: %v", ss)
	ss, err = trie.GetDurationSlice("app.duration-absent", defval, time.Microsecond)
	t.Logf("app.duration-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetDurationMap(t *testing.T) {
	trie := newTrieTree()

	const defval = (0*24+1)*time.Hour + 8*time.Millisecond + 719*time.Microsecond + 76*time.Nanosecond

	trie.Set("app.duration", "{a:1h8ms719us76ns}")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetDurationMap("app.duration")
	t.Logf("app.duration: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, map[string]time.Duration{"a": defval}) {
		t.Fail()
	}
	ss64 = trie.MustDurationMap("app.duration")
	t.Logf("app.duration: %v", ss64)
	ss := trie.MustDurationMap("app.duration-absent", map[string]time.Duration{"a": defval})
	t.Logf("app.duration-absent: %v", ss)
	ss, err = trie.GetDurationMap("app.duration-absent", map[string]time.Duration{"a": defval})
	t.Logf("app.duration-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetTime(t *testing.T) {
	trie := newTrieTree()

	defval := times.MustSmartParseTime("2021-1-29")

	trie.Set("app.time", "2021-1-29")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetTime("app.time")
	t.Logf("app.time: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, defval) {
		t.Fail()
	}
	ss64 = trie.MustTime("app.time")
	t.Logf("app.time: %v", ss64)
	ss := trie.MustTime("app.time-absent", defval)
	t.Logf("app.time-absent: %v", ss)
	ss, err = trie.GetTime("app.time-absent", defval)
	t.Logf("app.time-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetTimeSlice(t *testing.T) {
	trie := newTrieTree()

	defval := times.MustSmartParseTime("2021-1-29")
	defval2 := times.MustSmartParseTime("1979-1-29")

	trie.Set("app.time", "[2021-1-29,1979-01-29]")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetTimeSlice("app.time")
	t.Logf("app.time: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, []time.Time{defval, defval2}) {
		t.Fail()
	}
	ss64 = trie.MustTimeSlice("app.time")
	t.Logf("app.time: %v", ss64)
	ss := trie.MustTimeSlice("app.time-absent", defval, defval2)
	t.Logf("app.time-absent: %v", ss)
	ss, err = trie.GetTimeSlice("app.time-absent", defval, defval2)
	t.Logf("app.time-absent: %v, err: %v", ss, err)
}

func TestTrieS_GetTimeMap(t *testing.T) {
	trie := newTrieTree()

	defval := times.MustSmartParseTime("2021-1-29")
	defval2 := times.MustSmartParseTime("1979-1-29")
	mdef := map[string]time.Time{"a": defval, "b": defval2}

	trie.Set("app.time", "{a:2021-1-29,b:1979-01-29}")
	ret := trie.dump(false)
	t.Logf("\nPath\n%v\n", ret)

	ss64, err := trie.GetTimeMap("app.time")
	t.Logf("app.time: %v", ss64)
	if err != nil || !reflect.DeepEqual(ss64, mdef) {
		t.Fail()
	}
	ss64 = trie.MustTimeMap("app.time")
	t.Logf("app.time: %v", ss64)
	ss := trie.MustTimeMap("app.time-absent", mdef)
	t.Logf("app.time-absent: %v", ss)
	ss, err = trie.GetTimeMap("app.time-absent", mdef)
	t.Logf("app.time-absent: %v, err: %v", ss, err)
}
