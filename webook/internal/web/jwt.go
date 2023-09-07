package web

import "github.com/golang-jwt/jwt/v5"

type UserClaim struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

var JWTKey = []byte("moyn8y9abnd7q4zkq2m73yw8tu9j5ixm")
