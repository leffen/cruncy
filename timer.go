package cruncy


import (
  "time"
  "log"
)

type TimerData struct {
  StartTimeRun   time.Time
  StartTimeBatch time.Time
  EndTimeRun      time.Time

  BatchSize      int64
  PrevRows       int64
  Index          int64
}

func NewTimer() *TimerData {
  timer := &TimerData{}
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
  duration_sec := timer.TotalDuractionSeconds()
  if duration_sec > 0 {
    log.Printf("Total duration:, %v rows =%d row time=%d rows/sec ", duration, timer.Index, timer.Index / duration_sec )
  } else {
    log.Printf("Total duration:, %v rows =%d  SUPER FAST", duration, timer.Index)

  }
}

func (timer *TimerData) ShowBatchTime() {
  diff := timer.Index - timer.PrevRows

  t1 := time.Now()
  var duration time.Duration = t1.Sub(timer.StartTimeBatch)
  var d2 time.Duration = timer.TotalDuration()

  var ds int64 = int64(d2.Seconds())
  if (ds > 0) {
    log.Printf("%d rows avg flow %d/s - batch time %v batch size %d batch_flow %d \n", timer.Index, timer.Index / ds, duration,diff,diff/ds)
  } else {
    log.Printf("%d rows - batch time %v \n", timer.Index,  duration)
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


