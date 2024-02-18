package radix

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/evendeep"

	"gopkg.in/yaml.v3"
)

// TypedGetters makes a formal specification for Trie[any]
type TypedGetters[T any] interface {
	GetString(path string, defaultVal ...string) (ret string, err error)
	MustString(path string, defaultVal ...string) (ret string)
	GetStringSlice(path string, defaultVal ...string) (ret []string, err error)
	MustStringSlice(path string, defaultVal ...string) (ret []string)

	// GetStringMap locates a node and returns its data as a string map.
	//
	// Note that it doesn't care about the sub-nodes and these children's data.
	// If you want to extract a nodes tree from a given node, using GetM, or GetR.
	GetStringMap(path string, defaultVal ...map[string]string) (ret map[string]string, err error)
	MustStringMap(path string, defaultVal ...map[string]string) (ret map[string]string)

	TypedBooleanGetters
	TypedIntegersGetters
	MoreIntegersGetters
	TypedFloatsGetters
	TypedComplexesGetters
	TypedTimeGetters

	// GetR finds a given path recursively, and returns
	// the matched subtree as a map, which key is a dotted key path.
	//
	// See MustR for a sample result.
	//
	// GetR / MustR returns the whole tree when you gave the path "".
	GetR(path string, defaultVal ...map[string]any) (ret map[string]any, err error)
	// MustR finds a given path in tree recursively, and returns
	// the matched subtree as a map, which key is a dotted key path.
	//
	// A subtree will be flattened to dotted-key-path and value pairs.
	//
	// A sample is like (dumped by spew printer):
	//
	//     map[string]any{
	//       "app.debug": bool(false),
	//       "app.dump": int(3),
	//       "app.dump.to": "stdout",
	//       "app.logging.file": "/tmp/1.log",
	//       "app.logging.rotate": int(6),
	//       "app.logging.words": []string{
	//         "a",
	//         "1",
	//         "false",
	//       },
	//       "app.server.start": int(5),
	//       "app.verbose": bool(true),
	//     }
	//
	// Note the `GetR("app")` return the same result.
	//
	// GetR / MustR returns the whole tree when you gave the path "".
	MustR(path string, defaultVal ...map[string]any) (ret map[string]any)

	// GetM finds a given path recursively, and returns the matched
	// subtree as a map, which keys keep the original hierarchical
	// structure.
	//
	// GetM("") will return the whole tree.
	//
	// The optional MOpt operators could be:
	//  - WithKeepPrefix
	//  - WithFilter
	//  - WithoutFlattenKeys
	GetM(path string, opt ...MOpt[T]) (ret map[string]any, err error)
	// MustM finds a given path recursively, and returns the matched
	// subtree as a map, which keys keep the original hierarchical
	// structure.
	//
	// MustM("") will return the whole tree.
	//
	// The optional MOpt operators could be:
	//  - WithKeepPrefix
	//  - WithFilter
	//  - WithoutFlattenKeys
	//
	// MustM("app.logging") returns a subtree like:
	//
	//     map[string]any{
	//       "file": "/tmp/1.log",
	//       "rotate": int(6),
	//       "words": []string{
	//         "a",
	//         "1",
	//         "false",
	//       },
	//     }
	//
	// Note the key is without prefix 'app.logging.'.
	//
	// If you want to extract subtree with full prefixed keys, using
	// filter is a good idea:
	//
	//    m, err := trie.GetM("", WithFilter[any](func(node Node[any]){
	//       return strings.HasPrefix(node.Key(), "app.logging.")
	//    }))
	MustM(path string, opt ...MOpt[T]) (ret map[string]any)
	// GetSectionFrom finds a given path and loads the subtree into
	// 'holder', typically 'holder' could be a struct.
	//
	// For yaml input
	//
	//    app:
	//      server:
	//        sites:
	//          - name: default
	//            addr: ":7999"
	//            location: ~/Downloads/w/docs
	//
	// The following codes can load it into sitesS struct:
	//
	//	var sites sitesS
	//	err = store.WithPrefix("app").GetSectionFrom("server.sites", &sites)
	//
	//	type sitesS struct{ Sites []siteS }
	//
	//	type siteS struct {
	//	  Name        string
	//	  Addr        string
	//	  Location    string
	//	}
	//
	// In this above case, 'store' loaded yaml and built it
	// into memory, and extract 'server.sites' into 'sitesS'.
	// Since 'server.sites' is a yaml array, it was loaded
	// as a store entry and holds a slice value, so GetSectionFrom
	// extract it to sitesS.Sites field.
	//
	// The optional MOpt operators could be:
	//  - WithKeepPrefix
	//  - WithFilter
	GetSectionFrom(path string, holder any, opts ...MOpt[T]) (err error)
}

