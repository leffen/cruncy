package cruncy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNames(t *testing.T) {
	var testCases = []struct {
		src         string
		numExpected int
	}{
		{"Test1,Test2,Test3", 3},
		{"Test1\nTest2\nTest3", 3},
		{"Test1\n  #Test2\nTest3", 2},
		{"", 0},
	}

	for _, tt := range testCases {
		items := ParseNames(tt.src)
		assert.Equal(t, len(items), tt.numExpected)
	}

}
