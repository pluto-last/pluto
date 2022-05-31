package filter

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/table"
	"pluto/service"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var operationRecordService = service.ServiceGroupAPP.OperationRecordService

// OperationRecord 保存接口调用记录的中间件
func OperationRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var userID string
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = ioutil.ReadAll(c.Request.Body)
			if err != nil {
				global.GVA_LOG.Error("read body from request error:", zap.Any("err", err))
			} else {
				c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
			}
		}
		if claims, ok := c.Get("claims"); ok {
			waitUse := claims.(*params.CustomClaims)
			userID = waitUse.UserID
		} else {
			userID = c.Request.Header.Get("x-user-id")
		}

		ip := ClientIP(c.Request)
		c.Set("Real-Ip", ip)

		record := table.OperationRecord{
			Ip:     ip,
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
			UserID: userID,
		}
		// 存在某些未知错误 TODO
		//values := c.Request.Header.Values("content-type")
		//if len(values) >0 && strings.Contains(values[0], "boundary") {
		//	record.Body = "file"
		//}
		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		c.Next()

		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Resp = writer.body.String()

		if err := operationRecordService.CreateSysOperationRecord(record); err != nil {
			global.GVA_LOG.Error("create operation record error:", zap.Any("err", err))
		}
	}
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
