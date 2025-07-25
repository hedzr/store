package store

import (
	"context"
	"strings"
	"time"

	"github.com/hedzr/store/radix"
)

// NewDummyStore returns an empty store with dummy abilities implemented.
func NewDummyStore() *dummyS { return &dummyS{} }

type dummyS struct{ p string }

var _ radix.TypedGetters[any] = (*dummyS)(nil) // assertion helper

var _ Store = (*dummyS)(nil) // assertion helper

var _ MinimalStoreT[any] = (*dummyS)(nil) // assertion helper

func (s *dummyS) Close()                                                    {}
func (s *dummyS) MustGet(path string) (data any)                            { return }
func (s *dummyS) Get(path string) (data any, found bool)                    { return }
func (s *dummyS) Set(path string, data any) (node radix.Node[any], old any) { return }
func (s *dummyS) SetComment(path, description, comment string) (ok bool)    { return }
func (s *dummyS) SetTag(path string, tags any) (ok bool)                    { return } // set extra notable data bound to a key
func (s *dummyS) SetTTL(path string, ttl time.Duration, cb radix.OnTTLRinging[any]) (state int) {
	return
}
func (s *dummyS) SetTTLFast(node radix.Node[any], ttl time.Duration, cb radix.OnTTLRinging[any]) (state int) {
	return
}
func (s *dummyS) GetDesc(path string) (desc string, err error)       { return } // get tag field directly
func (s *dummyS) MustGetDesc(path string) (desc string)              { return } // mustget tag field directly
func (s *dummyS) GetTag(path string) (tag any, err error)            { return } // get tag field directly
func (s *dummyS) MustGetTag(path string) (tag any)                   { return } // mustget tag field directly
func (s *dummyS) GetComment(path string) (comment string, err error) { return } // get comment field directly
func (s *dummyS) MustGetComment(path string) (comment string)        { return } // mustget comment field directly
func (s *dummyS) GetEx(path string, cb func(node radix.Node[any], data any, branch bool, kvpair radix.KVPair)) {
}
func (s *dummyS) SetEx(path string, data any, cb radix.OnSetEx[any]) (old any)             { return }
func (s *dummyS) Remove(path string) (removed bool)                                        { return }
func (s *dummyS) RemoveEx(path string) (nodeRemoved, parent radix.Node[any], removed bool) { return }
func (s *dummyS) Merge(pathAt string, data map[string]any) (err error)                     { return }
func (s *dummyS) Has(path string) (found bool)                                             { return }
func (s *dummyS) Update(path string, cb func(node radix.Node[any], old any))               {}

// Locate provides an advanced interface for locating a path.
func (s *dummyS) Locate(path string, kvpair radix.KVPair) (node radix.Node[any], branch, partialMatched, found bool) {
	return
}

func (s *dummyS) GetString(path string, defaultVal ...string) (ret string, err error)        { return }
func (s *dummyS) MustString(path string, defaultVal ...string) (ret string)                  { return }
func (s *dummyS) GetStringSlice(path string, defaultVal ...string) (ret []string, err error) { return }
func (s *dummyS) MustStringSlice(path string, defaultVal ...string) (ret []string)           { return }
func (s *dummyS) GetStringMap(path string, defaultVal ...map[string]string) (ret map[string]string, err error) {
	return
}

func (s *dummyS) MustStringMap(path string, defaultVal ...map[string]string) (ret map[string]string) {
	return
}

