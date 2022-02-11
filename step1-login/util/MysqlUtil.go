package util

import (
	"bytedance-classmanager/step1-login/types"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var Db *sqlx.DB

//初始化连接数据库
func InitMysql() {

	dsn := "root:bytedancecamp@tcp(180.184.74.238:3306)/byteDanceProject"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		fmt.Println("connect database failed, ", err)
		return
	}
	Db = db
}

//登录查询
func LoginQuery(loginRequest types.LoginRequest) (types.TMember, error) {
	username := loginRequest.Username
	password := loginRequest.Password
	//按顺序查询三个表
	var tables = [...]string{"t_admin", "t_teacher", "t_student"}
	for _, table := range tables {
		member, err := QueryMember(username, password, table)
		if err != nil {
			fmt.Println("query database failed, ", err)
			return types.TMember{}, err
		}
		if member.UserID != "" {
			return member, nil
		}
	}
	return types.TMember{}, nil
}

//按照用户名密码查询表
func QueryMember(username string, password string, tableName string) (types.TMember, error) {
	var member []types.TMember
	err := Db.Select(&member,
		"select user_id, nick_name, user_name, user_type from "+tableName+
			" where user_name = ? and password = ? and is_del = 0",
		username, password)
	if err != nil {
		return types.TMember{}, err
	}
	if member == nil {
		return types.TMember{}, nil
	}
	//fmt.Println(member[0])
	return member[0], nil
}
