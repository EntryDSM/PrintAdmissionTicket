package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

var (
	DBConn *gorm.DB
)

func InitDB() {
	dsn := os.Getenv("MYSQL_URL")
	if dsn == "" {
		log.Panicln("[ERROR] Failed to configure a DataSource: 'MYSQL_URL' attribute is not specified.")
	}

	var err error
	DBConn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Panicln("[ERROR] Failed to Connect Database")
	}
	log.Println("[INFO] Connection Opened to Database")
}

func FindAllUserStatus() []User {
	var user []User
	DBConn.Table("user").
		Select("s.exam_code, user.receipt_code, user.apply_type, user.is_daejeon, user.name, user.grade_type, "+
			"user.user_photo, gas.school_name AS 'graduated_school_name', uas.school_name AS 'ungraduated_school_name'").
		Joins("JOIN status s on user.receipt_code = s.user_receipt_code").
		Joins("LEFT JOIN graduated_application ga on user.receipt_code = ga.user_receipt_code").
		Joins("LEFT JOIN ungraduated_application ua on user.receipt_code = ua.user_receipt_code").
		Joins("LEFT JOIN school gas on ga.school_code = gas.school_code").
		Joins("LEFT JOIN school uas on ua.school_code = uas.school_code").
		Where("s.is_final_submit = ?", 1).
		Find(&user)
	return user
}
