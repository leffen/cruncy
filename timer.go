package cruncy

import (
  "time"
  log "github.com/Sirupsen/logrus"
  "fmt"
  "github.com/satori/go.uuid"
)

type TimerData struct {
  Title          string
  Uuid           string
  StartTimeRun   time.Time
  StartTimeBatch time.Time
  EndTimeRun     time.Time

  BatchSize      int64
  PrevRows       int64
  Index          int64
}

func NewTimer(title string) *TimerData {
  timer := &TimerData{}
  timer.Title = title
  timer.Uuid = uuid.NewV4().String()
  timer.StartTimeRun = time.Now()
  timer.StartTimeBatch = timer.StartTimeRun
  timer.PrevRows = 0
  return timer
}

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
    msg := fmt.Sprintf("Total duration:, %v rows =%d row time=%d rows/sec ", duration, timer.Index, timer.Index / duration_sec)
    log.WithFields(log.Fields{
      "uuid": timer.Uuid,
      "title": timer.Title,
      "index": timer.Index,
      "total_flow":   timer.Index / ds,
      "State":   "stopped",
    }).Info(msg)
  } else {
    log.Printf("Total duration:, %v rows =%d  SUPER FAST", duration, timer.Index)

  }
}

func (timer *TimerData) ShowBatchTime() {
  diff := timer.Index - timer.PrevRows

  t1 := time.Now()
  var duration time.Duration = t1.Sub(timer.StartTimeBatch)
  var d2 time.Duration = timer.TotalDuration()

  ds := int64(d2.Seconds())
  ds_batch := int64(duration.Seconds())

  if (ds > 0) {
    msg := fmt.Sprintf("%d rows avg flow %d/s - batch time %v batch size %d batch_flow %d \n", timer.Index, timer.Index / ds, duration, diff, diff / ds_batch)
    log.WithFields(log.Fields{
      "uuid":         timer.Uuid,
      "title":        timer.Title,
      "index":        timer.Index,
      "total_flow":   timer.Index / ds,
      "batch_time":   duration,
      "batch_size":   diff,
      "batch_flow":   diff / ds_batch,
      "State":        "in_batch",
    }).Info(msg)
  } else {
    log.Printf("%d rows - batch time %v \n", timer.Index, duration)
  }
  timer.PrevRows = timer.Index
  timer.StartTimeBatch = time.Now()

}

func (timer *TimerData) Tick() {
  timer.Index++

  if timer.Index % 100000 == 0 {
    timer.ShowBatchTime()
  }
}

func (timer *TimerData) Stop() time.Time {
  timer.EndTimeRun = time.Now()
  return timer.EndTimeRun
}


