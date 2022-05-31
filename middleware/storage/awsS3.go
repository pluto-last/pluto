package storage

import (
	"fmt"
	"io"
	"net/http"
	"pluto/utils"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSS3 struct {
	AccessKey       string
	SecretKey       string
	Bucket          string
	Region          string // 地区 香港：ap-east-1， 参考 endpoints.ApEast1RegionID
	Endpoint        string
	TokenTimeOutSec int64
	Session         *session.Session
	Svc             *s3.S3
}

func (awss3 *AWSS3) Init() error {
	var err error
	awss3.Session, err = session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awss3.AccessKey, awss3.SecretKey, ""),
		Region:      aws.String(awss3.Region),
	})
	if err != nil {
		return err
	}

	// Create S3 service client
	awss3.Svc = s3.New(awss3.Session)

	return awss3.bucketInit()
}

func (awss3 *AWSS3) Upload(key string, reader io.Reader, size int64) error {
	uploader := s3manager.NewUploader(awss3.Session)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awss3.Bucket),
		Key:    aws.String(key),
		Body:   reader,
		ACL:    aws.String(s3.BucketCannedACLPublicReadWrite),
	})
	return err
}

func (awss3 *AWSS3) bucketInit() error {
	// 获取当前bucket是否存在
	bucketsResp, err := awss3.Svc.ListBuckets(&s3.ListBucketsInput{})
	for _, bucket := range bucketsResp.Buckets {
		fmt.Printf("current bucket %q ， need bucket name %q\n", *bucket.Name, awss3.Bucket)
		if *bucket.Name == awss3.Bucket {
			fmt.Printf("bucket %q 已存在\n", awss3.Bucket)
			return nil
		}
	}

	params := &s3.CreateBucketInput{
		Bucket: aws.String(awss3.Bucket),
		ACL:    aws.String(s3.BucketCannedACLPublicReadWrite),
	}

	_, err = awss3.Svc.CreateBucket(params)
	if err != nil {
		return err
	}

	// Wait until bucket is created before finishing
	fmt.Printf("Waiting for bucket %q to be created...\n", awss3.Bucket)

	err = awss3.Svc.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(awss3.Bucket),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Bucket %q successfully created\n", awss3.Bucket)
	return nil
}

func (awss3 *AWSS3) Getuptoken(key ...string) string {
	savekey := utils.RandStr(10)
	if len(key) > 0 {
		savekey = key[0]
	}

	params := &s3.PutObjectInput{
		Bucket: aws.String(awss3.Bucket), // Required
		Key:    aws.String(savekey),
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	}

	putReq, _ := awss3.Svc.PutObjectRequest(params)

	url, err := putReq.Presign(time.Duration(awss3.TokenTimeOutSec) * time.Second)
	if err != nil {
		panic(err)
	}

	return url
}

func (awss3 *AWSS3) IsExist(mediaid string) (bool, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(awss3.Bucket), // Required
		Key:    aws.String(mediaid),
	}
	_, err := awss3.Svc.GetObject(params)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (awss3 *AWSS3) FileStat(key string) (*Stat, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(awss3.Bucket), // Required
		Key:    aws.String(key),
	}
	_, err := awss3.Svc.GetObject(params)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (awss3 *AWSS3) FetchImage(key, url string) error {
	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 上传
	uploader := s3manager.NewUploader(awss3.Session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awss3.Bucket),
		Key:    aws.String(key),
		Body:   resp.Body,
		ACL:    aws.String(s3.BucketCannedACLPublicReadWrite),
	})

	return err
}

func (awss3 *AWSS3) FetchWeChatMedia(key, accesstoken, serverId string) error {
	mediaurl := fmt.Sprintf("http://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", accesstoken, serverId)

	// 下载url
	resp, err := http.Get(mediaurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 上传
	uploader := s3manager.NewUploader(awss3.Session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(awss3.Bucket),
		Key:    aws.String(key),
		Body:   resp.Body,
		ACL:    aws.String(s3.BucketCannedACLPublicReadWrite),
	})

	return err
}

func (awss3 *AWSS3) Pfop(key, fops, notifyurl string) (err error) {
	return nil
}

func (awss3 *AWSS3) DeleteCDNMedia(key string) error {
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(awss3.Bucket), // Required
		Key:    aws.String(key),
	}
	_, err := awss3.Svc.DeleteObject(params)
	return err

}

func (awss3 *AWSS3) MakeCdnUrl(key string) string {
	url := "https://%s.s3-%s.amazonaws.com/%s"
	return fmt.Sprintf(url, awss3.Bucket, awss3.Region, key)
}

func (awss3 *AWSS3) MakeCdnOrigin() string {
	url := "https://%s.s3-%s.amazonaws.com/"
	return fmt.Sprintf(url, awss3.Bucket, awss3.Region)
}

func (awss3 *AWSS3) SignRequest(req *http.Request) (token string, err error) {
	return "", nil
}

func (awss3 *AWSS3) GetAccessToken(data []byte) string {
	return ""
}

func (awss3 *AWSS3) GetAccessTokenWithData(data []byte) string {
	return ""
}

func (awss3 *AWSS3) Download(key string) (io.ReadCloser, error) {
	url := awss3.MakeCdnUrl(key)

	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
