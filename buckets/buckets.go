package buckets

import (
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type BucketList struct {
	Buckets []types.Bucket
}

func (bl *BucketList) PrintShort(w io.Writer) {
	for _, bucket := range bl.Buckets {
		fmt.Fprint(w, *bucket.Name, "\n")
	}
}

func (bl *BucketList) PrintLong(w io.Writer) {
	for _, bucket := range bl.Buckets {
		fmt.Fprintf(w, "%s\t%s\n", *bucket.CreationDate, *bucket.Name)
	}
}