type TypedBooleanGetters interface {
	GetBool(path string, defaultVal ...bool) (ret bool, err error)
	MustBool(path string, defaultVal ...bool) (ret bool)
	GetBoolSlice(path string, defaultVal ...bool) (ret []bool, err error)
	MustBoolSlice(path string, defaultVal ...bool) (ret []bool)
	GetBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool, err error)
	MustBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool)
}

type TypedIntegersGetters interface {
	GetInt64(path string, defaultVal ...int64) (ret int64, err error)
	MustInt64(path string, defaultVal ...int64) (ret int64)
	GetInt32(path string, defaultVal ...int32) (ret int32, err error)
	MustInt32(path string, defaultVal ...int32) (ret int32)
	GetInt16(path string, defaultVal ...int16) (ret int16, err error)
	MustInt16(path string, defaultVal ...int16) (ret int16)
	GetInt8(path string, defaultVal ...int8) (ret int8, err error)
	MustInt8(path string, defaultVal ...int8) (ret int8)
	GetInt(path string, defaultVal ...int) (ret int, err error)
	MustInt(path string, defaultVal ...int) (ret int)

	GetUint64(path string, defaultVal ...uint64) (ret uint64, err error)
	MustUint64(path string, defaultVal ...uint64) (ret uint64)
	GetUint32(path string, defaultVal ...uint32) (ret uint32, err error)
	MustUint32(path string, defaultVal ...uint32) (ret uint32)
	GetUint16(path string, defaultVal ...uint16) (ret uint16, err error)
	MustUint16(path string, defaultVal ...uint16) (ret uint16)
	GetUint8(path string, defaultVal ...uint8) (ret uint8, err error)
	MustUint8(path string, defaultVal ...uint8) (ret uint8)
	GetUint(path string, defaultVal ...uint) (ret uint, err error)
	MustUint(path string, defaultVal ...uint) (ret uint)

	GetInt64Slice(path string, defaultVal ...int64) (ret []int64, err error)
	MustInt64Slice(path string, defaultVal ...int64) (ret []int64)
	GetInt32Slice(path string, defaultVal ...int32) (ret []int32, err error)
	MustInt32Slice(path string, defaultVal ...int32) (ret []int32)
	GetInt16Slice(path string, defaultVal ...int16) (ret []int16, err error)
	MustInt16Slice(path string, defaultVal ...int16) (ret []int16)
	GetInt8Slice(path string, defaultVal ...int8) (ret []int8, err error)
	MustInt8Slice(path string, defaultVal ...int8) (ret []int8)
	GetIntSlice(path string, defaultVal ...int) (ret []int, err error)
	MustIntSlice(path string, defaultVal ...int) (ret []int)

	GetUint64Slice(path string, defaultVal ...uint64) (ret []uint64, err error)
	MustUint64Slice(path string, defaultVal ...uint64) (ret []uint64)
	GetUint32Slice(path string, defaultVal ...uint32) (ret []uint32, err error)
	MustUint32Slice(path string, defaultVal ...uint32) (ret []uint32)
	GetUint16Slice(path string, defaultVal ...uint16) (ret []uint16, err error)
	MustUint16Slice(path string, defaultVal ...uint16) (ret []uint16)
	GetUint8Slice(path string, defaultVal ...uint8) (ret []uint8, err error)
	MustUint8Slice(path string, defaultVal ...uint8) (ret []uint8)
	GetUintSlice(path string, defaultVal ...uint) (ret []uint, err error)
	MustUintSlice(path string, defaultVal ...uint) (ret []uint)

	GetInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64, err error)
	MustInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64)
	GetInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32, err error)
	MustInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32)
	GetInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16, err error)
	MustInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16)
	GetInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8, err error)
	MustInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8)
	GetIntMap(path string, defaultVal ...map[string]int) (ret map[string]int, err error)
	MustIntMap(path string, defaultVal ...map[string]int) (ret map[string]int)

	GetUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64, err error)
	MustUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64)
	GetUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32, err error)
	MustUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32)
	GetUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16, err error)
	MustUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16)
	GetUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8, err error)
	MustUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8)
	GetUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint, err error)
	MustUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint)
}

