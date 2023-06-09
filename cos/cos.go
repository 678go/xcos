package cos

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
)

type BucketServer interface {
	InitClient()
	Upload(context.Context, string) error
	UploadFolder(context.Context, string) error
	DownloadFile(ctx context.Context, path string) error
}

var buckets = map[string]BucketServer{
	"tencent": &TenCentBucket{},
	//"aliyun":  &AliBucket{},
}

func NewBucket(t string) (BucketServer, error) {
	server, ok := buckets[t]
	if !ok {
		slog.Error("not found bucket", "bucketname:", t)
		return nil, fmt.Errorf("not found bucket")
	}
	return server, nil
}

type BaseBucket struct {
	Secretid  string
	Secretkey string
	Bucketurl string
}
