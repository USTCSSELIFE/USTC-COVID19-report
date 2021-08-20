package main

import (
	"encoding/base64"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"os"
	"time"
)

const (
	loginUrl = "https://passport.ustc.edu.cn/login?service=https%3A%2F%2Fweixine.ustc.edu.cn%2F2020%2Fcaslogin"
	username = ""
	password = ""
)

var logger zerolog.Logger
var haveReportError = errors.New("have reported")

func main() {
	file, err := os.OpenFile("./log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	logger = zerolog.New(file).With().Timestamp().Caller().Logger()

	for i := 0; i < 5; i++ {
		reportTime, err := report()
		if err != nil {
			if errors.Is(err, haveReportError) {
				logger.Warn().Time("lastReportTime", reportTime).Msg(err.Error())
				return
			}
			logger.Error().Msgf("fail to report: %v", err)
			return
		}
		if isTimeValid(reportTime) {
			logger.Info().Msgf("succeed to report at %v", reportTime)
			return
		}
		// if one doesn't report will automatically receive a mail by USTC,
		// so not necessarily to tell the failure
		logger.Error().Time("lastReportTime", reportTime).
			Int("retry times", i+1).
			Msg("fail to report")
	}
}

func report() (time.Time, error) {
	u := launcher.New().NoSandbox(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustIncognito().MustPage(loginUrl).MustWindowFullscreen()
	page.MustElement("#username").MustInput(username)
	page.MustElement("#password").MustInput(password)
	bin := page.MustElementX("//img[@class='validate-img']").MustResource()
	validateCode, err := getValidateCode(base64.StdEncoding.EncodeToString(bin))
	if err != nil {
		return time.Time{}, err
	}
	page.MustElement("#validate").MustInput(validateCode)
	page.MustElement("#login").MustClick()
	timeText := page.MustWaitLoad().MustElementX("//div[@id='daliy-report']//span//strong").MustText()

	// if one has reported
	t, err := formatTime(timeText)
	if err != nil {
		return time.Time{}, err
	}
	if haveReported(t) {
		return t, haveReportError
	}

	page.MustElement("#report-submit-btn").MustClick()
	timeText = page.MustWaitLoad().MustElementX("//div[@id='daliy-report']//span//strong").MustText()
	t, err = formatTime(timeText)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}
