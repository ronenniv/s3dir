package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	clientData *s3.Client
	cfg        *aws.Config
	ctx        context.Context
}

func New() *S3Client {
	s := S3Client{
		ctx: context.Background(),
	}
	return &s
}

func (s *S3Client) connect() error {
	cfg, err := config.LoadDefaultConfig(s.ctx)
	if err != nil {
		err = fmt.Errorf("in connect(): %w", err)
		return err
	}
	s.cfg = &cfg
	s.clientData = s3.NewFromConfig(cfg)
	return nil
}

func (s *S3Client) listObjects(bucketName string) {
	input, err := s.clientData.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		fmt.Println("Got error retrieving list of objects:")
		fmt.Println(err)
		return
	}

	fmt.Println("Objects in " + bucketName + ":")

	for _, item := range input.Contents {
		fmt.Println("Name:          ", *item.Key)
		fmt.Println("Last modified: ", *item.LastModified)
		fmt.Println("Size:          ", item.Size)
		fmt.Println("Storage class: ", item.StorageClass)
		fmt.Println("")
	}

	fmt.Println("Found", len(input.Contents), "items in bucket", bucketName)
}

func main() {
	// bucket := flag.String("b", "", "The name of the bucket")
	// flag.Parse()
	s3client := New()
	if err := s3client.connect(); err != nil {
		log.Fatal(err)
	}
	s3client.listObjects("elasticbeanstalk-us-east-1-899792839109")

}
