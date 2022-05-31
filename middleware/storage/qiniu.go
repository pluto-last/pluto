package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"pluto/utils"
	"sync"
	"time"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
)

type CDN struct {
	bucketManager *storage.BucketManager
	config        *storage.Config
	uploader      *storage.ResumeUploaderV2
	mac           *auth.Credentials
	once          *sync.Once

	AccessKey       string
	SecretKey       string
	Bucket          string
	Domain          string
	SizeLimit       int64
	TokenTimeOutSec int64
	Zone            int // 您空间(Bucket)所在的区域 0:华东机房:1:华北机房;2:华南机房;3北美机房
	Pipeline        []string
	Scheme          string
}

func (cdn *CDN) Init() {

	cdn.config = &storage.Config{
		Zone:          cdn.getZone(),
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	cdn.uploaderInit()
	cdn.BucketInit()
}

func (cdn *CDN) BucketInit() {
	cdn.mac = auth.New(cdn.AccessKey, cdn.SecretKey)
	cdn.bucketManager = storage.NewBucketManager(cdn.mac, cdn.config)
}

func (cdn *CDN) uploaderInit() {
	cdn.uploader = storage.NewResumeUploaderV2(cdn.config)
}

func (cdn *CDN) getZone() *storage.Zone {
	zone := &storage.ZoneHuanan
	switch cdn.Zone {
	case 0:
		zone = &storage.ZoneHuadong
	case 1:
		zone = &storage.ZoneHuabei
	case 2:
		zone = &storage.ZoneHuanan
	case 3:
		zone = &storage.ZoneBeimei
	}

	return zone
}

func (cdn *CDN) Getuptoken(key ...string) string {
	savekey := utils.RandStr(10)
	if len(key) > 0 {
		savekey = key[0]
	}
	putPolicy := &storage.PutPolicy{
		Scope:      cdn.Bucket,
		Expires:    uint64(time.Now().Unix() + cdn.TokenTimeOutSec),
		FsizeLimit: cdn.SizeLimit,
		SaveKey:    savekey,
		ReturnBody: "{\"type\": $(mimeType),\"key\":$(key),\"link\":\"" + cdn.Scheme + "://" + cdn.Domain + "/$(key)\"}",
	}
	return putPolicy.UploadToken(cdn.mac)
}

func (cdn *CDN) Upload(mediaid string, reader io.Reader, size int64) error {
	// 上传凭证
	upToken := cdn.Getuptoken()
	ret := &storage.PutRet{}

	//  io.Reader to io.ReaderAt
	buff := bytes.NewBuffer([]byte{})
	size, err := io.Copy(buff, reader)
	if err != nil {
		return err
	}

	return cdn.uploader.Put(context.TODO(), ret, upToken, mediaid, bytes.NewReader(buff.Bytes()), size, nil)
}

func (cdn *CDN) IsExist(mediaid string) (bool, error) {
	_, err := cdn.bucketManager.Stat(cdn.Bucket, mediaid)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (cdn *CDN) FileStat(mediaid string) (*Stat, error) {
	entry, err := cdn.bucketManager.Stat(cdn.Bucket, mediaid)
	if err != nil {
		return nil, err
	}
	return &Stat{
		Hash:     entry.Hash,
		Fsize:    entry.Fsize,
		PutTime:  entry.PutTime,
		MimeType: entry.MimeType,
	}, nil
}

func (cdn *CDN) FetchImage(key, url string) error {
	_, err := cdn.bucketManager.Fetch(url, cdn.Bucket, key)
	return err
}

func (cdn *CDN) FetchWeChatMedia(key, accesstoken, serverId string) error {
	mediaurl := fmt.Sprintf("http://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", accesstoken, serverId)
	err := cdn.FetchImage(key, mediaurl)
	return err
}

func (cdn *CDN) Pfop(key, fops, notifyurl string) (err error) {
	ctx := context.Background()
	form := make(map[string][]string)
	form["bucket"] = []string{cdn.Bucket}
	form["key"] = []string{key}
	form["fops"] = []string{fops}
	form["notifyURL"] = []string{notifyurl}
	//form["force"] = []string{}
	form["pipeline"] = cdn.Pipeline //https://portal.qiniu.com/mps/upload
	//bucket=<urlEncodedBucket>&key=<urlEncodedKey>&fops=<urlEncodedFops>
	// &notifyURL=<urlEncodedPersistentNotifyUrl>&force=<Force>&pipeline=<Pipeline Name>
	resp := new(interface{})
	err = cdn.bucketManager.Client.CallWithForm(ctx, resp, "POST", "http://api.qiniu.com/pfop/", http.Header{}, form)
	return err
}

func (cdn *CDN) DeleteCDNMedia(key string) error {
	err := cdn.bucketManager.Delete(cdn.Bucket, key)
	return err
}

func (cdn *CDN) MakeCdnUrl(key string) string {
	return storage.MakePublicURLv2(cdn.MakeCdnOrigin(), key)
}

func (cdn *CDN) MakeCdnOrigin() string {
	return cdn.Scheme + "://" + cdn.Domain + "/"
}

func (cdn *CDN) SignRequest(req *http.Request) (token string, err error) {
	return cdn.mac.SignRequest(req)
}

func (cdn *CDN) GetAccessToken(data []byte) string {
	return cdn.mac.Sign(data)
}

func (cdn *CDN) GetAccessTokenWithData(data []byte) string {
	return cdn.mac.SignWithData(data)
}

func (cdn *CDN) Download(key string) (io.ReadCloser, error) {
	url := cdn.MakeCdnUrl(key)

	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
