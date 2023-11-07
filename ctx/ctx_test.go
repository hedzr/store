package ctx

import (
	"testing"
)

func TestCtxS_NamesCount(t *testing.T) {
	ctx := TODO()
	for ctx.Next() {
		t.Log(ctx.Key())
	}
}

func TestCtxS_Next(t *testing.T) {
	ctx := WithValue(TODO(), "k1", 1).(*ctxS)
	ctx.add("k2", 2)
	ctx.add("k3", 3)

	for ctx.Next() {
		t.Log(ctx.Key())
	}
}
