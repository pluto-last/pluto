package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"pluto/global"
	"pluto/model/constant"
	"sort"
	"strconv"
	"strings"
	"time"
)

// OpenApiCheck 校验签名
func OpenApiCheck(c *gin.Context) (res bool, err error) {

	appid, ok := c.GetPostForm("appid")
	if !ok {
		global.GVA_LOG.Error("appid不存在")
		err = errors.New("appid不存在")
		return
	}

	// appid base64解密
	bs4, _ := base64.StdEncoding.DecodeString(appid)
	if string(bs4) != GetConfigValueString(constant.CtAppidConfig, "Pm21eLMjxX4QW5VT") {
		global.GVA_LOG.Error("appid错误")
		err = errors.New("appid错误")
		return
	}

	sign, ok := c.GetPostForm("sn")
	if !ok {
		global.GVA_LOG.Error("sn不存在")
		err = errors.New("sn不存在")
		return
	}

	// 校验请求是否超市
	timestamp, ok := c.GetPostForm("timestamp")
	if !ok {
		global.GVA_LOG.Error("timestamp不存在")
		err = errors.New("timestamp不存在")
		return
	}

	if !CheckTimestamp(timestamp) {
		global.GVA_LOG.Error("请求已经超时")
		err = errors.New("请求已经超时")
		return
	}

	// 校验签名
	err = CheckSn(c, appid, sign)
	if err != nil {
		err = errors.New("验签失败")
		return
	}

	return true, nil
}

// CheckSn 校验签名
func CheckSn(c *gin.Context, appid, sign string) error {

	fmt.Println("appid", appid)

	if sign != getOpenApiSn(c, appid) {
		return errors.New("验签失败")
	}
	return nil
}

//根据参数map和密钥获取签名
func getOpenApiSn(c *gin.Context, appid string) string {

	//取出所有参数名
	var keys []string
	for k, _ := range c.Request.PostForm {
		if k != "sn" {
			keys = append(keys, k)
		}
	}

	//对参数名做升序排列
	sort.Strings(keys)

	// 参数拼接
	strs := ""
	for _, key := range keys {
		v, _ := c.GetPostForm(key)
		strs += key + "-" + v
	}
	strs = "api" + strs + "api"
	strs = strings.Replace(strs, "\n", "", -1)
	fmt.Println("拼装后的签名字符串为[%s]", strs)
	appid = Get16Key(appid)
	aes := AesEncrypt(strs, appid)

	bs4 := base64.StdEncoding.EncodeToString([]byte(aes))
	sn := MD5(bs4)
	return sn
}

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

//校验时间戳
func CheckTimestamp(tm string) bool {
	tmi, err := strconv.ParseInt(tm, 10, 64)
	if err != nil {
		return false
	}
	now := time.Now().UnixNano() / 1e6
	var inv int64
	inv = 120000
	x := now - tmi
	//前后差距
	if x < 0 {
		x = x * -1
	}
	//如果前后大于设定
	if x > inv {
		return false
	}
	return true
}

//获取16位的key,取左边16位，不足的补0
func Get16Key(key string) string {
	str := key + "0000000000000000"
	rs := []rune(str)
	rtn := string(rs[0:16])
	return rtn
}

//MD5 md5加密
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
