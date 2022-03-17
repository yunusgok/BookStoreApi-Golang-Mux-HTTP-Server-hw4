package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB(conString string) *gorm.DB {
	// conString = "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(conString), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Cannot connect to database : %s", err.Error()))
	}

	return db
}
