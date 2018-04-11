package cruncy

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

// AtomicIntCounter uses an int64 internally.
type AtomicIntCounter int64

// Get returns current counter
func (c *AtomicIntCounter) Get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

func (c *AtomicIntCounter) inc() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

// TimerData containes run time data
type TimerData struct {
	Title          string
	Uuid           string
	StartTimeRun   time.Time
	StartTimeBatch time.Time
	EndTimeRun     time.Time

	BatchSize  int64
	PrevRows   int64
	Index      AtomicIntCounter
	ErrorCount AtomicIntCounter
	mu         sync.RWMutex
}

// NewTimer creates a new timer struct
func NewTimer(title string) *TimerData {
	timer := &TimerData{}
	timer.Title = title
	timer.Uuid = ksuid.New().String()
	timer.StartTimeRun = time.Now()
	timer.StartTimeBatch = timer.StartTimeRun
	timer.PrevRows = 0
	timer.ErrorCount = 0
	return timer
}

// BatchDuractionSeconds returns durection in seconds
func (timer *TimerData) BatchDuractionSeconds() int64 {
	t1 := time.Now()
	duration := t1.Sub(timer.StartTimeBatch)
	return int64(duration.Seconds())
}

// TotalDuractionSeconds returns total duration in seconds
func (timer *TimerData) TotalDuractionSeconds() int64 {
	t1 := time.Now()
	duration := t1.Sub(timer.StartTimeRun)
	return int64(duration.Seconds())
}

// TotalDuration returns duration as a time.Duration
func (timer *TimerData) TotalDuration() time.Duration {
	t1 := time.Now()
	return t1.Sub(timer.StartTimeRun)
}

// ShowTotalDuration outputs duration to log with fields
func (timer *TimerData) ShowTotalDuration() {

	cnt := timer.Index.Get()
	timer.mu.RLock()
	uuid := timer.Uuid
	title := timer.Title
	startTime := timer.StartTimeRun
	timer.mu.RUnlock()

	t1 := time.Now()
	duration := t1.Sub(startTime)
	ds := int64(duration.Seconds())
	if ds > 0 {
		msg := fmt.Sprintf("Total duration:, %v rows =%d rate = %d rows/sec ", duration, cnt, cnt/ds)
		log.WithFields(log.Fields{
			"uuid":       uuid,
			"title":      title,
			"total_rows": cnt,
			"avg_flow":   cnt / ds,
			"State":      "stopped",
		}).Info(msg)
	} else {
		log.WithFields(log.Fields{
			"uuid":       uuid,
			"title":      title,
			"total_rows": cnt,
			"avg_flow":   cnt,
			"State":      "stopped",
		}).Infof("Total duration:, %v rows =%d  SUPER FAST", duration, cnt)

	}
}

// ShowBatchTime show averages to now
func (timer *TimerData) ShowBatchTime() {

	cnt := timer.Index.Get()
	timer.mu.RLock()
	uuid := timer.Uuid
	title := timer.Title
	prevRows := timer.PrevRows
	startTime := timer.StartTimeBatch
	timer.mu.RUnlock()

	diff := cnt - prevRows

	t1 := time.Now()
	duration := t1.Sub(startTime)
	d2 := timer.TotalDuration()

	ds := int64(d2.Seconds())
	dsBatch := int64(duration.Seconds())

	if ds > 0 && dsBatch > 0 {
		msg := fmt.Sprintf("%d rows avg flow %d/s - batch time %v batch size %d batch_flow %d \n", cnt, cnt/ds, duration, diff, diff/dsBatch)
		log.WithFields(log.Fields{
			"uuid":       uuid,
			"title":      title,
			"index":      cnt,
			"total_flow": cnt / ds,
			"batch_time": duration,
			"batch_size": diff,
			"batch_flow": diff / dsBatch,
			"State":      "in_batch",
		}).Info(msg)
	} else {
		log.Printf("%d rows - batch time %v \n", cnt, duration)
	}

	timer.mu.Lock()
	timer.PrevRows = cnt
	timer.StartTimeBatch = time.Now()
	timer.mu.Unlock()
}

// Tick increases tick with one
func (timer *TimerData) Tick() {
	cnt := timer.Index.inc()

	if cnt%100000 == 0 {
		timer.ShowBatchTime()
	}
}

// Stop stops the timer
func (timer *TimerData) Stop() time.Time {
	timer.mu.Lock()
	timer.EndTimeRun = time.Now()
	timer.mu.Unlock()
	return timer.EndTimeRun
}

// IncError adds one to number of errors
func (timer *TimerData) IncError() int64 {
	timer.ErrorCount.inc()
	return timer.ErrorCount.Get()
}
