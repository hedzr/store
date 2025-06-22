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

func ExampleStore_BR() {
	conf := New()

	conf.Set("path.to.someone.son.grand-son.mid-name", "Von")
	conf.Set("path.to.someone.son.grand-son-2.mid-name", "Von")
	conf.Set("path.mid-name", "Von")

	son := conf.WithPrefix("path.to.someone.son")
	sbr := son.BR()

	// In BR mode enabled, MustString will match these keys
	// if current node cannot hit 'mid-name' subkey:
	//    - path.to.someone.son.mid-name ? NOT
	//    - path.to.someone.mid-name.    ? NOT
	//    - path.to.mid-name             ? HIT!
	println("sbr[mid-name] =", sbr.MustString("mid-name"))
	// In default mode (RecursiveNone), only current node
	// joint into the matching turn.
	// So empty result returned.
	println("son[mid-name] =", sbr.MustString("mid-name"))

	// Outputs:
	// sbr[mid-name] = Von
	// son[mid-name] =
}
