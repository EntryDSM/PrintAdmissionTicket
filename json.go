package main

import (
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

var (
	applicationJSON = "application/json"
)

func DoJSONWrite(ctx *fasthttp.RequestCtx, error ErrorResponse) {
	ctx.Response.SetStatusCode(error.Code)

	start := time.Now()
	body, err := json.Marshal(error)
	if err != nil {
		elapsed := time.Since(start)
		fmt.Errorf("", elapsed, err.Error(), error)
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	ctx.SetContentType(applicationJSON)
	ctx.SetStatusCode(error.Code)
	ctx.Write(body)
}
