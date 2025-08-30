package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go_task3/chapter2"
	"github.com/jmoiron/sqlx"
)

func main() {
	// gorm链接
	// dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("数据库连接成功")

	// chapter1.Run(db)
	// chapter3.Run(db)

	// sqlx链接
	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	fmt.Println("数据库连接成功")

	chapter2.Run(db)
}
