package store

// func ExampleStoreS_Get() {
// 	trie := newBasicStore()
// 	fmt.Println(trie.MustInt("app.dump"))
// 	fmt.Println(trie.MustString("app.dump"))
// 	fmt.Println(trie.MustBool("app.dump")) // convert 3 to bool will get false, only 1 -> true.
// 	// Output:
// 	// 3
// 	// 3
// 	// false
// }
//
// func ExampleStoreS_Dump() {
// 	conf := New()
// 	conf.Set("app.debug", false)
// 	conf.Set("app.verbose", true)
// 	conf.Set("app.dump", 3)
// 	conf.Set("app.logging.file", "/tmp/1.log")
// 	conf.Set("app.server.start", 5)
//
// 	ss := conf.WithPrefix("app.logging")
// 	ss.Set("rotate", 6)
// 	ss.Set("words", []string{"a", "1", "false"})
//
// 	data, found := conf.Get("app.logging.rotate")
// 	println(data, found)
// 	data = conf.MustGet("app.logging.rotate")
// 	println(data)
// 	fmt.Println(conf.MustInt("app.dump"))
// 	fmt.Println(conf.MustString("app.dump"))
// 	fmt.Println(conf.MustBool("app.dump")) // convert 3 to bool will get false, only 1 -> true.
// 	// Output:
// 	// 6 true
// 	// 6
// 	// 3
// 	// 3
// 	// false
// }
