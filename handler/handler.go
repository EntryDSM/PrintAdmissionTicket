package handler

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"

	"github.com/entrydsm/printadmissionticket/db"
	"github.com/entrydsm/printadmissionticket/excel"
)

const (
	ContentType = "application/octet-stream"
	FileName    = "대덕소프트웨어마이스터고등학교_수험표.xlsx"
)

func PrintApplicantAdmission(ctx *fasthttp.RequestCtx, dbCon *gorm.DB, downloader *s3manager.Downloader) error {
	xlsx, err := excelize.OpenFile("template.xlsx")
	if err != nil {
		return err
	}

	users := ListSubmittedUsers(dbCon)
	xlsx, err = excel.WriteAdmissionTicketsToExcel(downloader, xlsx, users)
	if err != nil {
		return err
	}

	ctx.Response.Header.SetContentType(ContentType)
	ctx.Response.Header.Set("Content-Disposition", "attachment; filename="+FileName)
	ctx.Response.Header.Set("Content-Transfer-Encoding", "binary")
	ctx.Response.Header.Set("Expires", "0")
	writer := ctx.Response.BodyWriter()
	err = xlsx.Write(writer)
	if err != nil {
		return err
	}

	return nil
}

func ListSubmittedUsers(DBConn *gorm.DB) []db.UserModel {
	var users []db.UserModel
	DBConn.Table("user").
		Select("s.exam_code, user.receipt_code, user.apply_type, user.is_daejeon, user.name, user.grade_type, "+
			"user.user_photo, gas.school_name AS 'graduated_school_name', uas.school_name AS 'ungraduated_school_name'").
		Joins("JOIN status s on user.receipt_code = s.user_receipt_code").
		Joins("LEFT JOIN graduated_application ga on user.receipt_code = ga.user_receipt_code").
		Joins("LEFT JOIN ungraduated_application ua on user.receipt_code = ua.user_receipt_code").
		Joins("LEFT JOIN school gas on ga.school_code = gas.school_code").
		Joins("LEFT JOIN school uas on ua.school_code = uas.school_code").
		Where("s.is_final_submit = ?", 1).
		Find(&users)
	return users
}
