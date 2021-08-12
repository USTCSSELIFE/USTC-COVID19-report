package main

import (
	"github.com/golang-module/carbon"
	"regexp"
	"time"
)

func formatTime(timeText string) time.Time {
	regx := regexp.MustCompile("2021-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}")
	reportTimeString := regx.FindString(timeText)

	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		logger.Err(err)
	}
	reportTime, err := time.ParseInLocation("2006-01-02 15:04:05", reportTimeString, location)
	if err != nil {
		logger.Err(err)
	}
	return reportTime

}

func isTimeValid(reportTime time.Time) bool {
	return time.Since(reportTime) < time.Minute*10
}

func haveReported(reportTime time.Time) bool {
	cReportTime := carbon.Time2Carbon(reportTime)
	cNow := carbon.Time2Carbon(time.Now())

	if time.Since(reportTime) <= 0 {
		return false
	}
	if cReportTime.BetweenIncludedBoth(cNow.StartOfDay(), cNow) {
		return true
	}

	return false
}
