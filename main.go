package main

import (
	"log"
	"net/http"
	"os"

	"github.com/entrydsm/printadmissionticket/db"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"

	"github.com/entrydsm/printadmissionticket/handler"
)

// todo: cache directory
func main() {
	dsn := os.Getenv("MYSQL_URL")
	dbCon, err := db.InitDB(dsn)
	if err != nil {
		log.Panic("failed to connect db")
	}

	s3Downloader, err := InitS3Downloader()
	if err != nil {
		log.Panic("failed to connect s3")
	}

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Panic("failed to get jwt secret key from env")
	}

	r := router.New()
	r.GET("/api/v5/admin/excel/admission_ticket", func(ctx *fasthttp.RequestCtx) {
		if !IsValidToken(ctx, []byte(jwtSecretKey)) {
			ctx.Error("invalid token", http.StatusUnauthorized)
			return
		}
		if err := handler.PrintApplicantAdmission(ctx, dbCon, s3Downloader); err != nil {
			ctx.Error(err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(fasthttp.ListenAndServe(":8080", r.Handler))
}
