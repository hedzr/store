package ctx

import (
	"fmt"
	"sort"
)

// type Val *Value

// Val represents a value.
type Val any

// Ctx is like go context but has a stronger value api.
type Ctx interface {
	NamesCount() int
	Next() bool
	Entry() (key string, val Val)
	Key() string
	Value() Val

	WithValues(args ...any)
	WithValue(name string, value Val)

	// Reset resets iterating states so that you can re-iterate
	// the values via Next
	Reset()
}

func TODO() Ctx { return &ctxS{} }

func WithValue(parent Ctx, name string, value Val) Ctx {
	c := parent
	if c == nil {
		c = &ctxS{}
	}
	switch cc := c.(type) {
	case interface{ WithValue(name string, value Val) }:
		cc.WithValue(name, value)
	case interface{ Add(name string, value Val) }:
		cc.Add(name, value)
	}
	return c
}

func WithValues(parent Ctx, args ...any) Ctx {
	c := parent
	if c == nil {
		c = &ctxS{}
	}
	switch cc := c.(type) {
	case interface{ WithValues(args ...any) }:
		cc.WithValues(args...)
	case interface{ Add(args ...any) }:
		cc.Add(args...)
	}
	return c
}

type ctxS struct {
	values map[string]Val
	iter   []string
	picked string
}

func (c *ctxS) WithValues(args ...any) {
	var k string
	for _, t := range args {
		if k == "" {
			switch z := t.(type) {
			case string:
				k = z
			case fmt.Stringer:
				k = z.String()
			}
			continue
		}
		c.add(k, t)
		k = ""
	}
}

func (c *ctxS) WithValue(name string, value Val) { c.add(name, value) }

func (c *ctxS) add(name string, value Val) {
	if c.values == nil {
		c.values = make(map[string]Val)
	}
	c.values[name] = value
}

func (c *ctxS) NamesCount() int {
	return len(c.values)
}

func (c *ctxS) ValueBy(name string) Val {
	if v, ok := c.values[name]; ok {
		return v
	}
	return nil
}

func (c *ctxS) NextName() (ret string) {
	if c.iter == nil {
		c.iter = make([]string, 0, len(c.values))
		for k := range c.values {
			c.iter = append(c.iter, k)
		}
		sort.Strings(c.iter)
	}
	if len(c.iter) > 0 {
		ret, c.iter = c.iter[0], c.iter[1:]
	}
	return
}

// Reset resets iterating states so that you can re-iterate
// the values via Next
func (c *ctxS) Reset() {
	c.iter, c.picked = nil, ""
}

func (c *ctxS) Next() bool {
	if c.iter == nil {
		c.iter = make([]string, 0, len(c.values))
		for k := range c.values {
			c.iter = append(c.iter, k)
		}
		sort.Strings(c.iter)
	}
	if len(c.iter) > 0 {
		c.picked, c.iter = c.iter[0], c.iter[1:]
		return true
	}
	return false
}

func (c *ctxS) Entry() (key string, val Val) {
	if v, ok := c.values[c.picked]; ok {
		key, val = c.picked, v
	}
	return
}

func (c *ctxS) Key() (ret string) {
	if _, ok := c.values[c.picked]; ok {
		ret = c.picked
	}
	return
}

func (c *ctxS) Value() Val {
	if v, ok := c.values[c.picked]; ok {
		return v
	}
	return nil
}
