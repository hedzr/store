package radix

import (
	"fmt"
	"runtime"
	"strings"
)

// Source describes the location of a line of source code.
type Source struct {
	// Function is the package path-qualified function name containing the
	// source line. If non-empty, this string uniquely identifies a single
	// function in the program. This may be the empty string if not known.
	Function string `json:"function"`
	// File and Line are the file name and line number (1-based) of the source
	// line. These may be the empty string and zero, respectively, if not known.
	File string `json:"file"`
	Line int    `json:"line"`
}

// func (s Source) toGroup() (as Attr) {
// 	as = &gkvp{"source", Attrs{
// 		&kvp{"function", s.Function},
// 		&kvp{"file", s.File},
// 		&kvp{"line", s.Line},
// 	}}
// 	// as = Group("source",
// 	// 	NewAttr("function", s.Function),
// 	// 	NewAttr("file", s.File),
// 	// 	NewAttr("line", s.Line))
// 	return
// }

func getpc(skip int, extra ...int) (pc uintptr) {
	var pcs [1]uintptr
	for _, ee := range extra {
		if ee > 0 {
			skip += ee
		}
	}
	runtime.Callers(skip+1, pcs[:])
	pc = pcs[0]
	return
}

func getpcsource(pc uintptr) Source {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	return Source{
		Function: frame.Function,
		File:     checkpath(frame.File),
		Line:     frame.Line,
	}
}

func checkpath(file string) string { return file }

func stack(skip, nFrames int) string {
	pcs := make([]uintptr, nFrames+1)
	n := runtime.Callers(skip+1, pcs)
	if n == 0 {
		return "(no stack)"
	}
	frames := runtime.CallersFrames(pcs[:n])
	var b strings.Builder
	i := 0
	for {
		frame, more := frames.Next()
		fmt.Fprintf(&b, "called from %s (%s:%d)\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
		i++
		if i >= nFrames {
			fmt.Fprintf(&b, "(rest of stack elided)\n")
			break
		}
	}
	return b.String()
}

// func checkpath(file string) string {
// 	// if s.curdir == "" {
// 	// 	s.curdir, _ = os.Getwd()
// 	// }
// 	// if strings.HasPrefix(file, s.curdir) {
// 	// 	file= file[len(s.curdir)+1:]
// 	// }
//
// 	if IsAnyBitsSet(Lprivacypath) {
// 		for k, v := range knownPathMap {
// 			if strings.HasPrefix(file, k) {
// 				file = strings.ReplaceAll(file, k, v)
// 			}
// 		}
//
// 		if IsAnyBitsSet(Lprivacypathregexp) {
// 			for _, rpl := range knownPathRegexpMap {
// 				if rpl.expr.MatchString(file) {
// 					file = rpl.expr.ReplaceAllString(file, rpl.repl)
// 				}
// 			}
// 		} else {
// 			if strings.HasPrefix(file, "/Volumes/") {
// 				if pos := strings.IndexRune(file[9:], '/'); pos >= 0 {
// 					file = "~" + file[9+pos:]
// 				}
// 			}
// 		}
// 	}
// 	return file
// }
//
// func checkedfuncname(name string) string {
// 	// name := s.Function
// 	if IsAnyBitsSet(Lcallerpackagename) {
// 		for k, v := range codeHostingProvidersMap {
// 			name = strings.ReplaceAll(name, k, v) // replace github.com with "GH", ...
// 		}
// 	} else {
// 		if pos := strings.LastIndex(name, "/"); pos >= 0 {
// 			name = name[pos+1:] // strip the leading package names, eg: "GH/hedzr/logg/slog/" will be removed
// 		}
// 	}
// 	return name
// }
