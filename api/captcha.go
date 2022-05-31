package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"pluto/global"
	"pluto/middleware/captcha"
	"pluto/middleware/sms"
	"pluto/model/params"
	"pluto/model/reply"
	"pluto/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CaptchaCtl struct {
}

// @Tags Base
// @Summary 生成图形验证码
// @Produce application/json
// @Param captchaID query  string true "验证码ID"
// @Router /base/captcha [get]
func (b *CaptchaCtl) Captcha(c *gin.Context) {
	id := c.Query("captchaID")

	if len(id) < 8 {
		reply.FailWithMessage("captchaID 必须大于8位", c)
		return
	}

	buff := new(bytes.Buffer)
	err := captcha.Generate(id, buff)
	if err != nil {
		global.GVA_LOG.Error("验证码获取失败!", zap.Any("err", err))
		reply.FailWithMessage("验证码获取失败", c)
		return
	}
	c.Header("Content-Type", "image/jpeg")
	buff.WriteTo(c.Writer)
}

// @Tags Base
// @Summary 发送短信验证码
// @Produce application/json
// @Param  body body params.SMSReq true  "电话，验证码，验证码ID"
// @Success 200 {string} string "{"code":0,"data":{},"msg":"操作成功"}"
// @Router /base/sms [post]
func (b *CaptchaCtl) SMS(c *gin.Context) {

	var r params.SMSReq
	_ = c.Bind(&r)
	if err := utils.Verify(r, utils.SMSVerify); err != nil {
		reply.FailWithMessage(err.Error(), c)
		return
	}
	if !captcha.Veryfy(r.CaptchaID, r.Captcha) {
		reply.FailWithMessage("验证码错误", c)
		return
	}

	rkey := "sms_send:" + r.Mobile

	// 60秒内不能重复发送验证码
	if global.GVA_REDIS.IsExist(rkey) {
		lastTime := time.Time{}
		str, _ := global.GVA_REDIS.GetString(rkey)
		json.Unmarshal([]byte(str), &lastTime)
		wait := time.Now().Sub(lastTime)

		reply.FailWithMessage(fmt.Sprintf("请%d秒后重试", 60-int(wait.Seconds())), c)
		return
	}
	timeout := time.Minute
	err := sms.SendCaptcha(r.Mobile)
	if err != nil {
		global.GVA_LOG.Error("短信发送失败!", zap.Any("err", err))
		reply.FailWithMessage("短信发送失败", c)
		return
	}

	timestampBytes, _ := json.Marshal(time.Now())
	global.GVA_REDIS.SetStringExpire(rkey, string(timestampBytes), timeout)

	reply.Ok(c)
}
