package store

type Opt func(s *storeS) // Opt(ions) for New Store

// WithDelimiter sets the delimiter char.
//
// A delimiter char is generally used for extracting the key-value
// pair via GetXXX, MustXXX, e.g., MustInt, MustStringSlice, ....
func WithDelimiter(delimiter rune) Opt {
	return func(s *storeS) {
		s.SetDelimiter(delimiter)
	}
}

// WithPrefix sets the associated prefix for the tree path.
func WithPrefix(prefix ...string) Opt {
	return func(s *storeS) {
		s.SetPrefix(prefix...)
	}
}

// WithOnChangeHandlers allows user's handlers can be callback once a node changed.
func WithOnChangeHandlers(handlers ...OnChangeHandler) Opt {
	return func(s *storeS) {
		s.onChangeHandlers = append(s.onChangeHandlers, handlers...)
	}
}

// WithOnNewHandlers allows user's handlers can be callback if a new node has been creating.
func WithOnNewHandlers(handlers ...OnNewHandler) Opt {
	return func(s *storeS) {
		s.onNewHandlers = append(s.onNewHandlers, handlers...)
	}
}

// WithOnDeleteHandlers allows user's handlers can be callback once a node removed.
func WithOnDeleteHandlers(handlers ...OnDeleteHandler) Opt {
	return func(s *storeS) {
		s.OnDeleteHandlers = append(s.OnDeleteHandlers, handlers...)
	}
}

// WithFlattenSlice sets a bool flag to tell Store the slice value should be
// treated as node leaf. The index of the slice would be part of node path.
// For example, you're loading a slice []string{"A","B"} into node path
// "app.slice", the WithFlattenSlice(true) causes the following structure:
//
//	app.slice.0 => "A"
//	app.slice.1 => "B"
//
// Also, WithFlattenSlice makes the map values to be flattened into a tree.
func WithFlattenSlice(b bool) Opt {
	return func(s *storeS) {
		s.flattenSlice = b
	}
}

// WithWatchEnable allows watching the external source if its provider
// supports Watchable ability.
func WithWatchEnable(b bool) Opt {
	return func(s *storeS) {
		s.allowWatch = b
	}
}
