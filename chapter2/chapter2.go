package chapter2

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Employee struct {
	ID         uint   `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     uint   `db:"salary"`
}

type Book struct {
	ID     uint   `db:"id"`
	Title  string `db:"title"`
	Author string `db:"author"`
	Price  uint   `db:"price"`
}

var schemas = []string{
	`CREATE TABLE IF NOT EXISTS employees (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    department VARCHAR(255) NOT NULL,
	salary BIGINT UNSIGNED NOT NULL
);`,
	`CREATE TABLE IF NOT EXISTS books (
    id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(50) NOT NULL,
	price BIGINT UNSIGNED NOT NULL
);`,
}

func Run(db *sqlx.DB) {
	// 建表
	for _, schema := range schemas {
		_, err := db.Exec(schema)
		if err != nil {
			fmt.Println("Error creating schema:", err)
			return
		}
	}
	var sqlStr string = ""
	var err error
	var rows *sqlx.Rows

	// 测试数据
	// sqlStr = `insert into employees (name, department, salary) values ("张三", "财务部", 10000), ("李四", "测试部", 15000), ("Ben", "技术部", 20000), ("David", "技术部", 18000);`
	// db.Exec(sqlStr)

	// sqlStr = `insert into books (title, author, price) values ("计算机网络", "Jake", 30), ("DirectX 12", "Ben", 60), ("价值投资", "Kail", 100);`
	// db.Exec(sqlStr)

	// 使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中
	sqlStr = `select * from employees where department = ?`
	rows, err = db.Queryx(sqlStr, "技术部")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("技术部员工信息如下：\n")
	for rows.Next() {
		var e Employee
		err = rows.StructScan(&e)

		if err != nil {
			fmt.Printf("技术部员工信息查询出错：%v\n", err)
			continue
		}
		fmt.Printf("%+v\n", e)
	}

	// 使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中
	var es []Employee
	sqlStr = `select * from employees where salary = (select MAX(salary) as largestPrice from employees)`
	err = db.Select(&es, sqlStr)
	if err != nil {
		fmt.Printf("查询工资最高员工信息出错：%v\n", err)
		return
	}
	fmt.Printf("工资最高员工信息如下：\n")
	for _, e := range es {
		fmt.Printf("%+v\n", e)
	}

	// 使用Sqlx执行一个复杂的查询，例如查询价格大于 50 元的书籍，并将结果映射到 Book 结构体切片中，确保类型安全
	var bs []Book
	sqlStr = `select * from books where price > ?`
	err = db.Select(&bs, sqlStr, 50)
	if err != nil {
		fmt.Printf("查询书籍信息出错：%v\n", err)
		return
	}
	fmt.Printf("price > 50 书籍信息如下：\n")
	for _, b := range bs {
		fmt.Printf("%+v\n", b)
	}
}
