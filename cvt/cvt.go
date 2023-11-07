package cvt

import (
	"fmt"
	"strings"

	evendeep "github.com/hedzr/go-diff/v2"
)

func Copy(m map[string]any) (data map[string]any) {
	data = make(map[string]any)
	evendeep.Copy(&m, &data)
	return
}

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
		for k, v := range z {
			x[fmt.Sprintf("%v", k)] = v
		}
		m[key] = x
		Normalize(x, worker)
	case []any: // process json array structure
		for i, v := range z {
			switch sub := v.(type) {
			case map[any]any:
				x := make(map[string]any)
				for k, v := range sub {
					x[fmt.Sprintf("%v", k)] = v
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

func Deflate(m map[string]any, delim string) (data map[string]any) {
	data = make(map[string]any)

	// Iterate through the flat conf map.
	for k, v := range m {
		var (
			keys = strings.Split(k, delim)
			next = data
		)

		// Iterate through key parts, for eg:, parent.child.key
		// will be ["parent", "child", "key"]
		for _, k := range keys[:len(keys)-1] {
			sub, ok := next[k]
			if !ok {
				// If the key does not exist in the map, create it.
				sub = make(map[string]interface{})
				next[k] = sub
			}
			if n, ok := sub.(map[string]any); ok {
				next = n
			}
		}

		// Assign the value.
		next[keys[len(keys)-1]] = v
	}
	return
}
