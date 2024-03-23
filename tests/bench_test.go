package tests

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/env"
)

func BenchmarkTrieSingleGetForProfiling(b *testing.B) { //nolint:revive
	conf := newStoreGo()
	// b.Logf("conf tree:\n%v", conf.Dump())
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = conf.Get(getWord(-1))
		}
	})
}

func BenchmarkTrieGet(b *testing.B) { //nolint:revive
	b.Logf("Logging at a disabled level without any structured context.")
	elapsedTimes := make(map[string]time.Duration)

	b.Run("hedzr/storeT[any]", func(b *testing.B) {
		trie := newStoreT()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				trie.Get(getWord(-1))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
	b.Run("hedzr/store", func(b *testing.B) {
		conf := newStoreGo()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, _ = conf.Get(getWord(-1))
			}
		})
		elapsedTimes[b.Name()] = b.Elapsed()
	})
}

func TestStoreDump(t *testing.T) {
	conf := newStore()
	t.Log("\n", conf.Dump())
}

func newStoreGo() store.Store {
	conf := store.New()

	// conf.Set("app.debug", false)
	// // t.Logf("\nPath\n%v\n", conf.dump())
	// conf.Set("app.verbose", true)
	// // t.Logf("\nPath\n%v\n", conf.dump())
	// conf.Set("app.dump", 3) //nolint:revive
	// conf.Set("app.logging.file", "/tmp/1.log")
	// conf.Set("app.server.start", 5)
	//
	// // conf.Set("app.logging.rotate", 6)
	// // conf.Set("app.logging.words", []string{"a", "1", "false"})
	//
	// ss := conf.WithPrefix("app.logging")
	// ss.Set("rotate", 6)
	// ss.Set("words", []string{"a", "1", "false"})

	conf.Load(context.Background(),
		store.WithStorePrefix("env"),
		store.WithProvider(env.New(
			env.WithPrefix("go"),
			env.WithLowerCase(true),
			env.WithUnderlineToDot(true),
		)),

		store.WithStoreFlattenSlice(true), // decode and flatten slice into tree structure instead treat it as a simple value
	)

	// collect and update all the validate keys
	if len(words) < 16 {
		// conf.Walk("", func(path, fragment string, node radix.Node[any]) {
		// 	if strings.HasPrefix(path, "env.") && node.IsLeaf() {
		// 		words = append(words, path)
		// 		// if strings.HasPrefix(path, "env.go") {
		// 		// 	println(path, ",", fragment, ",")
		// 		// }
		// 	}
		// 	_ = node
		// })

		words = []string{
			"env.gomodcache",
			"env.gopath",
			"env.goproxy",
		}
	}
	return conf
}

func newStoreT() store.MinimalStoreT[any] {
	trie := store.NewStoreT[any]()
	trie.Set("app.debug", 1)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Set("app.verbose", 2)
	// t.Logf("\nPath\n%v\n", trie.dump())
	trie.Set("app.dump", 3)
	trie.Set("app.logging.file", 4)
	trie.Set("app.server.start", 5)
	trie.Set("app.logging.rotate", 6) //
	trie.Set("app.logging.words", []string{"a", "1", "false"})
	_, _ = trie.Get("app.logging.rotate")
	return trie
}

func newStore() store.Store {
	conf := store.New()
	conf.Set("app.debug", false)
	// t.Logf("\nPath\n%v\n", conf.dump())
	conf.Set("app.verbose", true)
	// t.Logf("\nPath\n%v\n", conf.dump())
	conf.Set("app.dump", 3) //nolint:revive
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	// conf.Set("app.logging.rotate", 6)
	// conf.Set("app.logging.words", []string{"a", "1", "false"})

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	_, _ = ss.Get("rotate")
	return conf
}

func intn(max int) int { //nolint:revive,unused
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	n := nBig.Int64()
	return int(n)
}

func getWord(index int) string {
	if index < 0 || index >= len(words) {
		// index = intn(len(words)) //nolint:revive
		index = iterIndex % len(words) //nolint:revive
		iterIndex++
	}
	return words[index]
}

var iterIndex int

var words = []string{
	"a",
	"app",
	"apk",
	"app.de",
	"app.debug",
	"app.dump",
	"app.verbose",
	"app.server.",
	"app.server.start",
	"app.logging.file",
	"app.logging",
	"app.logging.",
	"app.logging.rotate",
	"app.missed",
}
