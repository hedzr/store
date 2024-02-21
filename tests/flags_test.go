package tests

import (
	"context"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/flags"
)

func TestStore_flags_Load(t *testing.T) {
	wordPtr := flag.String("word", "foo", "a string")

	numbPtr := flag.Int("numb", 42, "an int")
	forkPtr := flag.Bool("fork", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	durPtr := flag.Duration("duration", 5*time.Hour, "a duration")
	timePtr := flag.String("time", "2020-01-01", "a time string")
	typPtr := flag.String("type", "xxx", "type of the app")

	s := newBasicStore()
	if _, err := s.Load(context.TODO(),
		store.WithStorePrefix("app.flags"),
		store.WithProvider(flags.New()),

		store.WithStoreFlattenSlice(true), // expand map or slice in value
	); err != nil {
		t.Fatalf("failed: %v", err)
	}

	ret := s.Dump()
	t.Logf("\nPath\n%v\n", ret)

	assertEqual(t, *timePtr, s.MustString("app.flags.time"))
	assertEqual(t, svar, s.MustString("app.flags.svar"))
	assert.Equal(t, *typPtr, s.MustString("app.flags.type"))
	assert.Equal(t, *durPtr, s.MustDuration("app.flags.duration"))
	assert.Equal(t, false, s.MustBool("app.flags.fork"))
	assert.Equal(t, 42, s.MustInt("app.flags.numb"))
	assert.Equal(t, "foo", s.MustString("app.flags.word"))

	t.Logf("word: %v", *wordPtr)
	t.Logf("numb: %v", *numbPtr)
	t.Logf("fork: %v", *forkPtr)
	t.Logf("svar: %v", svar)
	t.Logf("tail: %v", flag.Args())
}

func newMap() map[string]any {
	return map[string]any{
		"cool.station": map[any]any{
			8 * time.Second: true,
			"flush": map[string]any{
				"always":   false,
				"interval": 5 * time.Hour,
			},
		},
		"desc": "a desc string",
	}
}
