package jwt

import (
	"errors"
	"pluto/global"
	"pluto/model/params"
	"pluto/model/reply"

	"github.com/dgrijalva/jwt-go"
	jwtReq "github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("Token is expired")
	TokenNotValidYet = errors.New("Token not active yet")
	TokenMalformed   = errors.New("That's not even a token")
	TokenInvalid     = errors.New("Couldn't handle this token:")
)

func NewJWT() *JWT {
	return &JWT{
		[]byte(global.GVA_CONFIG.JWT.SignKey),
	}
}

// 创建一个token
func (j *JWT) CreateToken(claims params.CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// 解析 token
func (j *JWT) ParseToken(tokenString string) (*params.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &params.CustomClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, TokenMalformed
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*params.CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid

	} else {
		return nil, TokenInvalid

	}

}

// JWTAuth 校验jwt-token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := jwtReq.OAuth2Extractor.ExtractToken(c.Request)
		if token == "" {
			reply.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}
		j := NewJWT()
		// parseToken 解析token包含的信息
		claims, err := j.ParseToken(token)
		if err != nil {
			if err.Error() == TokenExpired.Error() {
				reply.FailWithCode(101, err.Error(), c)
			} else {
				reply.FailWithMessage(err.Error(), c)
			}
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Set("userID", claims.UserID)
		c.Next()
	}
}
