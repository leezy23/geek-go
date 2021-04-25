package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"log"
	"time"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	dsn := "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	err = db.Ping()
	if err != nil {
		return err
	}
	return nil
}

type Employee struct {
	id     int
	salary int
}

func (e Employee) String() string {
	return fmt.Sprintf("id: %d, salary: %d", e.id, e.salary)
}

func main() {
	err := initDB() // 调用输出化数据库的函数
	if err != nil {
		log.Fatalf("init db failed,err:%v\n", err)
	}
	e, err := GetEmployeeById(2)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(e)
}

func GetEmployeeById(id int) (Employee, error) {
	row := db.QueryRow("select id, salary from employee where id = ?", id)
	var e Employee
	err := row.Scan(&e.id, &e.salary)
	switch {
	case err == sql.ErrNoRows: // 不处理查不到的错误
		err = nil
	case err != nil:
		err = errors.Wrap(err, "查询失败")
	}
	return e, err
}
