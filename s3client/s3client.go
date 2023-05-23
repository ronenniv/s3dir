package s3client

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ronenniv/s3dir/buckets"
	"github.com/ronenniv/s3dir/objects"
)

type S3Client struct {
	client *s3.Client
	cfg    *aws.Config
}

func NewS3Client() *S3Client {
	return &S3Client{}
}

func (s *S3Client) Connect() error {
	cfg, err := config.LoadDefaultConfig(context.Background()) // config.WithRegion(awsRegion),
	// config.WithSharedConfigProfile("rniv"),
	if err != nil {
		return err
	}
	s.cfg = &cfg
	s.client = s3.NewFromConfig(cfg)
	return nil
}

// ListBuckets return the list of all buckets
func (s *S3Client) ListBuckets() (*buckets.BucketList, error) {
	var bl buckets.BucketList

	input, err := s.client.ListBuckets(context.Background(), &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("got error in receiving list of buckets: %w", err)
	}
	bl.CopyAWSBucketListToBucketList(input.Buckets)

	return &bl, nil
}

// ListObjects will return list of all objects in bucketName
func (s *S3Client) ListObjects(bucketName string, prefix string) (*objects.Objects, error) {
	var bo objects.Objects

	input, err := s.client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
		Prefix:    aws.String(".git/"),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		var nsb *types.NoSuchBucket
		if errors.As(err, &nsb) {
			msg := fmt.Sprintf("%s: %s", bucketName, nsb.ErrorCode())
			bo.ErrMsg = &msg
		} else {
			return nil, fmt.Errorf("got error retrieving list of objects: %w", err)
		}
	}

	log.Printf("keycount %v commonprefix %#v input %#v\n", input.KeyCount, input.CommonPrefixes, input)
	bo.Objects = input.Contents

	return &bo, nil
}
