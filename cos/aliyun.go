package cos

import (
	"context"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"golang.org/x/exp/slog"
	"log"
	"os"
	"strings"
)

type AliBucket struct {
	BaseBucket
	ProgressListener
	bucket *oss.Bucket
}

func (a *AliBucket) InitClient() {
	a.BaseBucket = BaseBucket{
		Secretid:  "a",
		Secretkey: "b",
		Bucketurl: "https://oss-cn-ccccc.aliyuncs.com",
	}
	client, err := oss.New(a.Bucketurl, a.Secretid, a.Secretkey)
	if err != nil {
		slog.Error("init client fail", "err", err)
		return
	}
	bucket, err := client.Bucket("test-aaa")
	if err != nil {
		slog.Error("init bucket fail", "err", err)
	}
	a.bucket = bucket
}

// Upload 如果上传为文件夹 则里面不能包括特殊文件夹 比如隐藏文件夹
func (a *AliBucket) Upload(ctx context.Context, file string) error {
	var bucketPath string
	// ./
	if strings.HasPrefix(strings.TrimLeft(file, "."), "/") {
		bucketPath = fmt.Sprintf("tmp%s", strings.TrimLeft(file, "."))
	} else {
		// 空的
		bucketPath = fmt.Sprintf("tmp/%s", file)
	}
	fileInfo, err := os.Stat(file)
	if err != nil {
		slog.Error("upload file fail", "info", "存在特殊文件夹", err, err.Error())
	}
	if !fileInfo.IsDir() {
		if err := a.bucket.PutObjectFromFile(bucketPath, file, oss.Progress(&a.ProgressListener)); err != nil {
			slog.Error("upload file failed,", "err", err)
			return err
		}
		fmt.Println("上传成功!")
	}
	var dirEntries []string
	for _, entry := range FileForEachComplete(file, dirEntries) {
		if err := a.bucket.PutObjectFromFile(bucketPath, entry, oss.Progress(&a.ProgressListener)); err != nil {
			slog.Error("upload file failed,", "info", "存在特殊文件夹", "err", err)
			return err
		}
	}
	return nil
}

func (a *AliBucket) DownloadFile(ctx context.Context, path string) error {
	source := fmt.Sprintf("tmp/%s", strings.TrimLeft(path, "/"))
	split := strings.Split(source, "/")
	dest := split[len(split)-1]
	if err := a.bucket.GetObjectToFile(source, dest, oss.Progress(&a.ProgressListener)); err != nil {
		slog.Error("download file failed,", "err", err)
		return err
	}
	fmt.Println("下载成功!")
	return nil
}

type ProgressListener struct {
}

func (listener *ProgressListener) ProgressChanged(event *oss.ProgressEvent) {
	switch event.EventType {
	case oss.TransferStartedEvent:
		fmt.Printf("Transfer Started, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferDataEvent:
		fmt.Printf("\rTransfer Data, ConsumedBytes: %d, TotalBytes %d, %d%%.",
			event.ConsumedBytes, event.TotalBytes, event.ConsumedBytes*100/event.TotalBytes)
	case oss.TransferCompletedEvent:
		fmt.Printf("\nTransfer Completed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	case oss.TransferFailedEvent:
		fmt.Printf("\nTransfer Failed, ConsumedBytes: %d, TotalBytes %d.\n",
			event.ConsumedBytes, event.TotalBytes)
	default:
	}
}

func FileForEachComplete(fileFullPath string, f []string) []string {
	files, err := os.ReadDir(fileFullPath)
	if err != nil {
		fmt.Println(1)
		fmt.Println(err)
		log.Fatal(err)
	}
	for _, file := range files {
		if file.IsDir() {
			f = FileForEachComplete(fileFullPath+"/"+file.Name(), f)
		} else {
			f = append(f, fileFullPath+"/"+file.Name())
		}
	}
	return f
}
