package main

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

type awsService struct {
	svc *iam.IAM
}

func (s *awsService) Has2FA(userID string) (bool, error) {
	params := &iam.ListMFADevicesInput{
		UserName: aws.String(userID),
	}
	resp, err := s.svc.ListMFADevices(params)

	if err != nil {
		return false, err
	}

	return len(resp.MFADevices) > 0, nil
}

func newAWSService(config string) (Service, error) {
	configParts := strings.Split(config, ":")

	svc := iam.New(session.New(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(configParts[0], configParts[1], ""),
	}))

	return &awsService{svc: svc}, nil
}

func init() {
	registerService("aws", "https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html", newAWSService)
}
