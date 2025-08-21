package main

import (
	"fmt"

	"github.com/go_task3/chapter3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("数据库连接成功")

	// chapter1.Run(db)
	chapter3.Run(db)
}
