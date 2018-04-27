package cruncy

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

// DataFields is map to interface def

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
	Logger     *log.Entry
}

// NewTimer creates a new timer struct
func NewTimer(title string) *TimerData {
	nw := time.Now()
	uuid := ksuid.New().String()
	return &TimerData{
		Title:          title,
		Uuid:           uuid,
		StartTimeRun:   nw,
		StartTimeBatch: nw,
		Logger:         log.WithFields(log.Fields{"uuid": uuid, "title": title}),
		BatchSize:      100000,
	}
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
	startTime := timer.StartTimeRun
	timer.mu.RUnlock()

	t1 := time.Now()
	duration := t1.Sub(startTime)
	ds := int64(duration.Seconds())
	if ds > 0 {
		msg := fmt.Sprintf("Total duration:, %v rows =%d rate = %d rows/sec ", duration, cnt, cnt/ds)
		timer.Logger.WithFields(log.Fields{
			"total_rows": cnt,
			"avg_flow":   cnt / ds,
			"State":      "stopped",
		}).Info(msg)
	} else {
		timer.Logger.WithFields(log.Fields{
			"total_rows": cnt,
			"avg_flow":   cnt,
			"State":      "stopped",
		}).Infof("Total duration:, %v rows =%d  SUPER FAST", duration, cnt)

	}
}

// LogFields returns summary data as fieldmap
func (timer *TimerData) LogFields() log.Fields {
	cnt := timer.Index.Get()
	timer.mu.RLock()
	uuid := timer.Uuid
	title := timer.Title
	startTime := timer.StartTimeRun
	timer.mu.RUnlock()

	t1 := time.Now()
	duration := t1.Sub(startTime)
	ds := int64(duration.Seconds())
	avgFlow := cnt
	if ds > 0 {
		avgFlow = cnt / ds
	}

	return log.Fields{
		"uuid":       uuid,
		"title":      title,
		"total_rows": cnt,
		"duration_s": ds,
		"avg_flow_s": avgFlow,
		"start_time": startTime.UTC().Format("2006-01-02T15:04:05-0700"),
		"end_time":   t1.UTC().Format("2006-01-02T15:04:05-0700"),
	}

}

// ShowBatchTime show averages to now
func (timer *TimerData) ShowBatchTime() {

	cnt := timer.Index.Get()
	timer.mu.RLock()
	prevRows := timer.PrevRows
	startTime := timer.StartTimeBatch
	batchSize := timer.BatchSize
	timer.mu.RUnlock()

	diff := cnt - prevRows

	t1 := time.Now()
	duration := t1.Sub(startTime)
	d2 := timer.TotalDuration()

	ds := int64(d2.Seconds())
	dsBatch := int64(duration.Seconds())
	unit := "s"
	flow := float64(diff) / float64(dsBatch)
	if batchSize < 1000 && flow < 1 && dsBatch > 0 {
		unit = "m"
		flow = float64(diff) / (float64(dsBatch) / 60.0)
	}

	if ds > 0 && dsBatch > 0 {
		msg := fmt.Sprintf("%d rows avg flow %d/%s - batch time %v batch size %d batch_flow %.4f \n", cnt, cnt/ds, unit, duration, diff, flow)
		timer.Logger.WithFields(log.Fields{
			"index":      cnt,
			"total_flow": cnt / ds,
			"batch_time": duration,
			"batch_size": diff,
			"batch_flow": diff / dsBatch,
			"State":      "in_batch",
		}).Info(msg)
	} else {
		timer.Logger.Printf("%d rows - batch time %v \n", cnt, duration)
	}

	timer.mu.Lock()
	timer.PrevRows = cnt
	timer.StartTimeBatch = time.Now()
	timer.mu.Unlock()
}

// Tick increases tick with one
func (timer *TimerData) Tick() {
	cnt := timer.Index.inc()

	if cnt%timer.BatchSize == 0 {
		timer.ShowBatchTime()
	}
}

// Start the timer
func (timer *TimerData) Start() {
	nw := time.Now()
	timer.mu.Lock()
	timer.StartTimeRun = nw
	timer.StartTimeBatch = nw
	timer.mu.Unlock()
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
