package objects

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketObjects struct {
	Objects []types.Object
	ErrMsg  *string
}

func (bo BucketObjects) Len() int {
	return len(bo.Objects)
}
