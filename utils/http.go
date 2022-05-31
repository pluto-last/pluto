package utils

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"golang.org/x/xerrors"
	"io/ioutil"
	"net/http"
	"pluto/global"
)

func DoPost(client *http.Client, url string, req interface{}) (respBytes []byte, err error) {
	reqBytes, err := json.Marshal(&req)
	if err != nil {
		global.GVA_LOG.Error("", zap.Any("err", err))
		return
	}

	postReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(reqBytes))
	if err != nil {
		global.GVA_LOG.Error("", zap.Any("err", err))
		return
	}

	postReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(postReq)
	if err != nil {
		global.GVA_LOG.Error("", zap.Any("err", err))
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, xerrors.New(resp.Status)
	}

	respBytes, err = ioutil.ReadAll(resp.Body)
	return
}