type MoreIntegersGetters interface {
	GetKibiBytes(key string, defaultVal ...uint64) (ret uint64, err error)
	MustKibibytes(key string, defaultVal ...uint64) (ret uint64)
	GetKiloBytes(key string, defaultVal ...uint64) (ret uint64, err error)
	MustKilobytes(key string, defaultVal ...uint64) (ret uint64)
}

type TypedFloatsGetters interface {
	GetFloat64(path string, defaultVal ...float64) (ret float64, err error)
	MustFloat64(path string, defaultVal ...float64) (ret float64)
	GetFloat32(path string, defaultVal ...float32) (ret float32, err error)
	MustFloat32(path string, defaultVal ...float32) (ret float32)
	GetFloat64Slice(path string, defaultVal ...float64) (ret []float64, err error)
	MustFloat64Slice(path string, defaultVal ...float64) (ret []float64)
	GetFloat32Slice(path string, defaultVal ...float32) (ret []float32, err error)
	MustFloat32Slice(path string, defaultVal ...float32) (ret []float32)
	GetFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64, err error)
	MustFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64)
	GetFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32, err error)
	MustFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32)
}

type TypedComplexesGetters interface {
	GetComplex128(path string, defaultVal ...complex128) (ret complex128, err error)
	MustComplex128(path string, defaultVal ...complex128) (ret complex128)
	GetComplex64(path string, defaultVal ...complex64) (ret complex64, err error)
	MustComplex64(path string, defaultVal ...complex64) (ret complex64)
	GetComplex128Slice(path string, defaultVal ...complex128) (ret []complex128, err error)
	MustComplex128Slice(path string, defaultVal ...complex128) (ret []complex128)
	GetComplex64Slice(path string, defaultVal ...complex64) (ret []complex64, err error)
	MustComplex64Slice(path string, defaultVal ...complex64) (ret []complex64)
	GetComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128, err error)
	MustComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128)
	GetComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64, err error)
	MustComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64)
}

type TypedTimeGetters interface {
	GetDuration(path string, defaultVal ...time.Duration) (ret time.Duration, err error)
	MustDuration(path string, defaultVal ...time.Duration) (ret time.Duration)
	GetDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration, err error)
	MustDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration)
	GetDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration, err error)
	MustDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration)

	GetTime(path string, defaultVal ...time.Time) (ret time.Time, err error)
	MustTime(path string, defaultVal ...time.Time) (ret time.Time)
	GetTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time, err error)
	MustTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time)
	GetTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time, err error)
	MustTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time)
}

var _ TypedGetters[any] = (*trieS[any])(nil) // assertion helper

var _ Trie[any] = (*trieS[any])(nil) // assertion helper

var converter = evendeep.Cvt{}

