package params

type SMSReq struct {
	Mobile    string `json:"mobile" form:"mobile"`
	Captcha   string `json:"captcha" form:"captcha"`     // 验证码
	CaptchaID string `json:"captchaID" form:"captchaID"` // 验证码ID
}

type SMSResp struct {
	Result bool `json:"result"`
	Wait   int  `json:"Wait"`
}

type SysCaptchaResponse struct {
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
}
