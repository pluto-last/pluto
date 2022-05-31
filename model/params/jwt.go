package params

import (
	"github.com/dgrijalva/jwt-go"
)

// Custom claims structure
type CustomClaims struct {
	UserID string
	jwt.StandardClaims
}
