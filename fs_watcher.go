package cruncy

import (
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
)

// FsWatcher monitors a given monPath, checks for file updates every
// pollDur and declares them stable if no updates are detected after stableDur.
type FsWatcher struct {
	monPath            string
	pattern            string
	pollDur, stableDur time.Duration
	files              map[string]time.Time
}

// NewFsWatcher watches a folder for filesystem events.
func NewFsWatcher(monPath, pattern string, pollDur time.Duration, stableDur time.Duration) *FsWatcher {
	if len(pattern) == 0 {
		pattern = "*"
	}

	return &FsWatcher{monPath: monPath, pattern: pattern, pollDur: pollDur, stableDur: stableDur, files: map[string]time.Time{}}
}

// Watch returns stable filenames on the channel.
func (w *FsWatcher) Watch() <-chan string {
	const NumStableFiles = 5000
	stableFilenameCh := make(chan string, NumStableFiles)
	go w.watch(stableFilenameCh)

	return stableFilenameCh
}

func (w *FsWatcher) watch(stableFilenameCh chan string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Errorf("Watcher.watch> Unable to start notify watcher with error %s", err)
		return
	}
	defer watcher.Close()

	done := make(chan bool)
	evch := make(chan fsnotify.Event)

	go w.innerWatch(stableFilenameCh, evch, watcher)

	err = watcher.Add(w.monPath)
	if err != nil {
		logrus.Errorf("Watcher.watch> Unable to sadd to watcher with error %s", err)
		return
	}
	<-done
}

func (w *FsWatcher) innerWatch(stableFilenameCh chan string, eventCh chan fsnotify.Event, watcher *fsnotify.Watcher) {
	monch := w.mon(eventCh)
	for {
		select {
		case event := <-watcher.Events:
			eventCh <- event
		case s := <-monch:
			fileName := filepath.Base(s)
			matched, err := filepath.Match(w.pattern, fileName)
			if err != nil {
				logrus.WithError(err).Infof("Error matching %s with %s gives error: %s", s, w.pattern, err)
				continue
			}
			if !matched {
				logrus.Debugf("%s dit not match %s", fileName, w.pattern)
				continue
			}

			stableFilenameCh <- s
			logrus.Debugf("DELIVERED STABLE: %v\n", filepath.Base(s))
		case err := <-watcher.Errors:
			if err != nil {
				logrus.WithError(err).Printf("Watcher.Watch FsNotify returned error: %v", err)
			}
		}
	}
}

// Watcher:  local un-exported methods
func (w *FsWatcher) timeStampFiles() {
	now := time.Now()
	g, err := filepath.Glob(filepath.Join(w.monPath, w.pattern))
	if err != nil {
		logrus.Fatal("filepath:", err)
	}
	for _, v := range g {
		w.files[v] = now
	}
}
func (w *FsWatcher) mon(ev <-chan fsnotify.Event) <-chan string {
	ch := make(chan string)
	poll := time.Tick(w.pollDur)
	go func(e <-chan fsnotify.Event, s chan<- string) {

		var now time.Time
		for {
			select {
			case event := <-e:
				eventName := filepath.Clean(event.Name)
				switch {
				case event.Op&fsnotify.Create == fsnotify.Create:
					w.files[eventName] = time.Now()
				case event.Op&fsnotify.Write == fsnotify.Write:
					w.files[eventName] = time.Now()
				case event.Op&fsnotify.Chmod == fsnotify.Chmod:
					w.files[eventName] = time.Now()
				case event.Op&fsnotify.Remove == fsnotify.Remove:
					delete(w.files, eventName)
				case event.Op&fsnotify.Rename == fsnotify.Rename:
					delete(w.files, eventName)
				}
			case <-poll:
				now = time.Now()
				for fn, ts := range w.files {
					if now.After(ts.Add(w.stableDur)) {
						delete(w.files, fn)
						s <- fn
					}
				}
			}
		}
	}(ev, ch)
	return ch
}