func (s *trieS[T]) GetR(path string, defaultVal ...map[string]any) (ret map[string]any, err error) { //nolint:revive
	var (
		found, partialMatched, branch bool
		node                          *nodeS[T]
	)

	if path == "" || path == "." || path == "(root)" {
		// node, _, _, _ = s.Locate(path)
		// // if len(node.children) > 0 {
		// // 	node = node.children[0]
		// // }
		// ret = make(map[string]any)
		// ret[node.pathS] = node.data

		ret = make(map[string]any)
		s.root.Walk(func(prefix, key string, node Node[T]) {
			if (path == "" || !s.endsWith(prefix, s.delimiter)) && !node.isBranch() {
				ret[prefix] = node.Data()
			}
		})
		return
	}

	node, branch, partialMatched, found = s.Locate(path)
	if found || partialMatched {
		_, _, ret = branch, partialMatched, make(map[string]any)

		node.Walk(func(prefix, key string, node Node[T]) {
			if !s.endsWith(prefix, s.delimiter) && !node.isBranch() {
				// For a trie like:
				//
				//     app.                          <B>
				//       d                           <B>
				//         ebug                      <L>, app.debug => false
				//         ump                       <L>, app.dump => 3
				//           .to                     <L>, app.dump.to => stdout
				//       verbose                     <L>, app.verbose => true
				//       logging.                    <B>
				//         file                      <L>, app.logging.file => /tmp/1.log
				//         rotate                    <L>, app.logging.rotate => 6
				//         words                     <L>, app.logging.words => [a 1 false]
				//       server.start                <L>, app.server.start => 5
				//
				// the node 'app.' is a branch and ends with '.',
				// the node 'app.' / 'd' is a branch and not ends with '.'.
				//
				// We keep those nodes which is leaf and not ends with '.', just
				// like 'app.dump' or 'app.dump.to'.
				//
				// See also TestStore_GetR()
				ret[prefix] = node.Data()
			}
		})
		logz.Debug("[GetR] ", "ret", ret)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustR(path string, defaultVal ...map[string]any) (ret map[string]any) {
	ret, _ = s.GetR(path, defaultVal...)
	return
}

//

func (s *trieS[T]) GetM(path string, opts ...MOpt[T]) (ret map[string]any, err error) { //nolint:revive
	var (
		found, partialMatched, branch bool
		node                          *nodeS[T]
	)

	if path == "" || path == "." || path == "(root)" {
		ret = make(map[string]any)
		putter := prefixPutter[T]{}
		for _, opt := range opts {
			opt(&putter)
		}
		s.root.Walk(func(prefix, key string, node Node[T]) {
			if (path == "" || !s.endsWith(prefix, s.delimiter)) && !node.isBranch() {
				if putter.filterFn != nil {
					if !putter.filterFn(node) {
						return
					}
				}
				ret[prefix] = node.Data()
			}
		})
		if putter.noFlatten {
			ret = s.splitCompactKeys(ret)
		}
		return
	}

	node, branch, partialMatched, found = s.Locate(path)
	if found || partialMatched {
		_, _, ret = branch, partialMatched, make(map[string]any)
		putter := prefixPutter[T]{prefix: strings.Split(s.join(s.prefix, path), string(s.delimiter))}
		for _, opt := range opts {
			opt(&putter)
		}
		logz.Verbose("[GetM] loop subtree and return as a map", "path", putter.prefix)
		node.Walk(func(prefix, key string, node Node[T]) {
			if !s.endsWith(prefix, s.delimiter) && !node.isBranch() {
				logz.Verbose("  - put into map", "prefix", prefix, "key", key)
				putter.put(ret, prefix, string(s.delimiter), node.Data())
			}
		})
		logz.Verbose("[GetM] ", "ret", ret)
		if putter.noFlatten {
			ret = s.splitCompactKeys(ret)
		}
	}
	return
}

func (s *trieS[T]) MustM(path string, opts ...MOpt[T]) (ret map[string]any) {
	ret, _ = s.GetM(path, opts...)
	return
}

func (s *trieS[T]) splitCompactKeys(in map[string]any) (out map[string]any) {
	out = make(map[string]any)
	for k, v := range in {
		if strings.ContainsRune(k, s.delimiter) {
			a := strings.Split(k, string(s.delimiter))
			s.submap(out, a, v)
		} else {
			out[k] = v
		}
	}
	return
}

func (s *trieS[T]) submap(src map[string]any, keys []string, v any) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		src[keys[0]] = v
		return
	}

	k, rest := keys[0], keys[1:]
	if m1, ok := src[k]; ok {
		if m2, ok := m1.(map[string]any); ok {
			s.submap(m2, rest, v)
		}
	} else {
		m2 := make(map[string]any)
		s.submap(m2, rest, v)
		src[k] = m2
	}
}

