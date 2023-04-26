package cos

import (
	"context"
	"fmt"
	"github.com/678go/xcos/config"
	"github.com/tencentyun/cos-go-sdk-v5"
	"golang.org/x/exp/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TenCentBucket struct {
	client  *cos.Client
	preSign *url.URL
	sigKey  string
}

func (b *TenCentBucket) Upload(ctx context.Context, file string) error {

	open, err := os.Open(file)
	if err != nil {
		slog.Error("read file err", err)
		return err
	}
	req, err := http.NewRequest(http.MethodPut, b.preSign.String(), open)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", b.sigKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	return nil
}

func NewTenCentBucket(bucket *config.Bucket) *TenCentBucket {
	u, _ := url.Parse(bucket.BucketUrl)
	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  os.Getenv(bucket.SecretId),
			SecretKey: os.Getenv(bucket.SecretKey),
		},
	})

	preSignedURL, err := client.Object.GetPresignedURL(context.Background(), http.MethodGet, "main", bucket.SecretId, bucket.SecretKey, 10*time.Minute, nil)
	if err != nil {
		slog.Error("get preSign failed", err)
		return nil
	}
	// 时间
	//s := carbon.Now().Timestamp()
	//d := carbon.Tomorrow().Timestamp()
	//h := sha1.New()
	//h.Write([]byte("get\\n/exampleobject\\n\\n\\n"))
	////sha1HttpString := h.Sum(nil)
	//
	//var hashFunc = sha1.New
	//h = hmac.New(hashFunc, []byte(bucket.SecretKey))
	//h.Write([]byte(fmt.Sprintf("%d;%d", s, d)))
	//signKey := h.Sum(nil)
	//a := fmt.Sprintf("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s&q-url-param-list=%s&q-signature=%v",
	//	bucket.SecretKey, fmt.Sprintf("%d;%d", s, d), fmt.Sprintf("%d;%d", s, d, "get\\n/exampleobject\\n\\n\\n", "delimiter;max-keys;prefix", signKey))
	//
	//fmt.Println(a)
	return &TenCentBucket{
		client:  client,
		preSign: preSignedURL,
		sigKey:  "q-sign-algorithm=sha1&q-ak=AKIDsdBUXKtEJAmAKFwPdcMr78uRjma9TYNx&q-sign-time=1682524476;1682528076&q-key-time=1682524476;1682528076&q-header-list=&q-url-param-list=&q-signature=645fd4f31e5573248edb586bcd32e977412d0816",
	}
}
