package s3

import (
	"bytes"
	"git.liebaopay.com/pigs/public/ksyun/aws"
	ksaws "github.com/ks3sdklib/aws-sdk-go/aws"
	"github.com/ks3sdklib/aws-sdk-go/aws/credentials"
	"github.com/ks3sdklib/aws-sdk-go/service/s3"
	"io"
	"time"
)

type S3 struct {
	aws.Auth
	aws.Region
}

type Bucket struct {
	*s3.S3
	Name string
}

func NewS3(auth aws.Auth, region aws.Region) *S3 {
	return &S3{auth, region}
}

func (this *S3) Bucket(name string) *Bucket {
	s3 := s3.New(&ksaws.Config{
		Region:     this.Name,
		Endpoint:   this.S3Endpoint, //s3地址
		DisableSSL: this.DisableSSL, //是否禁用https
		Credentials: credentials.NewStaticCredentials(
			this.AccessKey,
			this.SecretKey,
			this.TokenKey,
		),
	})
	return &Bucket{s3, name}
}

type ACL string

const (
	Private           = ACL("private")
	PublicRead        = ACL("public-read")
	PublicReadWrite   = ACL("public-read-write")
	AuthenticatedRead = ACL("authenticated-read")
	BucketOwnerRead   = ACL("bucket-owner-read")
	BucketOwnerFull   = ACL("bucket-owner-full-control")
)

func (this *Bucket) Put(path string, data []byte, contType string, perm ACL) error {
	body := bytes.NewReader(data)
	return this.PutReader(path, body, contType, perm)
}

func (this *Bucket) PutReader(path string, r io.ReadSeeker, contType string, perm ACL) error {
	params := &s3.PutObjectInput{
		Bucket:      ksaws.String(this.Name),    // bucket名称
		Key:         ksaws.String(path),         // object key
		ACL:         ksaws.String(string(perm)), //权限，支持private(私有)，public-read(公开读)
		Body:        r,                          //要上传的内容
		ContentType: ksaws.String(contType),     //设置content-type
		Metadata:    map[string]*string{},
	}
	if _, err := this.PutObject(params); err != nil {
		return err
	}
	return nil
}

func (this *Bucket) URL(path string) string {
	params := &s3.GetObjectInput{
		Bucket: ksaws.String(this.Name), // bucket名称
		Key:    ksaws.String(path),      // object key
		//ResponseCacheControl:       aws.String("ResponseCacheControl"),//控制返回的Cache-Control header
		//ResponseContentDisposition: aws.String("ResponseContentDisposition"),//控制返回的Content-Disposition header
		//ResponseContentEncoding:    aws.String("ResponseContentEncoding"),//控制返回的Content-Encoding header
		//ResponseContentLanguage:    aws.String("ResponseContentLanguage"),//控制返回的Content-Language header
		//ResponseContentType: aws.String("image/png"), //控制返回的Content-Type header
	}

	expires := time.Hour // 1d

	resp, err := this.GetObjectPresignedUrl(params, expires) //第二个参数为外链过期时间，第二个参数为time.Duration类型
	if err != nil {
		return ""
	}
	return resp.String()
}