func (s *dummyS) GetInt64(path string, defaultVal ...int64) (ret int64, err error)        { return }
func (s *dummyS) MustInt64(path string, defaultVal ...int64) (ret int64)                  { return }
func (s *dummyS) GetInt32(path string, defaultVal ...int32) (ret int32, err error)        { return }
func (s *dummyS) MustInt32(path string, defaultVal ...int32) (ret int32)                  { return }
func (s *dummyS) GetInt16(path string, defaultVal ...int16) (ret int16, err error)        { return }
func (s *dummyS) MustInt16(path string, defaultVal ...int16) (ret int16)                  { return }
func (s *dummyS) GetInt8(path string, defaultVal ...int8) (ret int8, err error)           { return }
func (s *dummyS) MustInt8(path string, defaultVal ...int8) (ret int8)                     { return }
func (s *dummyS) GetInt(path string, defaultVal ...int) (ret int, err error)              { return }
func (s *dummyS) MustInt(path string, defaultVal ...int) (ret int)                        { return }
func (s *dummyS) GetInt64Slice(path string, defaultVal ...int64) (ret []int64, err error) { return }
func (s *dummyS) MustInt64Slice(path string, defaultVal ...int64) (ret []int64)           { return }
func (s *dummyS) GetInt32Slice(path string, defaultVal ...int32) (ret []int32, err error) { return }
func (s *dummyS) MustInt32Slice(path string, defaultVal ...int32) (ret []int32)           { return }
func (s *dummyS) GetInt16Slice(path string, defaultVal ...int16) (ret []int16, err error) { return }
func (s *dummyS) MustInt16Slice(path string, defaultVal ...int16) (ret []int16)           { return }
func (s *dummyS) GetInt8Slice(path string, defaultVal ...int8) (ret []int8, err error)    { return }
func (s *dummyS) MustInt8Slice(path string, defaultVal ...int8) (ret []int8)              { return }
func (s *dummyS) GetIntSlice(path string, defaultVal ...int) (ret []int, err error)       { return }
func (s *dummyS) MustIntSlice(path string, defaultVal ...int) (ret []int)                 { return }

func (s *dummyS) GetInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64, err error) {
	return
}

func (s *dummyS) MustInt64Map(path string, defaultVal ...map[string]int64) (ret map[string]int64) {
	return
}

func (s *dummyS) GetInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32, err error) {
	return
}

func (s *dummyS) MustInt32Map(path string, defaultVal ...map[string]int32) (ret map[string]int32) {
	return
}

func (s *dummyS) GetInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16, err error) {
	return
}

func (s *dummyS) MustInt16Map(path string, defaultVal ...map[string]int16) (ret map[string]int16) {
	return
}

func (s *dummyS) GetInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8, err error) {
	return
}

func (s *dummyS) MustInt8Map(path string, defaultVal ...map[string]int8) (ret map[string]int8) {
	return
}

func (s *dummyS) GetIntMap(path string, defaultVal ...map[string]int) (ret map[string]int, err error) {
	return
}

func (s *dummyS) MustIntMap(path string, defaultVal ...map[string]int) (ret map[string]int) { return }

func (s *dummyS) GetUint64(path string, defaultVal ...uint64) (ret uint64, err error)        { return }
func (s *dummyS) MustUint64(path string, defaultVal ...uint64) (ret uint64)                  { return }
func (s *dummyS) GetUint32(path string, defaultVal ...uint32) (ret uint32, err error)        { return }
func (s *dummyS) MustUint32(path string, defaultVal ...uint32) (ret uint32)                  { return }
func (s *dummyS) GetUint16(path string, defaultVal ...uint16) (ret uint16, err error)        { return }
func (s *dummyS) MustUint16(path string, defaultVal ...uint16) (ret uint16)                  { return }
func (s *dummyS) GetUint8(path string, defaultVal ...uint8) (ret uint8, err error)           { return }
func (s *dummyS) MustUint8(path string, defaultVal ...uint8) (ret uint8)                     { return }
func (s *dummyS) GetUint(path string, defaultVal ...uint) (ret uint, err error)              { return }
func (s *dummyS) MustUint(path string, defaultVal ...uint) (ret uint)                        { return }
func (s *dummyS) GetUint64Slice(path string, defaultVal ...uint64) (ret []uint64, err error) { return }
func (s *dummyS) MustUint64Slice(path string, defaultVal ...uint64) (ret []uint64)           { return }
func (s *dummyS) GetUint32Slice(path string, defaultVal ...uint32) (ret []uint32, err error) { return }
func (s *dummyS) MustUint32Slice(path string, defaultVal ...uint32) (ret []uint32)           { return }
func (s *dummyS) GetUint16Slice(path string, defaultVal ...uint16) (ret []uint16, err error) { return }
func (s *dummyS) MustUint16Slice(path string, defaultVal ...uint16) (ret []uint16)           { return }
func (s *dummyS) GetUint8Slice(path string, defaultVal ...uint8) (ret []uint8, err error)    { return }
func (s *dummyS) MustUint8Slice(path string, defaultVal ...uint8) (ret []uint8)              { return }
func (s *dummyS) GetUintSlice(path string, defaultVal ...uint) (ret []uint, err error)       { return }
func (s *dummyS) MustUintSlice(path string, defaultVal ...uint) (ret []uint)                 { return }

