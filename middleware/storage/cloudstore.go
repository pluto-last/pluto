package storage

import (
	"io"
	"net/http"
)

type CloudStore interface {
	Getuptoken(key ...string) string
	Upload(mediaid string, reader io.Reader, size int64) error
	Download(mediaid string) (io.ReadCloser, error)
	IsExist(mediaid string) (bool, error)
	DeleteCDNMedia(key string) error
	MakeCdnUrl(key string) string
	MakeCdnOrigin() string
	FileStat(mediaid string) (*Stat, error)
	FetchImage(key, url string) error
	FetchWeChatMedia(key, accesstoken, serverId string) error
	Pfop(key, fops, notifyurl string) (err error)
	SignRequest(req *http.Request) (token string, err error)
	GetAccessToken(data []byte) string
	GetAccessTokenWithData(data []byte) string
}

type Stat struct {
	Hash     string `json:"hash"`
	Fsize    int64  `json:"fsize"`
	PutTime  int64  `json:"putTime"`
	MimeType string `json:"mimeType"`
	EndUser  string `json:"endUser"`
}
