package main

import (
	"github.com/pkg/errors"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ocr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ocr/v20181119"
)

const (
	secretId  = ""
	secretKey = ""
)

var lowConfidenceError = errors.New("confidence is lower than 85")

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

	validateCode := *response.Response.TextDetections[0].DetectedText
	confidence := *response.Response.TextDetections[0].Confidence
	if confidence < 85 {
		return "", lowConfidenceError
	}

	return validateCode, nil
}
