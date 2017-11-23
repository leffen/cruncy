package cruncy

import (
	"fmt"
	"sync"
	"time"

	"github.com/renstrom/shortuuid"
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
	timer.Uuid = shortuuid.New()
	timer.StartTimeRun = time.Now()
	timer.StartTimeBatch = timer.StartTimeRun
	timer.PrevRows = 0
	timer.ErrorCount = 0
	return timer
}

// BatchDuractionSeconds returns durection in seconds
func (timer TimerData) BatchDuractionSeconds() int64 {
	t1 := time.Now()
	var duration time.Duration = t1.Sub(timer.StartTimeBatch)
	return int64(duration.Seconds())
}

func (timer TimerData) TotalDuractionSeconds() int64 {
	t1 := time.Now()
	var duration time.Duration = t1.Sub(timer.StartTimeRun)
	return int64(duration.Seconds())
}

func (timer TimerData) TotalDuration() time.Duration {
	t1 := time.Now()
	return t1.Sub(timer.StartTimeRun)
}

func (timer TimerData) ShowTotalDuration() {
	duration := timer.TotalDuration()
	ds := timer.TotalDuractionSeconds()
	if ds > 0 {
		msg := fmt.Sprintf("Total duration:, %v rows =%d row time=%d rows/sec ", duration, timer.Index, timer.Index/ds)
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"index":      timer.Index,
			"total_flow": timer.Index / ds,
			"State":      "stopped",
		}).Info(msg)
	} else {
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"index":      timer.Index,
			"total_flow": timer.Index,
			"State":      "stopped",
		}).Infof("Total duration:, %v rows =%d  SUPER FAST", duration, timer.Index)

	}
}

func (timer *TimerData) ShowBatchTime() {
	timer.muShow.RLock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer timer.muShow.RUnlock()

	diff := timer.Index - timer.PrevRows

	t1 := time.Now()
	var duration time.Duration = t1.Sub(timer.StartTimeBatch)
	var d2 time.Duration = timer.TotalDuration()

	ds := int64(d2.Seconds())
	ds_batch := int64(duration.Seconds())

	if ds > 0 && ds_batch > 0 {
		msg := fmt.Sprintf("%d rows avg flow %d/s - batch time %v batch size %d batch_flow %d \n", timer.Index, timer.Index/ds, duration, diff, diff/ds_batch)
		log.WithFields(log.Fields{
			"uuid":       timer.Uuid,
			"title":      timer.Title,
			"index":      timer.Index,
			"total_flow": timer.Index / ds,
			"batch_time": duration,
			"batch_size": diff,
			"batch_flow": diff / ds_batch,
			"State":      "in_batch",
		}).Info(msg)
	} else {
		log.Printf("%d rows - batch time %v \n", timer.Index, duration)
	}
	timer.PrevRows = timer.Index
	timer.StartTimeBatch = time.Now()

}

func (timer *TimerData) Tick() {
	timer.mu.RLock() // Claim the mutex as a RLock - allowing multiple go routines to log simultaneously
	defer timer.mu.RUnlock()

	timer.Index++

	if timer.Index%100000 == 0 {
		timer.ShowBatchTime()
	}
}

func (timer *TimerData) Stop() time.Time {
	timer.EndTimeRun = time.Now()
	return timer.EndTimeRun
}

func (timer *TimerData) IncError() int64 {
	timer.ErrorCount++
	return timer.ErrorCount
}
