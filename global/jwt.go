package global

import (
	"github.com/golang-jwt/jwt/v4"
)

// CustomClaims Custom claims structure
type CustomClaims struct {
	BaseClaims
	BufferTime int64
	jwt.StandardClaims
}

type BaseClaims struct {
	Username string
}
