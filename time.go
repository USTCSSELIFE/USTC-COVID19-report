package main

import (
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
	return time.Since(reportTime).Hours()/24 >= 0 && time.Since(reportTime).Hours()/24 <= 1
}