// func (s *trieS[T]) mergemap(m map[string]any, key string, v map[string]any) {}

// func runeToString(runes ...rune) string { return string(runes) }

func (s *trieS[T]) GetSectionFrom(path string, holder any, opts ...MOpt[T]) (err error) {
	var ret map[string]any
	ret, err = s.GetM(path, opts...)
	if err == nil && ret != nil {
		defer handleSerializeError(&err)
		m := s.splitCompactKeys(ret)
		var b []byte
		b, err = yaml.Marshal(m)
		if err == nil {
			err = yaml.Unmarshal(b, holder)
			// if err == nil {
			// 	logrus.Debugf("configuration section got: %v", configHolder)
			// }
		}
	} else {
		ret, err = s.GetR(path)
		if err == nil && ret != nil {
			defer handleSerializeError(&err)
			m := s.splitCompactKeys(ret)
			var b []byte
			b, err = yaml.Marshal(m)
			if err == nil {
				err = yaml.Unmarshal(b, holder)
			}
		}
	}
	return
}

func handleSerializeError(err *error) { //nolint:gocritic //can't opt
	if v := recover(); v != nil {
		if e1, ok := v.(error); ok {
			*err = e1
		} else {
			*err = fmt.Errorf("%v", v)
		}
		// *err = v.(error) //nolint:errcheck    // errors.New("unexpected unknown error handled").WithData(v)
		// switch e := v.(type) {
		// case error:
		// 	*err = e
		// case string:
		// 	*err = errors.New(e)
		// // if s == "cannot marshal type: complex128" {
		// // 	err = errors.New(e)
		// // }
		// default:
		// 	panic(v)
		// }
	}
}

// WithKeepPrefix _
func WithKeepPrefix[T any](b bool) MOpt[T] {
	return func(s *prefixPutter[T]) {
		s.keepPrefix = b
	}
}

// WithFilter can be used in calling nodeS[T].GetM(path, ...)
func WithFilter[T any](filter FilterFn[T]) MOpt[T] {
	return func(s *prefixPutter[T]) {
		s.filterFn = filter
	}
}

// WithoutFlattenKeys allows returns a nested map.
// If the keys contain delimiter char, they will be split as
// nested sub-map.
func WithoutFlattenKeys[T any](b bool) MOpt[T] {
	return func(s *prefixPutter[T]) {
		s.noFlatten = b
	}
}

type FilterFn[T any] func(node Node[T]) bool // used by GetM, MustM, ...

type MOpt[T any] func(s *prefixPutter[T]) // used by GetM, MustM, ...

type prefixPutter[T any] struct {
	prefix     []string
	keepPrefix bool // constructing the result map by keeping prefix structure?
	noFlatten  bool // split key like 'app.logging.files' as nested sub-map
	filterFn   FilterFn[T]
}

func (s *prefixPutter[T]) put(m map[string]any, prefix, delimiter string, v any) {
	keys := strings.Split(prefix, delimiter)
	if s.keepPrefix {
		s.putKeys(m, keys, v)
		return
	}

	if len(keys) >= len(s.prefix) {
		kk := keys[len(s.prefix):]
		if len(kk) > 0 {
			s.putKeys(m, kk, v)
		} else {
			kk = keys[len(s.prefix)-1:]
			s.putKeys(m, kk, v)
		}
		return
	}
	s.putKeys(m, keys, v)
}

func (s *prefixPutter[T]) putKeys(m map[string]any, keys []string, v any) {
	if len(keys) == 0 {
		return
	}
	if len(keys) == 1 {
		m[keys[0]] = v
		return
	}

	k, rest := keys[0], keys[1:]
	if c, ok := m[k]; ok {
		if cm, ok := c.(map[string]any); ok {
			s.putKeys(cm, rest, v)
		} else { //nolint:staticcheck,revive
			// ? panic
		}
	} else {
		cm := make(map[string]any)
		m[k] = cm
		s.putKeys(cm, rest, v)
	}
}

//

