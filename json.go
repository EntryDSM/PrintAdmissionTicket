package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
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
		log.Println(fmt.Sprintf("[ERROR] %s %s %s", elapsed, err, error.Reason))
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
	ctx.SetContentType(applicationJSON)
	ctx.SetStatusCode(error.Code)
	ctx.Write(body)
}
