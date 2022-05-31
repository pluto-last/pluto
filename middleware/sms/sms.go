package sms

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"pluto/global"
	"pluto/utils"
	"time"

	"go.uber.org/zap"
)

type smsCaptcha struct {
	CaptchaValue string
	CaptchaLeft  int
}

func (sms *smsCaptcha) Json() []byte {
	str, _ := json.Marshal(sms)
	return str
}

func (sms *smsCaptcha) Unjson(jsondata []uint8) error {
	return json.Unmarshal(jsondata, sms)
}

type smsResponse struct {
	Code   int     `json:"code"`
	Msg    string  `json:"msg"`
	Count  int     `json:"count"`
	Fee    float64 `json:"fee"`
	Unit   string  `json:"unit"`
	Mobile string  `json:"mobile"`
	Sid    int64   `json:"sid"`
}

// SendSMS 短信验证码调用接口
func SendSMS(mobile string, text string, params ...map[string]interface{}) (err error) {
	data := url.Values{}
	data.Set("apikey", global.GVA_CONFIG.SMS.APIKey)
	data.Set("mobile", mobile)
	data.Set("text", text)
	resp, err := http.PostForm("https://sms.yunpian.com/v2/sms/single_send.json", data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	global.GVA_LOG.Info("发送短信验证码。", zap.Any("info", data))
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		resp := new(smsResponse)
		err = json.Unmarshal(body, resp)
		if err != nil {
			return err
		}
		if resp.Code != 0 {
			global.GVA_LOG.Info(fmt.Sprintf("%#v", resp))
			return errors.New(resp.Msg)
		}
	} else {
		err = errors.New(resp.Status)
		if err != nil {
			return err
		}
		return err
	}
	return err
}

func SendCaptcha(mobile string, keys ...[]string) (err error) {
	identifying := utils.RandDigitStr(6)

	//err = SendSMS(mobile, fmt.Sprintf(global.GVA_CONFIG.SMS.Template, identifying))
	//if err != nil {
	//	return
	//}
	global.GVA_LOG.Info("发送短信验证码成功", zap.Any("mobile", mobile), zap.Any("code", identifying))

	rkey := "sms_result:" + mobile
	if global.GVA_REDIS.IsExist(rkey) {
		global.GVA_REDIS.Delete(rkey)
	}
	global.GVA_REDIS.SetStringExpire(rkey, string((&smsCaptcha{identifying, 3}).Json()), time.Minute*5)
	return
}

// VerifyCaptcha 校验短信验证码是否正确
func VerifyCaptcha(mobile, smsCode string) (err error) {

	rkey := "sms_result:" + mobile
	if !global.GVA_REDIS.IsExist(rkey) {
		err = errors.New("请重新发送验证码")
		return
	}

	smscaptcha := new(smsCaptcha)
	str, _ := global.GVA_REDIS.GetString(rkey)
	err = smscaptcha.Unjson([]byte(str))
	if err != nil {
		err = errors.New("验证码错误,请重新发送验证码")
	}
	if smscaptcha.CaptchaValue != smsCode {
		smscaptcha.CaptchaLeft -= 1
		if smscaptcha.CaptchaLeft == 0 {
			global.GVA_REDIS.Delete(rkey)
			err = errors.New("验证码错误,请重新发送验证码")
		} else {
			err = errors.New("验证码错误")
		}
	}
	return
}

// 清除验证码缓存
func ClearCaptcha(mobile string) {
	if global.GVA_REDIS.IsExist(mobile) {
		global.GVA_REDIS.Delete(mobile)
	}
}
