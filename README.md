## ksyun
This package is a ksyun SDK for golang.
## Usage:
---------
```go
package api

import (
	"fmt"
	"github.com/928799934/ksyun/aws"
	"github.com/928799934/ksyun/s3"
	log "github.com/928799934/log4go.v1"
	"net/http"
	"path"
)

const (
	// 测试环境
	s3BucketTest = "test"
	s3UrlTest    = "https://test.ks3-cn-beijing.ksyun.com/"

	// 线上环境
	s3BucketOnline = "online"
	s3UrlOnline    = "https://online.ks3-cn-beijing.ksyun.com/"

	s3AccessKey = ""
	s3SecretKey = ""
)

var (
	// 默认测试环境
	s3Bucket = s3BucketTest
	s3Url    = s3UrlTest
)

type S3 struct {
	bucket *s3.Bucket
}

func InitVar(isOnline bool) {
	if isOnline {
		s3Bucket = s3BucketOnline
		s3Url = s3UrlOnline
	}
}

func NewS3() (*S3, error) {
	bucket := s3.NewS3(aws.Auth{
		AccessKey: s3AccessKey,
		SecretKey: s3SecretKey,
	}, aws.BEIJING).Bucket(s3Bucket)

	// s3 bucket
	return &S3{bucket}, nil
}

// 储存图片
func (s *S3) Put(data []byte, contType, fpath, name string) error {
	if err := s.bucket.Put(path.Join(fpath, name), data, contType, s3.PublicRead); err != nil {
		log.Error("bucket.Put(%s, %s) error(%v)", path.Join(fpath, name), contType, err)
		return err
	}
	return nil
}

// 获取图片
func (s *S3) Get(fpath, name string) string {
	return s3Url + path.Join(fpath, name)
}

func (s *S3) GetByWidthAndHeight(fpath string, width, height, mode int) string {
	return fmt.Sprintf("%s%s@base@tag=imgScale&w=%d&h=%d&m=%d", s3Url, fpath, width, height, mode)
}

func (s *S3) GetBySpecificParams(fpath string, m, w, h, q, c, f int) string {
	return fmt.Sprintf("%s%s@base@tag=imgScale&m=%d&w=%d&h=%d&q=%d&c=%d&f=%d", s3Url, fpath, m, w, h, q, c, f)
}

// 判断存在
func (s *S3) Exist(fpath, name string) bool {
	u := s.bucket.URL(path.Join(fpath, name))
	resp, err := http.Head(u)
	if err != nil {
		return false
	}
	return resp.StatusCode == http.StatusOK
}
```
