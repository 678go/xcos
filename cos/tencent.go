package cos

import (
	"context"
	"fmt"
	"github.com/678go/xcos/util"
	"github.com/schollz/progressbar/v3"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"os"
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
		Secretkey: "a",
		Bucketurl: "https://test-test.cos.ap-nanjing.myqcloud.com",
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
		progressbar.OptionSetDescription("running..."),
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
	//split := strings.Split(source, "/")
	//dest := split[len(split)-1]
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
	bucketObj, _, err := t.getObjectList(source)
	if err != nil {
		slog.Error("获取存储桶文件列表失败", "err", err)
	}
	for _, object := range bucketObj {
		localPath, bucketPath, err := util.DownloadPathFixed(object.Key)
		if err != nil {
			return err
		}
		download, err := t.client.Object.Download(ctx, bucketPath, localPath, opt)
		if download.StatusCode == 404 {
			slog.Warn("not found file", "filename", path)
			return err
		}
		if err != nil {
			slog.Error("download file error,", err, "filename:", path)
			return err
		}
	}
	check = true

	return nil
}

func (t *TenCentBucket) UploadFolder(ctx context.Context, filepath string) error {
	fileInfo, _ := os.Stat(filepath)
	if !fileInfo.IsDir() {
		if err := t.Upload(ctx, filepath); err != nil {
			slog.Error("单个文件上传失败,", "err", err)
			return err
		}
		return nil
	}
	for _, f := range util.GetLocalFilesListRecursive(filepath) {
		if err := t.Upload(ctx, filepath+"/"+f); err != nil {
			return err
		}
	}
	return nil
}

func (t *TenCentBucket) getObjectList(bucketPath string) (objects []cos.Object, commonPrefixes []string, err error) {
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:    bucketPath, // prefix 表示要查询的文件夹
		Delimiter: "",         // deliter 表示分隔符, 设置为/表示列出当前目录下的 object, 设置为空表示列出所有的 object
		MaxKeys:   1000,
		Marker:    "",
	}
	isTruncated := true
	marker = ""
	for isTruncated {
		opt.Marker = marker

		res, _, err := t.client.Bucket.Get(context.Background(), opt)
		if err != nil {
			slog.Error("err", err)
			os.Exit(1)
		}

		objects = append(objects, res.Contents...)
		commonPrefixes = res.CommonPrefixes

		isTruncated = res.IsTruncated
		marker = res.NextMarker

	}
	return
}
