package main

import (
	"github.com/golang-module/carbon"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeValid(t *testing.T) {
	assert.Equal(t, true, isTimeValid(time.Now()))
	assert.Equal(t, true, isTimeValid(carbon.Time2Carbon(time.Now()).AddMinute().Carbon2Time()))
}

func TestHaveReported(t *testing.T) {
	assert.Equal(t, false, haveReported(carbon.Time2Carbon(time.Now()).Yesterday().Carbon2Time()))
	assert.Equal(t, true, haveReported(carbon.Time2Carbon(time.Now()).StartOfMinute().Carbon2Time()))
	assert.Equal(t, false, haveReported(carbon.Time2Carbon(time.Now()).Tomorrow().Carbon2Time()))
	assert.Equal(t, false, haveReported(carbon.Time2Carbon(time.Now()).AddYear().Carbon2Time()))
}
