package main

import (
	"github.com/golang-module/carbon"
	"time"
)

func formatTime(timeText string) (time.Time, error) {
	c := carbon.ParseByFormat(timeText, "*上次上报时间：Y-m-d H:i:s，请每日按时打卡", carbon.Shanghai)
	if c.Error != nil {
		return time.Time{}, c.Error
	}
	return c.Carbon2Time(), nil
}

func isTimeValid(reportTime time.Time) bool {
	return time.Since(reportTime) < time.Minute*10
}

func haveReported(reportTime time.Time) bool {
	cReportTime := carbon.Time2Carbon(reportTime)
	cNow := carbon.Time2Carbon(time.Now())

	if cReportTime.IsFuture() {
		return false
	}
	if cReportTime.BetweenIncludedBoth(cNow.StartOfDay(), cNow) {
		return true
	}

	return false
}
