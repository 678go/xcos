package cos

import "context"

type BucketServer interface {
	upload(context.Context, string) error
}