func (s *trieS[T]) GetString(path string, defaultVal ...string) (ret string, err error) {
	var (
		found bool
		data  any
	)
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.String(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustString(path string, defaultVal ...string) (ret string) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.String(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetStringSlice(path string, defaultVal ...string) (ret []string, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.StringSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustStringSlice(path string, defaultVal ...string) (ret []string) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.StringSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetStringMap(path string, defaultVal ...map[string]string) (ret map[string]string, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.StringMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustStringMap(path string, defaultVal ...map[string]string) (ret map[string]string) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.StringMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetInt64(path string, defaultVal ...int64) (ret int64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt64(path string, defaultVal ...int64) (ret int64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt(path string, defaultVal ...int) (ret int, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = int(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt(path string, defaultVal ...int) (ret int) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = int(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt32(path string, defaultVal ...int32) (ret int32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = int32(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt32(path string, defaultVal ...int32) (ret int32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = int32(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt16(path string, defaultVal ...int16) (ret int16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = int16(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt16(path string, defaultVal ...int16) (ret int16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = int16(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt8(path string, defaultVal ...int8) (ret int8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = int8(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt8(path string, defaultVal ...int8) (ret int8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = int8(converter.Int(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint64(path string, defaultVal ...uint64) (ret uint64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint64(path string, defaultVal ...uint64) (ret uint64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint(path string, defaultVal ...uint) (ret uint, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = uint(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint(path string, defaultVal ...uint) (ret uint) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = uint(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint32(path string, defaultVal ...uint32) (ret uint32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = uint32(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint32(path string, defaultVal ...uint32) (ret uint32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = uint32(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint16(path string, defaultVal ...uint16) (ret uint16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = uint16(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint16(path string, defaultVal ...uint16) (ret uint16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = uint16(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint8(path string, defaultVal ...uint8) (ret uint8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = uint8(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint8(path string, defaultVal ...uint8) (ret uint8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = uint8(converter.Uint(data))
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetInt64Slice(path string, defaultVal ...int64) (ret []int64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustInt64Slice(path string, defaultVal ...int64) (ret []int64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetInt32Slice(path string, defaultVal ...int32) (ret []int32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustInt32Slice(path string, defaultVal ...int32) (ret []int32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetInt16Slice(path string, defaultVal ...int16) (ret []int16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int16Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustInt16Slice(path string, defaultVal ...int16) (ret []int16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int16Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetInt8Slice(path string, defaultVal ...int8) (ret []int8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int8Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustInt8Slice(path string, defaultVal ...int8) (ret []int8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int8Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetIntSlice(path string, defaultVal ...int) (ret []int, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.IntSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustIntSlice(path string, defaultVal ...int) (ret []int) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.IntSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

//

func (s *trieS[T]) GetUint64Slice(path string, defaultVal ...uint64) (ret []uint64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustUint64Slice(path string, defaultVal ...uint64) (ret []uint64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetUint32Slice(path string, defaultVal ...uint32) (ret []uint32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustUint32Slice(path string, defaultVal ...uint32) (ret []uint32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetUint16Slice(path string, defaultVal ...uint16) (ret []uint16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint16Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustUint16Slice(path string, defaultVal ...uint16) (ret []uint16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint16Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetUint8Slice(path string, defaultVal ...uint8) (ret []uint8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint8Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustUint8Slice(path string, defaultVal ...uint8) (ret []uint8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint8Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetUintSlice(path string, defaultVal ...uint) (ret []uint, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.UintSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustUintSlice(path string, defaultVal ...uint) (ret []uint) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.UintSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

//

func (s *trieS[T]) GetInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int16Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int16Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Int8Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Int8Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetIntMap(path string, defaultVal ...map[string]int) (ret map[string]int, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.IntMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustIntMap(path string, defaultVal ...map[string]int) (ret map[string]int) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.IntMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint16Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint16Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Uint8Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Uint8Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.UintMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.UintMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

// GetKibiBytes returns the uint64 value which is based kibi-byte format.
//
// kibi-byte format is human-readable.
//
// Within this format, number presentations are: 2k, 8m, 3g, 5t, 6p, 7e.
// An optional 'iB' can be appended to the unit suffix, such as 2KiB, 5TiB, 7EiB.
//
// Note that they are case-sensitive.
//
// kibi-byte is based 1024. It means:
//
//	1 KiB = 1k = 1024 bytes
//
// See also: https://en.wikipedia.org/wiki/Kibibyte
// Its related word is kilobyte, refer to: https://en.wikipedia.org/wiki/Kilobyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *trieS[T]) GetKibiBytes(key string, defaultVal ...uint64) (ir64 uint64, err error) {
	var sz string
	sz, err = s.GetString(key, "")
	if sz == "" || err != nil {
		for _, v := range defaultVal {
			ir64 = v
		}
		return
	}
	ir64 = s.FromKibiBytes(sz)
	return
}

func (s *trieS[T]) MustKibibytes(key string, defaultVal ...uint64) (ir64 uint64) {
	ir64, _ = s.GetKibiBytes(key, defaultVal...)
	return
}

// FromKibiBytes convert string to the uint64 value based kibi-byte format.
func (s *trieS[T]) FromKibiBytes(sz string) (ir64 uint64) {
	// var suffixes = []string {"B","KB","MB","GB","TB","PB","EB","ZB","YB"}
	const suffix = "kmgtpezyKMGTPEZY"
	sz = strings.TrimSpace(sz)       //nolint:revive
	sz = strings.TrimRight(sz, "iB") //nolint:revive
	sz = strings.TrimRight(sz, "ib") //nolint:revive
	szr := strings.TrimSpace(strings.TrimRightFunc(sz, func(r rune) bool {
		return strings.ContainsRune(suffix, r)
	}))

	var if64 float64
	var err error
	if strings.ContainsRune(szr, '.') {
		if if64, err = strconv.ParseFloat(szr, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 = uint64(if64 * float64(s.fromKibiBytes(r)))
		}
	} else {
		if ir64, err = strconv.ParseUint(szr, 0, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 *= s.fromKibiBytes(r)
		}
	}
	return
}

func (s *trieS[T]) fromKibiBytes(r rune) (times uint64) { //nolint:revive
	switch r {
	case 'k', 'K':
		return 1024
	case 'm', 'M':
		return 1024 * 1024
	case 'g', 'G':
		return 1024 * 1024 * 1024
	case 't', 'T':
		return 1024 * 1024 * 1024 * 1024
	case 'p', 'P':
		return 1024 * 1024 * 1024 * 1024 * 1024
	case 'e', 'E':
		return 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	// case 'z', 'Z':
	// 	ir64 *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	// case 'y', 'Y':
	// 	ir64 *= 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024
	default:
		return 1
	}
}

//

// GetKiloBytes returns the uint64 value which is based kilo-byte format.
//
// kilo-byte format is human-readable.
//
// Within this format, number presentations are: 2K, 8M, 3G, 5T, 6P, 7E.
// An optional 'B' can be appended to the unit suffix, such as 2KB, 5TB, 7EB.
// All of them are case-insensitive.
//
// kilo-byte is based 1000. It means:
//
//	1 KB = 1k = 1000 bytes
//
// See also: https://en.wikipedia.org/wiki/Kilobyte
// Its related word is kilo-byte, refer to: https://en.wikipedia.org/wiki/Kilobyte
//
// The pure number part can be golang presentation, such as 0x99, 0001b, 0700.
func (s *trieS[T]) GetKiloBytes(key string, defaultVal ...uint64) (ir64 uint64, err error) {
	var sz string
	sz, err = s.GetString(key, "")
	if sz == "" || err != nil {
		for _, v := range defaultVal {
			ir64 = v
		}
		return
	}
	ir64 = s.FromKiloBytes(sz)
	return
}

func (s *trieS[T]) MustKilobytes(key string, defaultVal ...uint64) (ir64 uint64) {
	ir64, _ = s.GetKiloBytes(key, defaultVal...)
	return
}

// FromKiloBytes converts the uint64 value which is based kilo-byte format.
func (s *trieS[T]) FromKiloBytes(sz string) (ir64 uint64) {
	// var suffixes = []string {"B","KB","MB","GB","TB","PB","EB","ZB","YB"}
	const suffix = "kmgtpezyKMGTPEZY"
	sz = strings.TrimSpace(sz)      //nolint:revive
	sz = strings.TrimRight(sz, "B") //nolint:revive
	sz = strings.TrimRight(sz, "b") //nolint:revive
	szr := strings.TrimSpace(strings.TrimRightFunc(sz, func(r rune) bool {
		return strings.ContainsRune(suffix, r)
	}))

	var if64 float64
	var err error
	if strings.ContainsRune(szr, '.') {
		if if64, err = strconv.ParseFloat(szr, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 = uint64(if64 * float64(s.fromKilobytes(r)))
		}
	} else {
		if ir64, err = strconv.ParseUint(szr, 0, 64); err == nil {
			r := []rune(sz)[len(sz)-1]
			ir64 *= s.fromKilobytes(r)
		}
	}
	return
}

func (s *trieS[T]) fromKilobytes(r rune) (times uint64) { //nolint:revive
	switch r {
	case 'k', 'K':
		return 1000
	case 'm', 'M':
		return 1000 * 1000
	case 'g', 'G':
		return 1000 * 1000 * 1000
	case 't', 'T':
		return 1000 * 1000 * 1000 * 1000
	case 'p', 'P':
		return 1000 * 1000 * 1000 * 1000 * 1000
	case 'e', 'E':
		return 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	// case 'z', 'Z':
	// 	ir64 *= 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	// case 'y', 'Y':
	// 	ir64 *= 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000
	default:
		return 1
	}
}

//

func (s *trieS[T]) GetFloat64(path string, defaultVal ...float64) (ret float64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float64(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustFloat64(path string, defaultVal ...float64) (ret float64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float64(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetFloat32(path string, defaultVal ...float32) (ret float32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float32(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustFloat32(path string, defaultVal ...float32) (ret float32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float32(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetFloat64Slice(path string, defaultVal ...float64) (ret []float64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustFloat64Slice(path string, defaultVal ...float64) (ret []float64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetFloat32Slice(path string, defaultVal ...float32) (ret []float32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustFloat32Slice(path string, defaultVal ...float32) (ret []float32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float32Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Float32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Float32Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetComplex128(path string, defaultVal ...complex128) (ret complex128, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex128(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustComplex128(path string, defaultVal ...complex128) (ret complex128) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex128(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetComplex64(path string, defaultVal ...complex64) (ret complex64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex64(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustComplex64(path string, defaultVal ...complex64) (ret complex64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex64(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetComplex128Slice(path string, defaultVal ...complex128) (ret []complex128, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex128Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustComplex128Slice(path string, defaultVal ...complex128) (ret []complex128) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex128Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetComplex64Slice(path string, defaultVal ...complex64) (ret []complex64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustComplex64Slice(path string, defaultVal ...complex64) (ret []complex64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex64Slice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex128Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex128Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Complex64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Complex64Map(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetBool(path string, defaultVal ...bool) (ret bool, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Bool(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustBool(path string, defaultVal ...bool) (ret bool) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Bool(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetBoolSlice(path string, defaultVal ...bool) (ret []bool, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.BoolSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustBoolSlice(path string, defaultVal ...bool) (ret []bool) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.BoolSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.BoolMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.BoolMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetDuration(path string, defaultVal ...time.Duration) (ret time.Duration, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Duration(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustDuration(path string, defaultVal ...time.Duration) (ret time.Duration) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Duration(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.DurationSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.DurationSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.DurationMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.DurationMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

//

func (s *trieS[T]) GetTime(path string, defaultVal ...time.Time) (ret time.Time, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.Time(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustTime(path string, defaultVal ...time.Time) (ret time.Time) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.Time(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) GetTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.TimeSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) MustTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.TimeSlice(data)
	} else if !found {
		ret = defaultVal
	}
	return
}

func (s *trieS[T]) GetTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time, err error) {
	var found bool
	var data any
	data, _, found, err = s.Query(path)
	if found && err == nil {
		ret = converter.TimeMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}

func (s *trieS[T]) MustTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time) {
	data, _, found, err := s.Query(path)
	if found && err == nil {
		ret = converter.TimeMap(data)
	} else if !found {
		for _, v := range defaultVal {
			ret = v
		}
	}
	return
}
