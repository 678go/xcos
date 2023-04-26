package config

import (
	"golang.org/x/exp/slog"
	"gopkg.in/ini.v1"
)

type Bucket struct {
	Name      string `ini:"name"`
	SecretId  string `ini:"secretid"`
	SecretKey string `ini:"secretkey"`
	BucketUrl string `ini:"bucketurl"`
}

func InitBucket(r string) (*Bucket, error) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		slog.Error("Fail to read file", err)
		return nil, err
	}
	b := new(Bucket)
	_ = cfg.Section(r).MapTo(b)
	b.Name = r
	return b, err
}
