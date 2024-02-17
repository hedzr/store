package cvt

import (
	"reflect"
	"testing"
)

func TestCopy(t *testing.T) {
	m := map[string]any{
		"k1": map[any]any{
			1.3: 1.313,
		},
		"k2": []string{"s", "w"},
	}
	mm := Copy(m)
	t.Logf("result: %v", mm)
	if !reflect.DeepEqual(mm, m) {
		t.Fatalf("wrong. map mm = %v", mm)
	}
}

func TestNormalize(t *testing.T) {
	m := map[string]any{
		"k1": map[any]any{
			1.3: 1.313,
		},
		"k2": []string{"s", "w"},
		"k3": []any{nil, map[any]any{2.7: 2.718}, map[string]any{"e": 2.71828}},
		"k4": map[any]any{
			1.3: map[any]any{
				1.4: 1.34,
			},
			2.5: map[any]any{
				5: map[any]any{
					uint(5): int64(6),
				},
			},
		},
	}

	mm := Normalize(m, func(k string, v any) {
		t.Logf("[worker] %q => %v", k, v)
	})
	t.Logf("result: %v", mm)

	mm = Normalize(m, nil)
	t.Logf("result: %v", mm)
	if !reflect.DeepEqual(mm["k1"], map[string]any{"1.3": 1.313}) {
		t.Fatalf("wrong. k1 = %v", mm["k1"])
	}
}

func TestDeflate(t *testing.T) {
	m := map[string]any{
		"k1.k2": map[any]any{
			1.3: 1.313,
		},
		"k2": []string{"s", "w"},
	}
	mm := Deflate(m, ".")
	t.Logf("result: %v", mm)
	if !reflect.DeepEqual(mm["k1"].(map[string]any)["k2"], map[any]any{1.3: 1.313}) {
		t.Fatalf("wrong. k1 = %v", mm["k1"])
	}
}
