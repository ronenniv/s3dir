package s3types

import "io"

type Lister interface {
	PrintShort(io.Writer)
	PrintLong(io.Writer)
}
