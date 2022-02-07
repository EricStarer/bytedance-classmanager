package model

import (
	"Register/utils"
)

type User struct {
	Nick_name string `json:"nick_name"`
	Type_user string `json:"type_user"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

//添加新成员
func Create(nick_name, type_user, username, password string) {
	user := User{
		Nick_name: nick_name,
		Type_user: type_user,
		Username:  username,
		Password:  password,
	}
	utils.Db.Create(&user)
}
