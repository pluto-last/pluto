package captcha

import (
	"io"
	"pluto/global"
	"strings"
	"time"

	"github.com/lifei6671/gocaptcha"
)

func Generate(id string, w io.Writer) error {
	//初始化一个验证码对象
	captchaImage := gocaptcha.New(300, 100, gocaptcha.RandLightColor())

	//画上随机直线
	captchaImage.DrawLine(2)

	//大波 浪线
	//captchaImage.DrawSineLine()

	//画边框
	captchaImage.DrawBorder(gocaptcha.ColorToRGB(0x17A7A7A))

	//画随机噪点
	//captchaImage.DrawNoise(gocaptcha.CaptchaComplexHigh)

	//画随机文字噪点
	//captchaImage.DrawTextNoise(gocaptcha.CaptchaComplexLower)
	//画验证码文字，可以预先保持到Session种或其他储存容器种
	text := gocaptcha.RandText(4)
	captchaImage.DrawText(text)

	captchaImage.DrawHollowLine()
	//将验证码保持到输出流种，可以是文件或HTTP流等
	err := captchaImage.SaveImage(w, gocaptcha.ImageFormatJpeg)
	if err != nil {
		return err
	}
	rkey := "captcha:" + id

	global.GVA_REDIS.SetStringExpire(rkey, text, time.Minute*5)
	return nil
}

func Veryfy(id, text string) (result bool) {
	rkey := "captcha:" + id
	if global.GVA_REDIS.IsExist(rkey) {
		val, _ := global.GVA_REDIS.GetString(rkey)
		if strings.ToLower(text) == strings.ToLower(val) {
			result = true
		}
		global.GVA_REDIS.Delete(rkey)
	}
	return result
}
