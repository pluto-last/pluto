package storage

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type AliOSS struct {
	client *oss.Client
	bucket *oss.Bucket

	Scheme          string
	AccessKey       string
	SecretKey       string
	Endpoint        string
	Bucket          string
	TokenTimeOutSec int64
}

type PolicyToken struct {
	AccessKeyId string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
	Callback    string `json:"callback"`
}

type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

func (ali *AliOSS) Init() error {
	var err error
	ali.client, err = oss.New(ali.Endpoint, ali.AccessKey, ali.SecretKey)
	if err != nil {
		return err
	}

	ali.bucket, err = ali.client.Bucket(ali.Bucket)
	if err != nil {
		return err
	}

	return nil
}

// Getuptoken 由服务端签名,web端直传oss
func (ali *AliOSS) Getuptoken(key ...string) string {

	// 过期时间
	expire_end := time.Now().Unix() + ali.TokenTimeOutSec
	var tokenExpire = get_gmt_iso8601(time.Now().Unix() + ali.TokenTimeOutSec)

	//create post policy json
	var config ConfigStruct
	config.Expiration = tokenExpire
	var condition []string
	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	config.Conditions = append(config.Conditions, condition)

	//calucate signature
	result, _ := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(ali.AccessKey))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	var policyToken PolicyToken
	policyToken.AccessKeyId = ali.AccessKey
	policyToken.Host = ali.MakeCdnOrigin()
	policyToken.Expire = expire_end
	policyToken.Signature = signedStr
	policyToken.Policy = debyte
	val, _ := json.Marshal(policyToken)

	return string(val)
}

func get_gmt_iso8601(expire_end int64) string {
	var tokenExpire = time.Unix(expire_end, 0).UTC().Format("2006-01-02T15:04:05Z")
	return tokenExpire
}

func (ali *AliOSS) Upload(mediaid string, reader io.Reader, size int64) error {
	return ali.bucket.PutObject(mediaid, reader)
}

func (ali *AliOSS) IsExist(mediaid string) (bool, error) {
	return ali.bucket.IsObjectExist(mediaid)
}

func (ali *AliOSS) Download(mediaid string) (io.ReadCloser, error) {
	return ali.bucket.GetObject(mediaid)
}

func (ali *AliOSS) DeleteCDNMedia(key string) error {
	return ali.bucket.DeleteObject(key)
}

func (ali *AliOSS) MakeCdnUrl(key string) string {
	return ali.MakeCdnOrigin() + "/" + key
}

func (ali *AliOSS) MakeCdnOrigin() string {
	return "https://" + ali.Bucket + "." + ali.Endpoint
}

func (ali *AliOSS) FileStat(mediaid string) (*Stat, error) {
	header, err := ali.bucket.GetObjectDetailedMeta(mediaid)
	if err != nil {
		return nil, err
	}

	fileStat := &Stat{}
	fileStat.Fsize, _ = strconv.ParseInt(header.Get("Content-Length"), 10, 64)
	fileStat.MimeType = header.Get("Content-Type")
	return fileStat, nil
}

func (ali *AliOSS) FetchImage(key, url string) error {

	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return ali.bucket.PutObject(key, resp.Body)
}

func (ali *AliOSS) FetchWeChatMedia(key, accesstoken, serverId string) error {
	mediaurl := fmt.Sprintf("http://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", accesstoken, serverId)
	err := ali.FetchImage(key, mediaurl)
	return err
}

func (ali *AliOSS) Pfop(key, fops, notifyurl string) (err error) {
	return nil
}

func (ali *AliOSS) SignRequest(req *http.Request) (token string, err error) {
	return "", nil

}

func (ali *AliOSS) GetAccessToken(data []byte) string {
	return ""
}

func (ali *AliOSS) GetAccessTokenWithData(data []byte) string {
	return ""
}
