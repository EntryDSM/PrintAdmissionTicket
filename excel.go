package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

var (
	AdmissionTicketSheet = "수험표"
)

func SetColumnWidth(xlsx *excelize.File) {
	xlsx.SetDefaultFont("맑은 고딕")
	setStyleWithFont(xlsx)
	for col := 1; col <= 11; col++ {
		colName, _ := excelize.ColumnNumberToName(col)
		switch index := col % 4; index {
		case 1:
			xlsx.SetColWidth(AdmissionTicketSheet, colName, colName, 14.97)
		case 2:
			xlsx.SetColWidth(AdmissionTicketSheet, colName, colName, 8.97)
		case 3:
			xlsx.SetColWidth(AdmissionTicketSheet, colName, colName, 17.97)
		case 0:
			xlsx.SetColWidth(AdmissionTicketSheet, colName, colName, 2.97)
		}
	}
}

func setStyleWithFont(xlsx *excelize.File) {
	font := &excelize.Font{Size: 10}
	alignment := &excelize.Alignment{Horizontal: "center", Vertical: "center"}

	alignment.ShrinkToFit = true
	defaultStyle, _ := xlsx.NewStyle(&excelize.Style{Font: font, Alignment: alignment})
	xlsx.SetColStyle(AdmissionTicketSheet, "A:K", defaultStyle)

	alignment.WrapText = true
	alignment.ShrinkToFit = false
	titleStyle, _ := xlsx.NewStyle(&excelize.Style{Font: font, Alignment: alignment})
	xlsx.SetColStyle(AdmissionTicketSheet, "A", titleStyle)
	xlsx.SetColStyle(AdmissionTicketSheet, "E", titleStyle)
	xlsx.SetColStyle(AdmissionTicketSheet, "I", titleStyle)
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

	xlsx.MergeCell(AdmissionTicketSheet, titleHCell, titleVCell)
	xlsx.MergeCell(AdmissionTicketSheet, footerHCell, footerVCell)
	xlsx.MergeCell(AdmissionTicketSheet, imageHCell, imageVCell)

	xlsx.SetCellStr(AdmissionTicketSheet, titleHCell, "2021학년도 대덕소프트웨어마이스터고등학교\n입학전형 수험표")
	xlsx.SetCellStr(AdmissionTicketSheet, footerHCell, "대덕소프트웨어마이스터고등학교장")

	xlsx.SetCellStr(AdmissionTicketSheet, examCodeCell, "수험번호")
	xlsx.SetCellStr(AdmissionTicketSheet, nameCell, "성명")
	xlsx.SetCellStr(AdmissionTicketSheet, middleSchoolCell, "출신중학교")
	xlsx.SetCellStr(AdmissionTicketSheet, isDaejeonCell, "지역")
	xlsx.SetCellStr(AdmissionTicketSheet, applyTypeCell, "전형유형")
	xlsx.SetCellStr(AdmissionTicketSheet, receiptCodeCell, "접수번호")

	xlsx.SetCellValue(AdmissionTicketSheet, examCodeValueCell, ticket.ExamCode)
	xlsx.SetCellValue(AdmissionTicketSheet, nameValueCell, ticket.Name)
	xlsx.SetCellValue(AdmissionTicketSheet, middleSchoolValueCell, ticket.MiddleSchool)
	xlsx.SetCellValue(AdmissionTicketSheet, isDaejeonValueCell, ticket.IsDaejeon)
	xlsx.SetCellValue(AdmissionTicketSheet, applyTypeValueCell, ticket.ApplyType)
	xlsx.SetCellValue(AdmissionTicketSheet, receiptCodeValueCell, ticket.ReceiptCode)

	if ticket.ImageURI != "" {
		SaveIfEmpty(ticket.ImageURI)
		yScale := 0.358
		if col == 1 {
			yScale = 0.3
		}

		err := xlsx.AddPicture(AdmissionTicketSheet, imageHCell, "./cache/"+ticket.ImageURI, fmt.Sprintf(`{
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
		xlsx.SetRowHeight(AdmissionTicketSheet, i, 18)
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
	xlsx.SetCellStyle(AdmissionTicketSheet, hCell, vCell, style)
}
