package main

import (
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

func InitS3Downloader() (*s3manager.Downloader, error) {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		return nil, err
	}

	Downloader := s3manager.NewDownloader(cfg)
	return Downloader, nil
}
