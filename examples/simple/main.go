package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/arl/statsviz"

	"github.com/hedzr/is"

	"github.com/hedzr/store"

	"net/http/pprof"
)

func main() {
	println(`
> To view the instant heap stats graph, open the page:

    http://localhost:8080/debug/statsviz

> To view the pprof endpoints, try to run this command in a shell:

    $ go tool pprof -http=:6060 http://localhost:8080/debug/pprof/heap

> Refs:
  - https://hedzr.com/golang/profiling/golang-pprof/ (Chinese)
  - https://stackoverflow.com/questions/24863164/how-to-analyze-golang-memory
  - https://pkg.go.dev/runtime#MemStats
  - https://pkg.go.dev/net/http/pprof

	`)

	svr := prepareHTTPServer()
	conf := prepareStore()

	// go func() { log.Println(http.ListenAndServe("localhost:8080", mux)) }()

	ctx := context.Background()
	catcher := is.Signals().Catch()
	catcher.
		WithOnSignalCaught(func(ctx context.Context, sig os.Signal, wg *sync.WaitGroup) {
			println()
			_, _ = sig, wg
			// logger.Debug("signal caught", "sig", sig)
			if err := svr.Shutdown(context.TODO()); err != nil {
				wg.Done() // just for programmatic safety
				log.Fatal("server shutdown error", "err", err)
			}
			// wg.Done() // http server shutdown ok, so closing catcher's wg counter.
			testStore(conf)
		}).
		WaitFor(ctx, func(ctx context.Context, closer func()) {
			// server.Debug("entering looper's loop...")
			defer closer()
			err := svr.ListenAndServe()
			if err != nil {
				log.Fatal("server serve failed", "err", err)
			}
		})
}

func prepareHTTPServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))

	// see: https://github.com/arl/statsviz
	// another: https://github.com/go-echarts/statsview
	_ = statsviz.Register(mux) // check the debug page at: http://localhost:8080/debug/statsviz

	svr := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return svr
}

func prepareStore() store.Store {
	conf := newBasicStore()
	conf.Set("app.logging.words", []any{"a", 1, false})
	conf.Set("app.server.sites", -1)
	return conf
}

func testStore(conf store.Store) {
	fmt.Printf("\nPath\n%v\n", conf.Dump())
}

func newBasicStore(opts ...store.Opt) store.Store {
	conf := store.New(opts...)
	conf.Set("app.debug", false)
	conf.Set("app.verbose", true)
	conf.Set("app.dump", 3)
	conf.Set("app.logging.file", "/tmp/1.log")
	conf.Set("app.server.start", 5)

	// conf.Set("app.logging.rotate", 6)
	// conf.Set("app.logging.words", []string{"a", "1", "false"})

	ss := conf.WithPrefix("app.logging")
	ss.Set("rotate", 6)
	ss.Set("words", []string{"a", "1", "false"})
	return conf
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
