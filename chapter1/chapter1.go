package chapter1

import (
	"fmt"

	"gorm.io/gorm"
)

type Student struct {
	ID    uint
	Name  string
	Age   uint8
	Grade string
}

type Account struct {
	ID      uint
	Balance uint
}

type Transaction struct {
	ID            uint
	FromAccountID uint
	ToAccountID   uint
	Amount        uint
}

func Run(db *gorm.DB) {
	db.AutoMigrate(&Student{})
	db.AutoMigrate(&Account{})
	db.AutoMigrate(&Transaction{})
	var result *gorm.DB

	// 编写SQL语句向 students 表中插入一条新记录，学生姓名为 "张三"，年龄为 20，年级为 "三年级"
	s1 := Student{
		Name:  "张三",
		Age:   20,
		Grade: "三年级",
	}
	result = db.Create(&s1)
	if result.Error == nil {
		fmt.Printf("插入成功，插入信息：%+v\n", s1)
	}

	// 编写SQL语句查询 students 表中所有年龄大于 18 岁的学生信息
	var s2 []Student
	result = db.Where("age > ?", 18).Find(&s2)
	if result.Error == nil {
		fmt.Printf("查询成功，查询结果：%+v\n", s2)
	}

	// 编写SQL语句将 students 表中姓名为 "张三" 的学生年级更新为 "四年级"
	result = db.Model(&Student{}).Where("name = ?", "张三").Update("grade", "四年级")
	if result.Error == nil {
		fmt.Printf("更新成功\n")
	}

	// 编写SQL语句删除 students 表中年龄小于 15 岁的学生记录
	result = db.Where("age < ?", 15).Delete(&Student{})
	if result.Error == nil {
		fmt.Printf("删除成功\n")
	}

	// 事务语句
	var fromID, toID, amount uint = 1, 2, 100
	db.Create(&Account{Balance: amount + 50})
	db.Create(&Account{Balance: amount + 70})

	needRollback := false
	tx := db.Begin()
	t1 := Transaction{
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
	}
	result = tx.Create(&t1)
	if result.Error != nil {
		// 交易记录创建失败
		fmt.Println(result.Error)
		needRollback = true
	}

	a1 := Account{
		ID: fromID,
	}
	result = tx.First(&a1)
	if result.Error != nil {
		// A 账户不存在
		fmt.Println(result.Error)
		needRollback = true
	} else {
		if a1.Balance > amount {
			a1.Balance -= amount
			result = tx.Updates(&a1)
			if result.Error != nil {
				// 更新 A 账户失败
				fmt.Println(result.Error)
				needRollback = true
			} else {
				a2 := Account{
					ID: toID,
				}
				result = tx.First(&a2)
				if result.Error != nil {
					// B 账户不存在
					fmt.Println(result.Error)
					needRollback = true
				} else {
					a2.Balance += amount
					result = tx.Updates(&a2)
					if result.Error != nil {
						// 更新 B 账户失败
						fmt.Println(result.Error)
						needRollback = true
					}
				}
			}
		} else {
			needRollback = true
		}
	}

	if needRollback {
		fmt.Println("事物语句执行失败")
		tx.Rollback()
	} else {
		fmt.Println("事物语句执行成功")
		tx.Commit()
	}
}
