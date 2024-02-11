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
	}
	mm := Normalize(m, nil)
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
