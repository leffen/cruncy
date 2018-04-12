package cruncy

import (
	"testing"

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
