package excel

import (
	"fmt"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/s3manager"
)

func PrintTicket(s3Downloader *s3manager.Downloader, xlsx *excelize.File, titleHCell string, ticket *Ticket) {
	col, row, _ := excelize.CellNameToCoordinates(titleHCell)

	examCodeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+2)
	nameValueCell, _ := excelize.CoordinatesToCellName(col+2, row+3)
	middleSchoolValueCell, _ := excelize.CoordinatesToCellName(col+2, row+4)
	isDaejeonValueCell, _ := excelize.CoordinatesToCellName(col+2, row+5)
	receiptCodeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+7)
	applyTypeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+6)
	imageHCell, _ := excelize.CoordinatesToCellName(col, row+2)

	xlsx.SetCellValue("Sheet1", examCodeValueCell, ticket.ExamCode)
	xlsx.SetCellValue("Sheet1", nameValueCell, ticket.Name)
	xlsx.SetCellValue("Sheet1", middleSchoolValueCell, ticket.MiddleSchool)
	xlsx.SetCellValue("Sheet1", isDaejeonValueCell, ticket.IsDaejeon)
	xlsx.SetCellValue("Sheet1", applyTypeValueCell, ticket.ApplyType)
	xlsx.SetCellValue("Sheet1", receiptCodeValueCell, ticket.ReceiptCode)

	if ticket.ImageURI != "" {
		SaveIfEmpty(s3Downloader, ticket.ImageURI)
		yScale := 0.358
		if col == 1 {
			yScale = 0.3
		}

		err := xlsx.AddPicture("Sheet1", imageHCell, "./cache/"+ticket.ImageURI, fmt.Sprintf(`{
			"x_offset": 1, 
			"y_offset": 1, 
			"x_scale": 0.369, 
			"y_scale": %f
		}`, yScale))
		if err != nil {
			log.Println(err)
		}
	}

}

func SaveIfEmpty(s3Downloader *s3manager.Downloader, key string) {
	if key == "" {
		log.Fatal()
	}

	filename := "./cache/" + key
	if exists(filename) {
		return
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
		log.Println(err)
	}
}

func exists(name string) bool {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
