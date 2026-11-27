package main

import (
	"context"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/etcd"
)

func main() {
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

	err = simpleTest(ctx, conf)

	if err != nil {
		// slog.Error("wrong", "err", err)
		logz.Fatal("Load etcd source failed.", "err", err)
	}
}

func simpleTest(ctx context.Context, conf store.Store) (err error) {
	_, err = conf.Load(ctx,
		store.WithStorePrefix("app.conf"),
		store.WithPosition("/app"),
		store.WithProvider(etcd.New(
			etcd.WithEndpoints("127.0.0.1:2379"),
			etcd.WithWatchEnabled(true),
			etcd.WithRecursive(true),
			etcd.WithProcessMeta(false),
			etcd.WithStripPrefix("/app"),
			etcd.WithPrependPrefix("etcd"),
			etcd.WithDelimiter(string([]rune{conf.Delimiter()})),
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

func init() {
	logz.AddFlags(logz.Lprivacypathregexp | logz.Lprivacypath)
	logz.SetLevel(logz.InfoLevel)
}
