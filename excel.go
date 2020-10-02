package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	_ "image/jpeg"
	_ "image/png"
	"log"
)

const (
	AdmissionTicketSheet        = "수험표"
	RegionalClassificationSheet = "지원자 지역구분통계표"
)

func InitSheet(xlsx *excelize.File) {
	xlsx.NewSheet(AdmissionTicketSheet)
	xlsx.NewSheet(RegionalClassificationSheet)
	xlsx.DeleteSheet("Sheet1")
}

func SetPageLayout(xlsx *excelize.File) {
	xlsx.SetPageLayout(
		AdmissionTicketSheet,
		excelize.PageLayoutOrientation(excelize.OrientationLandscape),
		excelize.PageLayoutPaperSize(9),
	)

	xlsx.SetPageMargins(AdmissionTicketSheet,
		excelize.PageMarginHeader(0.3),
		excelize.PageMarginFooter(0.3),
		excelize.PageMarginTop(0.25),
		excelize.PageMarginBottom(0.25),
		excelize.PageMarginLeft(0.25),
		excelize.PageMarginRight(0.25),
	)
}

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

	values := map[string]interface{}{
		titleHCell:            "2021학년도 대덕소프트웨어마이스터고등학교\n입학전형 수험표",
		footerHCell:           "대덕소프트웨어마이스터고등학교장",
		examCodeCell:          "수험번호",
		nameCell:              "성명",
		middleSchoolCell:      "출신중학교",
		isDaejeonCell:         "지역",
		applyTypeCell:         "전형유형",
		receiptCodeCell:       "접수번호",
		examCodeValueCell:     ticket.ExamCode,
		nameValueCell:         ticket.Name,
		middleSchoolValueCell: ticket.MiddleSchool,
		isDaejeonValueCell:    ticket.IsDaejeon,
		applyTypeValueCell:    ticket.ApplyType,
		receiptCodeValueCell:  ticket.ReceiptCode,
	}

	for axis, value := range values {
		xlsx.SetCellValue(AdmissionTicketSheet, axis, value)
	}

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
