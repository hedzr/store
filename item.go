package store

import (
	"fmt"
	"time"
)

// An Item is a key-value pair.
type Item struct {
	Key   string
	Value Value
}

// String returns an Item for a string value.
func String(key, value string) Item {
	return Item{key, StringValue(value)}
}

// Int64 returns an Item for an int64.
func Int64(key string, value int64) Item {
	return Item{key, Int64Value(value)}
}

// Int converts an int to an int64 and returns
// an Item with that value.
func Int(key string, value int) Item {
	return Int64(key, int64(value))
}

// Uint64 returns an Item for uint64.
func Uint64(key string, v uint64) Item {
	return Item{key, Uint64Value(v)}
}

// Uint returns an Item for uint64.
func Uint(key string, v uint) Item {
	return Item{key, UintValue(v)}
}

// Float64 returns an Item for a floating-point number.
func Float64(key string, v float64) Item {
	return Item{key, Float64Value(v)}
}

// Bool returns an Item for a bool.
func Bool(key string, v bool) Item {
	return Item{key, BoolValue(v)}
}

// Time returns an Item for a time.Time.
// It discards the monotonic portion.
func Time(key string, v time.Time) Item {
	return Item{key, TimeValue(v)}
}

// Duration returns an Item for a time.Duration.
func Duration(key string, v time.Duration) Item {
	return Item{key, DurationValue(v)}
}

// Group returns an Item for a Group Value.
// The first argument is the key; the remaining arguments
// are converted to Attrs as in [Logger.Log].
//
// Use Group to collect several key-value pairs under a single
// key on a log line, or as the result of StoreValue
// in order to log a single value as multiple Attrs.
func Group(key string, args ...any) Item {
	return Item{key, GroupValue(argsToAttrSlice(args)...)}
}

func argsToAttrSlice(args []any) []Item {
	var (
		attr  Item
		attrs []Item
	)
	for len(args) > 0 {
		attr, args = argsToAttr(args)
		attrs = append(attrs, attr)
	}
	return attrs
}

// Any returns an Item for the supplied value.
// See [AnyValue] for how values are treated.
func Any(key string, value any) Item {
	return Item{key, AnyValue(value)}
}

// Equal reports whether a and b have equal keys and values.
func (a Item) Equal(b Item) bool {
	return a.Key == b.Key && a.Value.Equal(b.Value)
}

func (a Item) String() string {
	return fmt.Sprintf("%s=%s", a.Key, a.Value)
}

// isEmpty reports whether a has an empty key and a nil value.
// That can be written as Item{} or Any("", nil).
func (a Item) isEmpty() bool {
	return a.Key == "" && a.Value.num == 0 && a.Value.any == nil
}

func (a Item) Kind() Kind {
	if !a.isEmpty() {
		return a.Value.Kind()
	}
	return KindAny
}
