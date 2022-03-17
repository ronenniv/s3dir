package objects

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketObjects struct {
	Objects []types.Object
	ErrMsg  *string
}

func (bo *BucketObjects) PrintShort(w io.Writer) {
	if bo.ErrMsg != nil {
		fmt.Fprintf(w, "%s\n", *bo.ErrMsg)
	} else {
		for _, obj := range bo.Objects {
			fmt.Fprintf(w, "%s\n", *obj.Key)
		}
	}
}

func (bo *BucketObjects) PrintLong(w io.Writer) {
	if bo.ErrMsg != nil {
		fmt.Fprintf(w, "%s\n", *bo.ErrMsg)
	} else {
		for _, obj := range bo.Objects {
			fmt.Fprintf(w, "%s\t%d\t%s\t%s\n", obj.StorageClass, obj.Size, *obj.LastModified, *obj.Key)
		}
	}
}
