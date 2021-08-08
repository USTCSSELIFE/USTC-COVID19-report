package main

import (
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/rs/zerolog"
	"os"
	"regexp"
	"time"
)

var (
	loginUrl = "https://passport.ustc.edu.cn/login?service=https%3A%2F%2Fweixine.ustc.edu.cn%2F2020%2Fcaslogin"
	username = ""
	password = ""
	logger   zerolog.Logger
)

func main() {
	file, err := os.OpenFile("./log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger = zerolog.New(file).With().Timestamp().Caller().Logger()

	for i := 0; i < 5; i++ {
		reportTime := report()
		if time.Now().Sub(reportTime) > time.Minute*10 {
			// if one doesn't report will automatically receive a mail by USTC,
			// so not necessarily to tell the failure
			logger.Error().Msgf("fail to report: %v\n", reportTime)
		} else {
			logger.Info().Msgf("succeed to report: %v\n", reportTime)
			return
		}
		logger.Info().Msgf("retry to report %d times\n", i+1)
	}
}

func report() time.Time {
	u := launcher.New().NoSandbox(true).Headless(false).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustIncognito().MustPage(loginUrl)
	page.MustElement("#username").MustInput(username)
	page.MustElement("#password").MustInput(password)
	page.MustElement("#login").MustClick()
	page.MustWaitLoad().MustElement("#report-submit-btn").MustClick()
	timeText, err := page.MustWaitLoad().MustElementX("//div[@id='daliy-report']//span//strong").Text()
	if err != nil {
		logger.Err(err)
	}

	return formatTime(timeText)
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