func (s *dummyS) GetUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64, err error) {
	return
}

func (s *dummyS) MustUint64Map(path string, defaultVal ...map[string]uint64) (ret map[string]uint64) {
	return
}

func (s *dummyS) GetUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32, err error) {
	return
}

func (s *dummyS) MustUint32Map(path string, defaultVal ...map[string]uint32) (ret map[string]uint32) {
	return
}

func (s *dummyS) GetUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16, err error) {
	return
}

func (s *dummyS) MustUint16Map(path string, defaultVal ...map[string]uint16) (ret map[string]uint16) {
	return
}

func (s *dummyS) GetUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8, err error) {
	return
}

func (s *dummyS) MustUint8Map(path string, defaultVal ...map[string]uint8) (ret map[string]uint8) {
	return
}

func (s *dummyS) GetUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint, err error) {
	return
}

func (s *dummyS) MustUintMap(path string, defaultVal ...map[string]uint) (ret map[string]uint) {
	return
}

func (s *dummyS) GetKibiBytes(path string, defaultVal ...uint64) (ret uint64, err error) { return }
func (s *dummyS) MustKibiBytes(path string, defaultVal ...uint64) (ret uint64)           { return }
func (s *dummyS) GetKiloBytes(path string, defaultVal ...uint64) (ret uint64, err error) { return }
func (s *dummyS) MustKiloBytes(path string, defaultVal ...uint64) (ret uint64)           { return }

func (s *dummyS) GetFloat64(path string, defaultVal ...float64) (ret float64, err error) { return }
func (s *dummyS) MustFloat64(path string, defaultVal ...float64) (ret float64)           { return }
func (s *dummyS) GetFloat32(path string, defaultVal ...float32) (ret float32, err error) { return }
func (s *dummyS) MustFloat32(path string, defaultVal ...float32) (ret float32)           { return }

func (s *dummyS) GetFloat64Slice(path string, defaultVal ...float64) (ret []float64, err error) {
	return
}

func (s *dummyS) MustFloat64Slice(path string, defaultVal ...float64) (ret []float64) { return }

func (s *dummyS) GetFloat32Slice(path string, defaultVal ...float32) (ret []float32, err error) {
	return
}

func (s *dummyS) MustFloat32Slice(path string, defaultVal ...float32) (ret []float32) { return }

func (s *dummyS) GetFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64, err error) {
	return
}

func (s *dummyS) MustFloat64Map(path string, defaultVal ...map[string]float64) (ret map[string]float64) {
	return
}

func (s *dummyS) GetFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32, err error) {
	return
}

func (s *dummyS) MustFloat32Map(path string, defaultVal ...map[string]float32) (ret map[string]float32) {
	return
}

func (s *dummyS) GetComplex128(path string, defaultVal ...complex128) (ret complex128, err error) {
	return
}

func (s *dummyS) MustComplex128(path string, defaultVal ...complex128) (ret complex128) { return }

func (s *dummyS) GetComplex64(path string, defaultVal ...complex64) (ret complex64, err error) {
	return
}

func (s *dummyS) MustComplex64(path string, defaultVal ...complex64) (ret complex64) { return }

func (s *dummyS) GetComplex128Slice(path string, defaultVal ...complex128) (ret []complex128, err error) {
	return
}

func (s *dummyS) MustComplex128Slice(path string, defaultVal ...complex128) (ret []complex128) {
	return
}

func (s *dummyS) GetComplex64Slice(path string, defaultVal ...complex64) (ret []complex64, err error) {
	return
}

func (s *dummyS) MustComplex64Slice(path string, defaultVal ...complex64) (ret []complex64) { return }

func (s *dummyS) GetComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128, err error) {
	return
}

func (s *dummyS) MustComplex128Map(path string, defaultVal ...map[string]complex128) (ret map[string]complex128) {
	return
}

