package file

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/hedzr/store"
)

func New(file string, opts ...Opt) *pvdr {
	s := &pvdr{file: file}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func WithCodec(codec store.Codec) Opt {
	return func(s *pvdr) {
		s.codec = codec
	}
}

func WithPosition(prefix string) Opt {
	return func(s *pvdr) {
		s.prefix = prefix
	}
}

func WithWatchEnabled(b bool) Opt {
	return func(s *pvdr) {
		s.watchEnabled = b
	}
}

type Opt func(s *pvdr)

type pvdr struct {
	file         string
	watchEnabled bool
	watching     int32
	codec        store.Codec
	prefix       string
}

func (s *pvdr) Count() int {
	return 0
}

func (s *pvdr) Has(key string) bool {
	return false
}

func (s *pvdr) Next() (key string, eol bool) {
	key, eol = "", true
	return
}

func (s *pvdr) Keys() (keys []string, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) Value(key string) (value any, ok bool) {
	value, ok = nil, false
	return
}

func (s *pvdr) MustValue(key string) (value any) {
	value = nil
	return
}

func (s *pvdr) Reader() (r *store.Reader, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) Read() (data map[string]any, err error) {
	err = store.NotImplemented
	return
}

func (s *pvdr) ReadBytes() (data []byte, err error) {
	data, err = os.ReadFile(s.file)
	return
}

func (s *pvdr) Write(data []byte) (err error) {
	err = os.WriteFile(s.file, data, 0644)
	return
}

func (s *pvdr) GetCodec() (codec store.Codec) { return s.codec }
func (s *pvdr) GetPosition() (pos string)     { return s.prefix }
func (s *pvdr) WithCodec(codec store.Codec)   { s.codec = codec }
func (s *pvdr) WithPosition(prefix string)    { s.prefix = prefix }

func (s *pvdr) Close() {
	atomic.CompareAndSwapInt32(&s.watching, 1, 0)
}

func (s *pvdr) Watch(ctx context.Context, cb func(event any, err error)) (err error) {
	if s.watchEnabled == false {
		return nil
	}

	// Resolve symlinks and save the original path so that changes to symlinks
	// can be detected.
	var realPath string
	realPath, err = filepath.EvalSymlinks(s.file)
	if err != nil {
		return
	}
	realPath = filepath.Clean(realPath)

	// Although only a single file is being watched, fsnotify has to watch
	// the whole parent directory to pick up all events such as symlink changes.
	fDir, _ := filepath.Split(s.file)

	var w *fsnotify.Watcher
	w, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go s.watchRunner(ctx, cb, w, realPath)

	// Watch the directory for changes.
	return w.Add(fDir)
}

func (s *pvdr) watchRunner(ctx context.Context, cb func(event any, err error), w *fsnotify.Watcher, realPath string) {
	var (
		// lastEvent     string
		// lastEventTime time.Time
		lastChange = changeS{provider: s}
	)

	if !atomic.CompareAndSwapInt32(&s.watching, 0, 1) {
		return
	}

	var ok bool
	var err error
	var event fsnotify.Event
loop:
	for atomic.LoadInt32(&s.watching) == 1 {
		select {
		case <-ctx.Done():
			return

		case event, ok = <-w.Events:
			if !ok {
				cb(nil, errors.New("fsnotify watch channel closed"))
				break loop
			}

			// Use a simple timer to buffer events as certain events fire
			// multiple times on some platforms.
			if event.String() == lastChange.lastEvent && time.Since(lastChange.lastEventTime) < time.Millisecond*5 {
				continue
			}
			lastChange.lastEvent = event.String()
			lastChange.lastEventTime = time.Now()

			evFile := filepath.Clean(event.Name)

			// Since the event is triggered on a directory, is this
			// one on the file being watched?
			if evFile != realPath && evFile != s.file {
				continue
			}

			// The file was removed.
			if event.Op&fsnotify.Remove != 0 {
				cb(nil, fmt.Errorf("file %s was removed", event.Name))
				break loop
			}

			// Resolve symlink to get the real path, in case the symlink's
			// target has changed.
			var curPath string
			curPath, err = filepath.EvalSymlinks(s.file)
			if err != nil {
				cb(nil, err)
				break loop
			}
			lastChange.realPath = filepath.Clean(curPath)

			// Finally, we only care about create and write.
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			// Trigger event.
			cb(lastChange, nil)

		// There's an error.
		case err, ok = <-w.Errors:
			if !ok {
				cb(nil, errors.New("fsnotify err channel closed"))
				break loop
			}

			// Pass the error to the callback.
			cb(nil, err)
			break loop
		}
	}

	w.Close()
}
