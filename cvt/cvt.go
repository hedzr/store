package cvt

import (
	"fmt"
	"strings"

	"github.com/hedzr/evendeep"
)

// Copy copies a map deeply
func Copy(m map[string]any) (data map[string]any) {
	data = make(map[string]any)
	evendeep.Copy(&m, &data)
	return
}

// Normalize make a map formal with nested map[string]any objects.
func Normalize(m map[string]any, worker func(k string, v any)) map[string]any {
	for k, v := range m {
		normalize(m, k, v, worker)
	}
	return m
}

func normalize(m map[string]any, key string, val any, worker func(k string, v any)) {
	switch z := val.(type) {
	case map[any]any:
		x := make(map[string]any)
		for k, vv := range z {
			kk := fmt.Sprintf("%v", k)
			x[kk] = vv
			if worker != nil {
				worker(kk, vv)
			}
		}
		m[key] = x
		Normalize(x, worker)
	case []any: // process json array structure
		for i, v := range z {
			switch sub := v.(type) {
			case map[any]any:
				x := make(map[string]any)
				for k, vv := range sub {
					kk := fmt.Sprintf("%v", k)
					x[kk] = vv
					if worker != nil {
						worker(kk, vv)
					}
				}
				z[i] = x
				Normalize(x, worker)
			case map[string]any:
				Normalize(sub, worker)
			}
		}
	case map[string]any:
		Normalize(z, worker)
	}
}

// Deflate split dotted key as a sub-map, which dot char can be another delimiter.
func Deflate(m map[string]any, delimiter string) (data map[string]any) {
	data = make(map[string]any)

	// Iterate through the flat conf map.
	for k, v := range m {
		var (
			keys = strings.Split(k, delimiter)
			next = data
		)

		// Iterate through key parts, for eg:, parent.child.key
		// will be ["parent", "child", "key"]
		for _, kk := range keys[:len(keys)-1] {
			sub, ok := next[kk]
			if !ok {
				// If the key does not exist in the map, create it.
				sub = make(map[string]any)
				next[kk] = sub
			}
			if n, ok1 := sub.(map[string]any); ok1 {
				next = n
			}
		}

		// Assign the value.
		next[keys[len(keys)-1]] = v
	}
	return
}