func (s *dummyS) GetComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64, err error) {
	return
}

func (s *dummyS) MustComplex64Map(path string, defaultVal ...map[string]complex64) (ret map[string]complex64) {
	return
}

func (s *dummyS) GetBool(path string, defaultVal ...bool) (ret bool, err error)        { return }
func (s *dummyS) MustBool(path string, defaultVal ...bool) (ret bool)                  { return }
func (s *dummyS) GetBoolSlice(path string, defaultVal ...bool) (ret []bool, err error) { return }
func (s *dummyS) MustBoolSlice(path string, defaultVal ...bool) (ret []bool)           { return }

func (s *dummyS) GetBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool, err error) {
	return
}

func (s *dummyS) MustBoolMap(path string, defaultVal ...map[string]bool) (ret map[string]bool) {
	return
}

func (s *dummyS) GetDuration(path string, defaultVal ...time.Duration) (ret time.Duration, err error) {
	return
}

func (s *dummyS) MustDuration(path string, defaultVal ...time.Duration) (ret time.Duration) { return }

func (s *dummyS) GetDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration, err error) {
	return
}

func (s *dummyS) MustDurationSlice(path string, defaultVal ...time.Duration) (ret []time.Duration) {
	return
}

func (s *dummyS) GetDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration, err error) {
	return
}

func (s *dummyS) MustDurationMap(path string, defaultVal ...map[string]time.Duration) (ret map[string]time.Duration) {
	return
}

func (s *dummyS) GetTime(path string, defaultVal ...time.Time) (ret time.Time, err error) { return }
func (s *dummyS) MustTime(path string, defaultVal ...time.Time) (ret time.Time)           { return }

func (s *dummyS) GetTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time, err error) {
	return
}

func (s *dummyS) MustTimeSlice(path string, defaultVal ...time.Time) (ret []time.Time) { return }

func (s *dummyS) GetTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time, err error) {
	return
}

func (s *dummyS) MustTimeMap(path string, defaultVal ...map[string]time.Time) (ret map[string]time.Time) {
	return
}

func (s *dummyS) GetR(path string, defaultVal ...map[string]any) (ret map[string]any, err error) {
	return
}

func (s *dummyS) MustR(path string, defaultVal ...map[string]any) (ret map[string]any) { return }

func (s *dummyS) GetM(path string, opt ...radix.MOpt[any]) (ret map[string]any, err error) {
	return
}

func (s *dummyS) MustM(path string, opt ...radix.MOpt[any]) (ret map[string]any)              { return }
func (s *dummyS) GetSectionFrom(path string, holder any, opts ...radix.MOpt[any]) (err error) { return }
func (s *dummyS) To(path string, target any, opts ...radix.MOpt[any]) (err error)             { return }

func (s *dummyS) Dump() (text string)                                                    { return }
func (s *dummyS) Clone() (newStore Store)                                                { return }
func (s *dummyS) Dup() (newStore Store)                                                  { return }
func (s *dummyS) Walk(path string, cb func(path, fragment string, node radix.Node[any])) {}
func (s *dummyS) WithPrefix(prefix ...string) (newStore Store)                           { return s }
func (s *dummyS) WithPrefixReplaced(prefix ...string) (newStore Store)                   { return s }
func (s *dummyS) SetPrefix(prefix ...string)                                             { s.p = strings.Join(prefix, ".") }
func (s *dummyS) Prefix() string                                                         { return s.p }
func (s *dummyS) Delimiter() rune                                                        { return '.' }
func (s *dummyS) SetDelimiter(delimiter rune)                                            {}
func (s *dummyS) Load(ctx context.Context, opts ...LoadOpt) (wr Writeable, err error)    { return }
func (s *dummyS) WithinLoading(fn func())                                                { fn() }

func (s *dummyS) SaveAs(ctx context.Context, file string, opts ...SaveAsOpt) (err error) { return }

func (s *dummyS) N() (newStore Store)               { return s }
func (s *dummyS) R() (newStore Store)               { return s }
func (s *dummyS) BR() (newStore Store)              { return s }
func (s *dummyS) RecursiveMode() radix.RecusiveMode { return radix.RecusiveNone }
