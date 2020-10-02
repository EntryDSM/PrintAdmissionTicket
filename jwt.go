package main

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
)

func IsValidToken(ctx *fasthttp.RequestCtx, jwtSecretKey []byte) bool {
	authHeader := string(ctx.Request.Header.Peek(fasthttp.HeaderAuthorization))
	if authHeader == "" {
		return false
	}

	token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return jwtSecretKey, nil
	})

	if err != nil || !token.Valid {
		return false
	}

	return true
}
