package db

type UserModel struct {
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
