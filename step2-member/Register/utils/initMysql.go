package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Db  *gorm.DB
	err error
)

func InitMysql() {
	fmt.Println("初始化连接数据库")
	//db,err := gorm.Open("mysql","root:bgbiao.top@(127.0.0.1:13306)/test_api?charset=utf8&parseTime=True&loc=Local")
	//Db, err = sql.Open("mysql", "root:bytedancecamp@tcp(180.184.74.238:3306)/test")
	Db, err = gorm.Open("mysql", "root:bytedancecamp@tcp(180.184.74.238:3306)/byteDanceProject?charset=utf8&parseTime=True&loc=Local")
	Db.SingularTable(true)
	if err != nil {
		panic(err.Error())
	}
}
