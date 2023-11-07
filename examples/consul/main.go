package main

import (
	"github.com/hashicorp/consul/api"

	logz "github.com/hedzr/logg/slog"
	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/consul"
)

func main() {
	conf := store.New()
	defer conf.Close()

	var err error
	err = conf.Load(
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
		logz.Fatal("Load consul source failed.", "err", err)
	}

	println(conf.Dump())

	done := make(chan struct{})
	<-done
}
