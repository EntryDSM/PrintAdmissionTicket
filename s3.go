package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"log"
	"os"
)

var (
	Downloader *s3manager.Downloader
	Bucket     string
)

func InitS3() {
	Bucket = os.Getenv("AWS_S3_BUCKET")
	if Bucket == "" {
		log.Panicln("[ERROR] Cannot Find 'AWS_S3_BUCKET' Property.")
	}

	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		log.Panicln(err)
	}
	Downloader = s3manager.NewDownloader(cfg)

	cacheDir := "./cache"
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		os.Mkdir(cacheDir, 0700)
	}
}

func SaveIfEmpty(key string) {
	if key == "" {
		log.Fatal()
	}

	filename := "./cache/" + key
	if exists(filename) {
		return
	}

	file, err := os.Create("./cache/" + key)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	_, err = Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Println(err)
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
