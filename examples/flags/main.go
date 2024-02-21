package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	logz "github.com/hedzr/logg/slog"

	"github.com/hedzr/store"
	"github.com/hedzr/store/providers/flags"
)

func main() {
	storeTest()

	// if err := root(os.Args[1:]); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
}

func storeTest() {
	wordPtr := flag.String("word", "foo", "a string")

	numbPtr := flag.Int("numb", 42, "an int")
	forkPtr := flag.Bool("fork", false, "a bool")

	var svar string
	flag.StringVar(&svar, "svar", "bar", "a string var")

	durPtr := flag.Duration("duration", 5*time.Hour, "a duration")
	timePtr := flag.String("time", "2020-01-01", "a time string")
	typPtr := flag.String("type", "xxx", "type of the app")

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var conf store.Store
	conf = newStore(ctx)
	defer conf.Close()

	err := sample(ctx, conf)
	if err != nil {
		// slog.Error("wrong", "err", err)
		logz.Fatal("Load etcd source failed.", "err", err)
	}

	println(conf.Dump())

	fmt.Println("word:", *wordPtr)
	fmt.Println("numb:", *numbPtr)
	fmt.Println("fork:", *forkPtr)
	fmt.Println("svar:", svar)
	fmt.Println("tail:", flag.Args())
	fmt.Printf("dura: %v\n", *durPtr)
	fmt.Printf("time: %v\n", *timePtr)
	fmt.Printf("type: %v\n", *typPtr)
}

func newStore(ctx context.Context) (conf store.Store) {
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
	return
}

func sample(ctx context.Context, conf store.Store) (err error) {
	_, err = conf.Load(ctx,
		store.WithStorePrefix("app.flags"),
		store.WithPosition("/app"),
		store.WithProvider(flags.New()),

		store.WithStoreFlattenSlice(true), // expand map or slice in value
	)
	if err != nil {
		return
	}
	return
}

// sub-command:

// func root(args []string) error {
// 	if len(args) < 1 {
// 		return errors.New("You must pass a sub-command")
// 	}
//
// 	cmds := []Runner{
// 		NewGreetCommand(),
// 	}
//
// 	subcommand := os.Args[1]
//
// 	for _, cmd := range cmds {
// 		if cmd.Name() == subcommand {
// 			cmd.Init(os.Args[2:])
// 			return cmd.Run()
// 		}
// 	}
//
// 	return fmt.Errorf("Unknown subcommand: %s", subcommand)
// }
//
// func NewGreetCommand() *GreetCommand {
// 	gc := &GreetCommand{
// 		fs: flag.NewFlagSet("greet", flag.ContinueOnError),
// 	}
//
// 	gc.fs.StringVar(&gc.name, "name", "World", "name of the person to be greeted")
//
// 	return gc
// }
//
// type GreetCommand struct {
// 	fs *flag.FlagSet
//
// 	name string
// }
//
// func (g *GreetCommand) Name() string {
// 	return g.fs.Name()
// }
//
// func (g *GreetCommand) Init(args []string) error {
// 	return g.fs.Parse(args)
// }
//
// func (g *GreetCommand) Run() error {
// 	fmt.Println("Hello", g.name, "!")
// 	return nil
// }
//
// type Runner interface {
// 	Init([]string) error
// 	Run() error
// 	Name() string
// }
