package main

type Ticket struct {
	ExamCode     string
	Name         string
	MiddleSchool string
	IsDaejeon    string
	ApplyType    string
	ReceiptCode  int
	ImageURI     string
}

type User struct {
	ExamCode              string
	ReceiptCode           int
	ApplyType             string
	IsDaejeon             bool `gorm:"column:is_daejeon"`
	Name                  string
	GradeType             string
	ImageURI              string `gorm:"column:user_photo"`
	GraduatedSchoolName   string `gorm:"column:graduated_school_name"`
	UngraduatedSchoolName string `gorm:"column:ungraduated_school_name"`
}

func (user *User) ToTicket() *Ticket {
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
