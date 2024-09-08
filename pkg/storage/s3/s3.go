package s3

import (
	"bank-authentication-system/pkg/config"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type S3 struct {
	Cfg     config.S3
	Session *session.Session
}

func NewSession(cfg config.S3) (*S3, error) {
	var s3 S3

	newSession, err := session.NewSession(
		&aws.Config{
			Region: aws.String(cfg.Region),
			Credentials: credentials.NewStaticCredentials(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
			Endpoint: &cfg.Endpoint,
		},
	)
	if err != nil {
		return nil, err
	}

	s3.Session = newSession
	s3.Cfg = cfg

	return &s3, nil
}

func GetImageFromS3(imageKey string) string {
	baseUrl := "https://publisher.s3.ir-thr-at1.arvanstorage.ir"
	urlStr := fmt.Sprintf("%s/%s", baseUrl, imageKey)

	return urlStr
}
