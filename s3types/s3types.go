package s3types

import "io"

type BucketLister interface {
	PrintShort(io.Writer)
	PrintLong(io.Writer)
}

type ObjectLister interface {
	PrintShort(io.Writer)
	PrintLong(io.Writer)
}
