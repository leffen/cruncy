package cruncy

import (
	"context"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	d, err := ioutil.TempDir("", "test")
	assert.Nil(t, err)
	defer os.RemoveAll(d)

	CreateDirUnlessExists(d)
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func(ctx context.Context, cancel context.CancelFunc, d string) {
		numFiles := 0
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Nanosecond * 100):
				f, err := ioutil.TempFile(d, "test")
				assert.Nil(t, err)
				time.Sleep(time.Millisecond * 100)
				os.Remove(f.Name())
				numFiles++
				if numFiles > 3 {
					cancel()
					return
				}
			}

		}
	}(ctx, cancel, d)

	w := NewFsWatcher(d, "*", time.Nanosecond*200, time.Millisecond*10)

	rFiles := 0
	watch := w.Watch()
	for {
		select {
		case <-ctx.Done():
			return
		case fileName := <-watch:
			logrus.Debugf("filename: %s", fileName)
			rFiles++
		case <-time.After(time.Second * 2):
			if rFiles < 3 {
				t.Fatal("Missed files")
			}
		}

	}

}
