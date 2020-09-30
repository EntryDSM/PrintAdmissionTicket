package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"os"
)

func IsValidToken(ctx *fasthttp.RequestCtx) bool {
	accessToken := string(ctx.Request.Header.Peek(fasthttp.HeaderAuthorization))
	if accessToken == "" {
		DoJSONWrite(ctx, ErrorResponse{Reason: "Invalid Token", Code: fasthttp.StatusUnauthorized})
		return false
	}

	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
		return secretKey, nil
	})

	if err != nil || !token.Valid {
		DoJSONWrite(ctx, ErrorResponse{Reason: "Invalid Token", Code: fasthttp.StatusUnauthorized})
		return false
	}

	return true
}
