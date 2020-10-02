package excel

import "github.com/entrydsm/printadmissionticket/db"

type Ticket struct {
	ExamCode     string
	Name         string
	MiddleSchool string
	IsDaejeon    string
	ApplyType    string
	ReceiptCode  int
	ImageURI     string
}

func UserToTicket(user db.UserModel) *Ticket {
	var isDaejeon string
	if user.IsDaejeon {
		isDaejeon = "대전"
	} else {
		isDaejeon = "전국"
	}

	applyType := ""
	switch user.ApplyType {
	case "COMMON":
		applyType = "일반전형"
	case "MEISTER":
		applyType = "마이스터인재전형"
	default:
		applyType = "사회통합전형"
	}

	schoolName := ""
	switch user.GradeType {
	case "GRADUATED":
		schoolName = user.GraduatedSchoolName
	case "UNGRADUATED":
		schoolName = user.UngraduatedSchoolName
	}

	return &Ticket{
		ExamCode:     user.ExamCode,
		Name:         user.Name,
		MiddleSchool: schoolName,
		IsDaejeon:    isDaejeon,
		ApplyType:    applyType,
		ReceiptCode:  user.ReceiptCode,
		ImageURI:     user.ImageURI,
	}
}
