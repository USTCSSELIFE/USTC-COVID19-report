package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTimeValid(t *testing.T) {
	assert.Equal(t, false, isTimeValid(formatTime("2006-01-02 15:04:05")))
	//assert.Equal(t, true, isTimeValid(formatTime("2021-08-11 16:32:05")))
}

func TestHaveReported(t *testing.T) {
	assert.Equal(t, true, haveReported(formatTime("2021-08-11 01:04:05")))
	assert.Equal(t, true, haveReported(formatTime("2021-08-11 15:04:05")))
	assert.Equal(t, false, haveReported(formatTime("2021-08-12 15:04:05")))
}
