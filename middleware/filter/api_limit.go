package filter

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"pluto/global"
	"sync"

	"golang.org/x/time/rate"
)

// 每60s产生12个令牌
var limiter = NewIPRateLimiter(0.2, 10)

func LimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := ClientIP(c.Request)
		path := c.Request.URL.Path
		global.GVA_LOG.Info(fmt.Sprintf("ip:%s,path:%s", ip, path))

		limit := limiter.GetLimiter(ip + path)
		if !limit.Allow() {
			global.GVA_LOG.Info(fmt.Sprintf("ip:%s,path:%s 被限制访问", ip, path))
			c.Abort()
			return
		}

		c.Next()
	}
}

// IPRateLimiter .
type IPRateLimiter struct {
	ips map[string]*rate.Limiter
	mu  *sync.RWMutex
	r   rate.Limit
	b   int
}

// NewIPRateLimiter .
func NewIPRateLimiter(r rate.Limit, b int) *IPRateLimiter {
	i := &IPRateLimiter{
		ips: make(map[string]*rate.Limiter),
		mu:  &sync.RWMutex{},
		r:   r,
		b:   b,
	}

	return i
}

// AddIP 创建了一个新的速率限制器，并将其添加到 ips 映射中,
// 使用 IP地址作为密钥
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)

	i.ips[ip] = limiter

	return limiter
}

// GetLimiter 返回所提供的IP地址的速率限制器(如果存在的话).
// 否则调用 AddIP 将 IP 地址添加到映射中
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.mu.Unlock()

	return limiter
}
