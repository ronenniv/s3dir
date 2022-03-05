package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/ronenniv/s3dir/buckets"
	"github.com/ronenniv/s3dir/objects"
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

func (s *S3Client) listBuckets() (*buckets.BucketList, error) {
	var bl buckets.BucketList
	input, err := s.clientData.ListBuckets(s.ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("got error in receiving list of buckets: %w", err)
	}
	bl.Buckets = input.Buckets

	return &bl, nil
}

func (s *S3Client) listObjects(bucketName string) (*objects.BucketObjects, error) {
	var bo objects.BucketObjects
	input, err := s.clientData.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		return nil, fmt.Errorf("got error retrieving list of objects: %w", err)
	}
	bo.Objects = input.Contents

	return &bo, nil
}

func main() {
	// bucket := flag.String("b", "", "The name of the bucket")
	// flag.Parse()
	s3client := New()
	if err := s3client.connect(); err != nil {
		log.Fatal(err)
	}
	buckets, _ := s3client.listBuckets()
	buckets.PrintShort(os.Stdout)
	buckets.PrintLong(os.Stdout)

	obj, _ := s3client.listObjects("tess-qa-checks")
	obj.PrintShort(os.Stdout)
	obj.PrintLong(os.Stdout)

}
