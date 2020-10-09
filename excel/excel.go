package excel

import (
	"errors"
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"

	"github.com/entrydsm/printadmissionticket/db"
)

func WriteAdmissionTicketsToExcel(s3Downloader *s3manager.Downloader, xlsx *excelize.File, users []db.UserModel) (*excelize.File, error) {
	for i, user := range users {
		err := PrintTicket(s3Downloader, xlsx, i, UserToTicket(user))
		if err != nil {
			return nil, err
		}
	}
	return xlsx, nil
}

func PrintTicket(s3Downloader *s3manager.Downloader, xlsx *excelize.File, index int, ticket *Ticket) error {
	imageColumns := []string{"A", "E", "I"}
	infoColumns := []string{"B", "G", "K"}
	colNum, rowNum := index%3, index/3
	row := (rowNum * 10) + 3 + (index+1)/9

	if ticket.ImageURI != "" {
		if err := SaveIfEmpty(s3Downloader, ticket.ImageURI); err != nil {
			return err
		}
		if err := xlsx.AddPicture("Sheet1", fmt.Sprintf("%s%d", imageColumns[colNum], row), "./cache/"+ticket.ImageURI, `{
			"x_offset": 1, 
			"y_offset": 1, 
			"x_scale": 0.369, 
			"y_scale": 0.358
		}`); err != nil {
			return err
		}

		// NOTE: 왜 col == 1 일 때만 yScale이 다른지 확인
		//yScale := 0.358
		//if col == 1 {
		//	yScale = 0.3
		//}
		//
		//err := xlsx.AddPicture("Sheet1", imageHCell, "./cache/"+ticket.ImageURI, fmt.Sprintf(`{
		//	"x_offset": 1,
		//	"y_offset": 1,
		//	"x_scale": 0.369,
		//	"y_scale": %f
		//}`, yScale))
		//if err != nil {
		//	return err
		//}
	}

	col := infoColumns[colNum]
	infos := []string{ticket.ExamCode, ticket.Name, ticket.MiddleSchool, ticket.IsDaejeon, ticket.ApplyType, fmt.Sprintf("%d", ticket.ReceiptCode)}
	for i := 0; i < len(infos); i++ {
		if err := xlsx.SetCellValue("Sheet1", fmt.Sprintf("%s%d", col, row+i), infos[i]); err != nil {
			return err
		}
	}

	return nil
}

func SaveIfEmpty(s3Downloader *s3manager.Downloader, key string) error {
	if key == "" {
		return errors.New("empty key")
	}

	filename := "./cache/" + key
	if exists(filename) {
		return nil
	}

	file, err := os.Create("./cache/" + key)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// todo: set bucket name on env
	_, err = s3Downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String("image.entrydsm.hs.kr"),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
