package main

import (
	"context"
	"flag"

	"github.com/hashicorp/consul/api"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/consul"
)

var complex = flag.Bool("complex", false, "another complex testing")

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var conf store.Store
	conf = store.New(
		store.WithOnNewHandlers(func(path string, value any, mergingMapOrLoading bool) {
			logz.InfoContext(ctx, "created", "path", path, "value", value, "mergingMapOrLoading", mergingMapOrLoading)
			println(conf.Dump())
		}),
		store.WithOnChangeHandlers(func(path string, value, oldValue any, mergingMapOrLoading bool) {
			logz.InfoContext(ctx, "modified", "path", path, "value", value, "oldVal", oldValue, "mergingMapOrLoading", mergingMapOrLoading)
			println(conf.Dump())
		}),
		store.WithOnDeleteHandlers(func(path string, value any, mergingMapOrLoading bool) {
			logz.InfoContext(ctx, "removed", "path", path, "value", value, "mergingMapOrLoading", mergingMapOrLoading)
			println(conf.Dump())
		}),
	)
	defer conf.Close()

	var err error
	if *complex {
		err = complexTest(ctx, conf) // not yet, reserved for the future
	} else {
		err = simpleTest(ctx, conf)
	}
	if err != nil {
		// slog.Error("wrong", "err", err)
		logz.Fatal("Load consul source failed.", "err", err)
	}
}

func simpleTest(ctx context.Context, conf store.Store) (err error) {
	_, err = conf.Load(ctx,
		store.WithStorePrefix("app.conf"),
		store.WithPosition("/testconsul"),
		store.WithProvider(consul.New(
			consul.WithWatchEnabled(true),
			consul.WithConsulConfig(api.DefaultConfig()),
			consul.WithRecursive(true),
			consul.WithProcessMeta(false),
			consul.WithStripPrefix("testconsul/conf"),
			consul.WithPrependPrefix("consul"),
			consul.WithDelimiter(string([]rune{conf.Delimiter()})),
		)),
	)
	if err != nil {
		return
	}

	// println(conf.Dump())

	done := make(chan struct{})
	<-done
	return
}

func complexTest(ctx context.Context, conf store.Store) (err error) {
	return
}

func init() {
	logz.AddFlags(logz.Lprivacypathregexp | logz.Lprivacypath)
	logz.SetLevel(logz.InfoLevel)
}
