package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(dsn string) (*gorm.DB, error) {
	DBConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return DBConn, nil
}
