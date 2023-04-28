package cos

import (
	"context"
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TenCentBucket struct {
	BaseBucket
	client *cos.Client
	bar    *progressbar.ProgressBar
}

func (t *TenCentBucket) InitClient() {
	t.BaseBucket = BaseBucket{
		Secretid:  "a",
		Secretkey: "b",
		Bucketurl: "c",
	}
	u, _ := url.Parse(t.Bucketurl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  t.Secretid,
			SecretKey: t.Secretkey,
		},
	})
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetWidth(10),
		progressbar.OptionSetDescription("uploading..."),
		progressbar.OptionShowCount(),
		progressbar.OptionSpinnerCustom([]string{"🐰", "🐰", "🥕", "🥕"}),
	)
	t.client = client
	t.bar = bar
}

// Upload 默认上传路径buckeurl/tmp/文件
func (t *TenCentBucket) Upload(ctx context.Context, file string) error {
	var bucketPath string
	// ./
	if strings.HasPrefix(strings.TrimLeft(file, "."), "/") {
		bucketPath = fmt.Sprintf("tmp%s", strings.TrimLeft(file, "."))
	} else {
		// 空的
		bucketPath = fmt.Sprintf("tmp/%s", file)
	}

	var check bool
	go func() {
		for {
			time.Sleep(1 * time.Second)
			_ = t.bar.Add(-1)
			if check {
				break
			}
		}
	}()

	opt := &cos.MultiUploadOptions{
		ThreadPoolSize: 5,
	}
	upload, response, err := t.client.Object.Upload(ctx, bucketPath, file, opt)
	if err != nil {
		slog.Error("upload file,", "err", err, "file", file)
		return err
	}
	if response.StatusCode == 200 {
		check = true
	}

	fmt.Println("\t", upload.Location)
	return nil
}
func (t *TenCentBucket) DownloadFile(ctx context.Context, path string) error {
	opt := &cos.MultiDownloadOptions{
		ThreadPoolSize: 5,
	}
	source := fmt.Sprintf("tmp/%s", strings.TrimLeft(path, "/"))
	split := strings.Split(source, "/")
	dest := split[len(split)-1]
	var check bool
	go func() {
		for {
			time.Sleep(1 * time.Second)
			_ = t.bar.Add(-1)
			if check {
				break
			}
		}
	}()
	download, err := t.client.Object.Download(ctx, source, dest, opt)
	if download.StatusCode == 404 {
		slog.Warn("not found file", "filename", path)
		return err
	}
	if err != nil {
		slog.Error("download file error,", err, "filename:", path)
		return err
	}
	if download.StatusCode == 200 {
		check = true
	}
	return nil
}
