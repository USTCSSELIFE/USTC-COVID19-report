package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/rs/zerolog"
	"os"
	"regexp"
	"time"
)

const (
	loginUrl = "https://passport.ustc.edu.cn/login?service=https%3A%2F%2Fweixine.ustc.edu.cn%2F2020%2Fcaslogin"
	username = ""
	password = ""
)

var logger zerolog.Logger

func main() {
	file, err := os.OpenFile("./log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger = zerolog.New(file).With().Timestamp().Caller().Logger()

	for i := 0; i < 5; i++ {
		reportTime, ok := report()
		if !ok {
			logger.Info().Msgf("have reported at %v\n", reportTime)
			return
		}
		if isTimeValid(reportTime) {
			logger.Info().Msgf("succeed to report at %v\n", reportTime)
			return
		}
		// if one doesn't report will automatically receive a mail by USTC,
		// so not necessarily to tell the failure
		logger.Error().Msgf("fail to report: %v\n", reportTime)
		logger.Info().Msgf("retry to report %d times\n", i+1)
	}
}

func report() (time.Time, bool) {
	u := launcher.New().NoSandbox(true).Headless(false).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustIncognito().MustPage(loginUrl).MustWindowFullscreen()
	page.MustElement("#username").MustInput(username)
	page.MustElement("#password").MustInput(password)
	page.MustElement("#login").MustClick()

	timeText := page.MustWaitLoad().MustElementX("//div[@id='daliy-report']//span//strong").MustText()

	// if one has reported
	if !isTimeValid(formatTime(timeText)) {
		return formatTime(timeText), false
	}
	page.MustElement("#report-submit-btn").MustClick()
	timeText = page.MustWaitLoad().MustElementX("//div[@id='daliy-report']//span//strong").MustText()

	return formatTime(timeText), true
}

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
	return time.Now().Sub(reportTime) < time.Minute*10
}
