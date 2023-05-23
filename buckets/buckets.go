package buckets

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Bucket struct {
	Name         string
	CreationDate time.Time
}

type BucketList struct {
	Buckets []Bucket
}

func CopyAWSBucketToBucket(bucket types.Bucket) Bucket {
	b := Bucket{}
	b.Name = *bucket.Name
	b.CreationDate = *bucket.CreationDate
	return b
}

func (bl *BucketList) CopyAWSBucketListToBucketList(bucketList []types.Bucket) {
	for _, bucket := range bucketList {
		bl.Buckets = append(bl.Buckets, CopyAWSBucketToBucket(bucket))
	}
}

func (bl BucketList) Len() int {
	return len(bl.Buckets)
}
