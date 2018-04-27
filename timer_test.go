package cruncy

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFields(t *testing.T) {
	tmr := NewTimer("Test")
	tmr.Tick()
	fields := tmr.LogFields()
	assert.NotNil(t, fields)
	assert.Equal(t, "Test", fields["title"])
	assert.Equal(t, int64(1), fields["total_rows"])
}

func TestTimer(t *testing.T) {
	tmr := NewTimer("test2")
	tmr.Start()
	tmr.Tick()
	<-time.After(100 * time.Millisecond)
	tmr.Stop()
	tmr.ShowBatchTime()
	tmr.ShowTotalDuration()
	assert.Equal(t, int64(1), tmr.Index.Get())
	assert.Equal(t, int64(0), tmr.BatchDuractionSeconds())
	assert.Equal(t, int64(0), tmr.TotalDuractionSeconds())
	x := tmr.TotalDuration()
	assert.True(t, x.Seconds() < float64(1.0))
}
