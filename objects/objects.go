package objects

import (
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Objects struct {
	Objects []types.Object
	ErrMsg  *string
}

func (bo Objects) Len() int {
	return len(bo.Objects)
}
