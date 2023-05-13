package s3client

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/ronenniv/s3dir/buckets"
	"github.com/ronenniv/s3dir/objects"
)

type S3Client struct {
	clientData *s3.Client
	cfg        *aws.Config
	ctx        context.Context
}

func NewS3Client() *S3Client {
	s := S3Client{
		ctx: context.Background(),
	}
	return &s
}

const awsEndpoint = "http://localhost:4566"
const awsRegion = "us-east-1"

func (s *S3Client) Connect() error {
	customResolver := aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})
	cfg, err := config.LoadDefaultConfig(s.ctx,
		// config.WithRegion(awsRegion),
		config.WithSharedConfigProfile("localstack"),
		config.WithEndpointResolver(customResolver),
	)
	if err != nil {
		err = fmt.Errorf("in connect(): %w", err)
		return err
	}
	s.cfg = &cfg
	s.clientData = s3.NewFromConfig(cfg)
	return nil
}

func (s *S3Client) ListBuckets() (*buckets.BucketList, error) {
	var bl buckets.BucketList

	input, err := s.clientData.ListBuckets(s.ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("got error in receiving list of buckets: %w", err)
	}
	bl.CopyAWSBucketListToBucketList(input.Buckets)

	return &bl, nil
}

func (s *S3Client) ListObjects(bucketName string) (*objects.BucketObjects, error) {
	var bo objects.BucketObjects

	input, err := s.clientData.ListObjectsV2(s.ctx, &s3.ListObjectsV2Input{Bucket: &bucketName})
	if err != nil {
		var nsb *types.NoSuchBucket
		if errors.As(err, &nsb) {
			msg := fmt.Sprintf("%s: %s", bucketName, nsb.ErrorCode())
			bo.ErrMsg = &msg
		} else {
			return nil, fmt.Errorf("got error retrieving list of objects: %w", err)
		}
	}
	bo.Objects = input.Contents

	return &bo, nil
}
