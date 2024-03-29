package ctx

import (
	"fmt"
	"testing"
	"time"
)

func TestCtxS_Next_Short(t *testing.T) {
	c := TODO() //nolint:gocritic
	for c.Next() {
		t.Log(c.Key())
	}
}

func TestCtxS_Next_Simple(t *testing.T) {
	c := WithValue(TODO(), "k1", 1)
	c.WithValues("k2", 2, "k3", 3)

	for c.Next() {
		t.Log(c.Key())
	}
}

func TestCtxS_NamesCount(t *testing.T) {
	ctx := WithValues(nil, "k1", 1, "k2", 2, "k3", time.Now(), time.Now(), "now").(*ctxS) //nolint:revive
	if ctx.NamesCount() != 4 {
		t.Fail()
	}

	for {
		n := ctx.NextName()
		if n == "" {
			break
		}
	}

	ctx.Reset()

	t.Log(ctx.Value()) // should be nil

	for ctx.Next() {
		t.Log(ctx.Entry())
	}

	ctx = WithValue(nil, "k1", 1).(*ctxS) //nolint:revive
	if ctx.NamesCount() != 1 {
		t.Fail()
	}
}

func TestCtxS_Next(t *testing.T) {
	ctx := WithValues(TODO(), "k1", 1, "k2", 2, "k3", 3).(*ctxS) //nolint:revive
	for ctx.Next() {
		t.Log(ctx.Key(), ctx.Value())
	}

	ctx.Reset()
	for ctx.Next() {
		t.Log(ctx.Key())
	}
}

func TestCtxS_ValueBy(t *testing.T) {
	ctx := WithValue(TODO(), "k1", 1).(*ctxS) //nolint:revive
	v := ctx.ValueBy("k1")
	if v != 1 {
		t.Fatalf("want 1 but got %v", v)
	}

	ctx.Reset()
	v = ctx.ValueBy("k1")
	if v != 1 {
		t.Fatalf("want 1 but got %v", v)
	}

	v = ctx.ValueBy("k2")
	if v != nil {
		t.Fatalf("want nil but got %v", v)
	}
}

func ExamplectxS_Next() {
	ctx := WithValue(TODO(), "k1", 1).(*ctxS) //nolint:revive
	ctx.add("k2", 2)
	ctx.add("k3", 3)

	for ctx.Next() {
		_, _ = fmt.Println(ctx.Key())
	}

	// Output:
	// k1
	// k2
	// k3
}

func ExampleWithValues() {
	ctx := WithValues(TODO(), "k1", 1, "k2", 2, "k3", 3)
	for ctx.Next() {
		_, _ = fmt.Println(ctx.Key())
	}

	// Output:
	// k1
	// k2
	// k3
}
