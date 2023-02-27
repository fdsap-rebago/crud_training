package goutils

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DBConnect *gorm.DB
	DBErr     error
)

func PostgresConnection(host, user, pass, dbname, port string) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s, sslmode=disable TimeZone=Asia/Manila",
		host, user, pass, dbname, port)

	DBConnect, DBErr = gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func MigrateModel(model interface{}) {
	migratErr := DBConnect.AutoMigrate(model)
	if migratErr != nil {
		fmt.Println("MIGRATE ERR:", migratErr)
	}
}
