package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

func SetColumnWidth(xlsx *excelize.File) {
	xlsx.SetDefaultFont("맑은 고딕")
	setStyleWithFont(xlsx)
	for col := 1; col <= 11; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		switch index := col % 4; index {
		case 1:
			xlsx.SetColWidth("Sheet1", colName, colName, 14.97)
		case 2:
			xlsx.SetColWidth("Sheet1", colName, colName, 8.97)
		case 3:
			xlsx.SetColWidth("Sheet1", colName, colName, 17.97)
		case 0:
			xlsx.SetColWidth("Sheet1", colName, colName, 2.97)
		}
	}
}

func setStyleWithFont(xlsx *excelize.File) {
	font := &excelize.Font{Size: 10}
	alignment := &excelize.Alignment{Horizontal: "center", Vertical: "center"}

	alignment.ShrinkToFit = true
	defaultStyle, _ := xlsx.NewStyle(&excelize.Style{Font: font, Alignment: alignment})
	xlsx.SetColStyle("Sheet1", "A:K", defaultStyle)

	alignment.WrapText = true
	alignment.ShrinkToFit = false
	titleStyle, _ := xlsx.NewStyle(&excelize.Style{Font: font, Alignment: alignment})
	xlsx.SetColStyle("Sheet1", "A", titleStyle)
	xlsx.SetColStyle("Sheet1", "E", titleStyle)
	xlsx.SetColStyle("Sheet1", "I", titleStyle)
}

func PrintTicket(xlsx *excelize.File, titleHCell string, ticket *Ticket) {
	col, row, _ := excelize.CellNameToCoordinates(titleHCell)

	titleVCell, _ := excelize.CoordinatesToCellName(col+2, row+1)
	footerHCell, _ := excelize.CoordinatesToCellName(col, row+8)
	footerVCell, _ := excelize.CoordinatesToCellName(col+2, row+8)
	imageHCell, _ := excelize.CoordinatesToCellName(col, row+2)
	imageVCell, _ := excelize.CoordinatesToCellName(col, row+7)
	examCodeCell, _ := excelize.CoordinatesToCellName(col+1, row+2)
	examCodeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+2)
	nameCell, _ := excelize.CoordinatesToCellName(col+1, row+3)
	nameValueCell, _ := excelize.CoordinatesToCellName(col+2, row+3)
	middleSchoolCell, _ := excelize.CoordinatesToCellName(col+1, row+4)
	middleSchoolValueCell, _ := excelize.CoordinatesToCellName(col+2, row+4)
	isDaejeonCell, _ := excelize.CoordinatesToCellName(col+1, row+5)
	isDaejeonValueCell, _ := excelize.CoordinatesToCellName(col+2, row+5)
	receiptCodeCell, _ := excelize.CoordinatesToCellName(col+1, row+7)
	receiptCodeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+7)
	applyTypeCell, _ := excelize.CoordinatesToCellName(col+1, row+6)
	applyTypeValueCell, _ := excelize.CoordinatesToCellName(col+2, row+6)

	xlsx.MergeCell("Sheet1", titleHCell, titleVCell)
	xlsx.MergeCell("Sheet1", footerHCell, footerVCell)
	xlsx.MergeCell("Sheet1", imageHCell, imageVCell)

	xlsx.SetCellStr("Sheet1", titleHCell, "2021학년도 대덕소프트웨어마이스터고등학교\n입학전형 수험표")
	xlsx.SetCellStr("Sheet1", footerHCell, "대덕소프트웨어마이스터고등학교장")

	xlsx.SetCellStr("Sheet1", examCodeCell, "수험번호")
	xlsx.SetCellStr("Sheet1", nameCell, "성명")
	xlsx.SetCellStr("Sheet1", middleSchoolCell, "출신중학교")
	xlsx.SetCellStr("Sheet1", isDaejeonCell, "지역")
	xlsx.SetCellStr("Sheet1", applyTypeCell, "전형유형")
	xlsx.SetCellStr("Sheet1", receiptCodeCell, "접수번호")

	xlsx.SetCellValue("Sheet1", examCodeValueCell, ticket.ExamCode)
	xlsx.SetCellValue("Sheet1", nameValueCell, ticket.Name)
	xlsx.SetCellValue("Sheet1", middleSchoolValueCell, ticket.MiddleSchool)
	xlsx.SetCellValue("Sheet1", isDaejeonValueCell, ticket.IsDaejeon)
	xlsx.SetCellValue("Sheet1", applyTypeValueCell, ticket.ApplyType)
	xlsx.SetCellValue("Sheet1", receiptCodeValueCell, ticket.ReceiptCode)

	if ticket.ImageURI != "" {
		SaveIfEmpty(ticket.ImageURI)
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

	for i := row; i <= row+9; i++ {
		xlsx.SetRowHeight("Sheet1", i, 18)
	}

	setBorderStyle(xlsx, titleHCell, titleVCell, 2, 2, 2, 1, true)
	setBorderStyle(xlsx, imageHCell, receiptCodeValueCell, 1, 1, 1, 1, false)
	setBorderStyle(xlsx, imageHCell, imageVCell, 2, 1, 1, 1, false)
	setBorderStyle(xlsx, examCodeValueCell, receiptCodeValueCell, 1, 1, 2, 1, false)
	setBorderStyle(xlsx, footerHCell, footerVCell, 2, 1, 2, 2, false)
}

func setBorderStyle(xlsx *excelize.File, hCell string, vCell string, leftStyle int, topStyle int, rightStyle int, bottomStyle int, isWrapText bool) {
	left := excelize.Border{Type: "left", Color: "#000000", Style: leftStyle}
	right := excelize.Border{Type: "right", Color: "#000000", Style: rightStyle}
	top := excelize.Border{Type: "top", Color: "#000000", Style: topStyle}
	bottom := excelize.Border{Type: "bottom", Color: "#000000", Style: bottomStyle}

	font := &excelize.Font{Size: 10}
	alignment := &excelize.Alignment{Horizontal: "center", Vertical: "center"}

	if isWrapText {
		alignment.WrapText = true
	} else {
		alignment.ShrinkToFit = true
	}

	style, _ := xlsx.NewStyle(&excelize.Style{
		Font:      font,
		Alignment: alignment,
		Border:    []excelize.Border{left, right, top, bottom},
	})
	xlsx.SetCellStyle("Sheet1", hCell, vCell, style)
}
