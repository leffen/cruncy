package cruncy

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/segmentio/ksuid"
	log "github.com/sirupsen/logrus"
)

// TimerData containes run time data
type TimerData struct {
	Title          string
	Uuid           string
	StartTimeRun   time.Time
	StartTimeBatch time.Time
	EndTimeRun     time.Time

	BatchSize  int64
	PrevRows   int64
	Index      int64
	ErrorCount int64
	mu         sync.RWMutex
	muShow     sync.RWMutex
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
	duration := timer.TotalDuration()
	ds := timer.TotalDuractionSeconds()
	if ds > 0 {
		msg := fmt.Sprintf("Total duration:, %v rows =%d row time=%d rows/sec ", duration, timer.Index, timer.Index/ds)
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"total_rows": timer.Index,
			"avg_flow":   timer.Index / ds,
			"State":      "stopped",
		}).Info(msg)
	} else {
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"total_rows": timer.Index,
			"avg_flow":   timer.Index,
			"State":      "stopped",
		}).Infof("Total duration:, %v rows =%d  SUPER FAST", duration, timer.Index)

	}
}

// ShowBatchTime show averages to now
func (timer *TimerData) ShowBatchTime() {
	timer.muShow.RLock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer timer.muShow.RUnlock()

	diff := timer.Index - timer.PrevRows

	t1 := time.Now()
	duration := t1.Sub(timer.StartTimeBatch)
	d2 := timer.TotalDuration()

	ds := int64(d2.Seconds())
	dsBatch := int64(duration.Seconds())

	if ds > 0 && dsBatch > 0 {
		msg := fmt.Sprintf("%d rows avg flow %d/s - batch time %v batch size %d batch_flow %d \n", timer.Index, timer.Index/ds, duration, diff, diff/dsBatch)
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"index":      timer.Index,
			"total_flow": timer.Index / ds,
			"batch_time": duration,
			"batch_size": diff,
			"batch_flow": diff / dsBatch,
			"State":      "in_batch",
		}).Info(msg)
	} else {
		log.Printf("%d rows - batch time %v \n", timer.Index, duration)
	}
	timer.PrevRows = timer.Index
	timer.StartTimeBatch = time.Now()

}

// Tick increases tick with one
func (timer *TimerData) Tick() {
	timer.mu.Lock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer timer.mu.Unlock()

	atomic.AddInt64(&timer.Index, 1)

	if timer.Index%100000 == 0 {
		timer.ShowBatchTime()
	}
}

// Stop stops the timer
func (timer *TimerData) Stop() time.Time {
	timer.EndTimeRun = time.Now()
	return timer.EndTimeRun
}

// IncError adds one to number of errors
func (timer *TimerData) IncError() int64 {
	timer.mu.Lock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer timer.mu.Unlock()

	timer.ErrorCount++
	return timer.ErrorCount
}
