package main

import (
	"github.com/pkg/errors"
	"strconv"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

const (
	secretId  = ""
	secretKey = ""
)

var (
	lowConfidenceError = errors.New("confidence is lower than 95")
	detectedTextError  = errors.New("detected text is not a pure number")
)

func getValidateCode(imageBase64 string) (string, error) {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "ocr.tencentcloudapi.com"
	client, _ := ocr.NewClient(credential, "ap-shanghai", cpf)

	request := ocr.NewGeneralBasicOCRRequest()
	request.ImageBase64 = common.StringPtr(imageBase64)

	response, err := client.GeneralBasicOCR(request)
	if err != nil {
		return "", err
	}
	logger.Info().Msg(response.ToJsonString())

	validateCode := *response.Response.TextDetections[0].DetectedText
	if _, err = strconv.Atoi(validateCode); err != nil {
		return "", detectedTextError
	}
	confidence := *response.Response.TextDetections[0].Confidence
	if confidence < 95 {
		return "", lowConfidenceError
	}

	return validateCode, nil
}
