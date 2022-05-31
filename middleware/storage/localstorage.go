package storage

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
)

type LocalStorage struct {
	Domain string
	Scheme string
}

type resData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	RouteID string `json:"route_id"`
}

func (localS *LocalStorage) Upload(key string, reader io.Reader, size int64) error {
	origin := localS.MakeCdnOrigin()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	formFile, err := writer.CreateFormFile("file", key)
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, reader)
	if err != nil {
		return err
	}

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", origin, body)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())

	query := make(url.Values)
	query.Add("key", key)
	req.URL.RawQuery = query.Encode()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	statusCode := resp.StatusCode
	if statusCode != 200 {
		log.Println("Code:", statusCode)
		log.Println("LocalStorage Upload ERROR:", string(content))
		return errors.New(fmt.Sprintf("上传失败，code： %d， 返回：%s", statusCode, string(content)))
	}
	resMsg := &resData{}
	err = json.Unmarshal(content, resMsg)
	if err != nil {
		return err
	}
	if resMsg.Code != 0 {
		return errors.New(fmt.Sprintf("上传失败，返回：%s", resMsg.Message))
	}

	return nil
}

func (localS *LocalStorage) IsExist(mediaid string) (bool, error) {
	res, err := http.Get(localS.MakeCdnUrl(mediaid))
	if err != nil {
		return false, err
	}
	if res.StatusCode != 200 {
		return false, err
	}
	defer res.Body.Close()
	return false, err
}

func (localS *LocalStorage) Getuptoken(key ...string) string {
	return ""
}

func (localS *LocalStorage) FileStat(key string) (*Stat, error) {
	return nil, nil
}

func (localS *LocalStorage) FetchImage(key, url string) error {
	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	// 上传
	return localS.Upload(key, reader, int64(reader.Size()))
}

func (localS *LocalStorage) FetchWeChatMedia(key, accesstoken, serverId string) error {
	mediaurl := fmt.Sprintf("http://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s", accesstoken, serverId)

	// 下载url
	resp, err := http.Get(mediaurl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	// 上传
	return localS.Upload(key, reader, int64(reader.Size()))
}

func (localS *LocalStorage) Pfop(key, fops, notifyurl string) (err error) {
	return nil
}

func (localS *LocalStorage) DeleteCDNMedia(key string) error {
	return nil
}

func (localS *LocalStorage) MakeCdnUrl(key string) string {
	url := "%s://%s/%s"
	return fmt.Sprintf(url, localS.Scheme, localS.Domain, key)
}

func (localS *LocalStorage) MakeCdnOrigin() string {
	url := "%s://%s/"
	return fmt.Sprintf(url, localS.Scheme, localS.Domain)
}

func (localS *LocalStorage) SignRequest(req *http.Request) (token string, err error) {
	return "", nil
}

func (localS *LocalStorage) GetAccessToken(data []byte) string {
	return ""
}

func (localS *LocalStorage) GetAccessTokenWithData(data []byte) string {
	return ""
}

func (localS *LocalStorage) Download(key string) (io.ReadCloser, error) {
	url := localS.MakeCdnUrl(key)
	// 下载url
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
