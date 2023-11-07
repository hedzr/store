package main

import (
	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/etcd"
)

func main() {
	conf := store.New()
	defer conf.Close()

	var err error
	err = conf.Load(
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
		logz.Fatal("Load consul source failed.", "err", err)
	}

	println(conf.Dump())

	done := make(chan struct{})
	<-done
}
