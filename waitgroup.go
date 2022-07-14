package cruncy

// DEBUG waitgroup implementation with information strings

import (
	"sync/atomic"
	"time"

	"github.com/sirupsen/logrus"
)

type WaitGroup struct {
	count int64
}

func (r *WaitGroup) Add(delta int, info string) {
	logrus.Infof("WaitGroup>ADD  %s  %d Count before: %d", info, delta, r.count)
	atomic.AddInt64(&r.count, int64(delta))
}

func (r *WaitGroup) Done(info string) {
	logrus.Infof("WaitGroup>Done Count before: %d", r.count)
	atomic.AddInt64(&r.count, -1)
}

func (r *WaitGroup) Wait() {
	logrus.Infof("WaitGroup>Wait Count : %d", r.count)
	for atomic.LoadInt64(&r.count) > 0 {
		//		logrus.Infof("WaitGroup>Wait ... Count : %d", r.count)
		time.Sleep(1)
	}
}

func (r *WaitGroup) TryWait() bool {
	return atomic.LoadInt64(&r.count) == 0
}
