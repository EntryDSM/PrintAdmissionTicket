package main

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
	"log"
)

var (
	ContentType = "application/octet-stream"
	FileName    = "대덕소프트웨어마이스터고등학교_수험표.xlsx"
)

func main() {
	InitDB()
	InitS3()
	r := router.New()
	r.GET("/api/v5/admin/excel/admission_ticket", printApplicantAdmission)
	log.Fatal(fasthttp.ListenAndServe(":8080", CORS(r.Handler)))
}

func printApplicantAdmission(ctx *fasthttp.RequestCtx) {
	if !IsValidToken(ctx) {
		return
	}

	xlsx := excelize.NewFile()

	InitSheet(xlsx)
	SetColumnWidth(xlsx)
	SetPageLayout(xlsx)

	users := FindAllUserStatus()
	axis := "A1"
	for index := 1; index <= len(users); index++ {
		user := users[index-1]
		PrintTicket(xlsx, axis, user.ToTicket())

		if index == len(users) {
			break
		}

		switch column := index % 3; column {
		case 1:
			fallthrough
		case 2:
			col, row, _ := excelize.CellNameToCoordinates(axis)
			axis, _ = excelize.CoordinatesToCellName(col+4, row)
		case 0:
			_, row, _ := excelize.CellNameToCoordinates(axis)
			if index%9 == 0 {
				row += 1
			}
			axis, _ = excelize.CoordinatesToCellName(1, row+10)
		}
	}

	returnXlsx(ctx, xlsx)
}

func returnXlsx(ctx *fasthttp.RequestCtx, xlsx *excelize.File) {
	ctx.Response.Header.SetContentType(ContentType)
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename="+FileName)
	ctx.Response.Header.Set("Content-Transfer-Encoding", "binary")
	ctx.Response.Header.Set("Expires", "0")

	writer := ctx.Response.BodyWriter()
	xlsx.Write(writer)
}
