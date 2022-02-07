package controller

import (
	"Register/request"
	"Register/response"
	"Register/types"
	"Register/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	登陆
	需要添加cookies
	如果是管理员登陆则允许后续注册操作
	SetCookie(key, value string, maxAge int, path, domain string, secure, httpOnly bool)
	key = "login" cookie 名称
	value = "1" or "2" or "3" 1代表管理员 2代表学生 3代表教师 只有管理员才有后续用户crud操作
- 用户输入账号和密码后点击登录。
- 用户名或者密码错误均返回密码错误。
- 登录成功后需要设置 Cookie，Cookie 名称为 camp-session。
*/
func Login(c *gin.Context) {
	var params request.LoginRequest
	jsons := utils.GetParams(c, params)
	json.Unmarshal(jsons, &params)
	
	fmt.Println(params.Username, params.Password)
	var res response.LoginResponse
	var admin types.TAdmin
	var student types.TStudent
	var teacher types.TTeacher

	//依次检查三张数据表，管理员，学生，教师， 如果登陆成功则添加cookies
	user1 := utils.Db.Where("user_name = ? and password = ?", params.Username, params.Password).Find(&admin)
	fmt.Println(user1)
	if user1.Error == nil && admin.IsDel == 0 {
		cookie, err := c.Cookie("camp-session")
		_ = err
		c.SetCookie("camp-session", "1", 3600, "/", "localhost", false, true)
		c.SetCookie("username", params.Username, 3600, "/", "localhost", false, true)
		c.SetCookie("password", params.Password, 3600, "/", "localhost", false, true)
		fmt.Println(cookie)
		res.Data.UserID = "管理员登陆"
		c.JSON(http.StatusOK, res)
		return
	} else if admin.IsDel == 1 {
		res.Code = types.UserHasDeleted
		res.Data.UserID = "管理员用户不存在, 已经被删除"
		c.JSON(http.StatusOK, res)
		return
	}

	user2 := utils.Db.Where("user_name = ? and password = ?", params.Username, params.Password).Find(&student)
	if user2.Error == nil && student.IsDel == 0 {
		cookie, err := c.Cookie("camp-session")
		_ = err
		c.SetCookie("camp-session", "2", 3600, "/", "localhost", false, true)
		c.SetCookie("username", params.Username, 3600, "/", "localhost", false, true)
		c.SetCookie("password", params.Password, 3600, "/", "localhost", false, true)
		fmt.Println(cookie)
		res.Data.UserID = "学生登陆"
		c.JSON(http.StatusOK, res)
		return
	} else if student.IsDel == 1 {
		res.Code = types.UserHasDeleted
		res.Data.UserID = "学生用户不存在"
		c.JSON(http.StatusOK, res)
		return
	}

	user3 := utils.Db.Where("user_name = ? and password = ?", params.Username, params.Password).Find(&teacher)
	if user3.Error == nil && teacher.IsDel == 0 {
		cookie, err := c.Cookie("camp-session")
		_ = err
		c.SetCookie("camp-session", "3", 3600, "/", "localhost", false, true)
		c.SetCookie("username", params.Username, 3600, "/", "localhost", false, true)
		c.SetCookie("password", params.Password, 3600, "/", "localhost", false, true)
		fmt.Println(cookie)
		res.Data.UserID = "教师登陆"
		c.JSON(http.StatusOK, res)
		return
	} else if teacher.IsDel == 1 {
		res.Code = types.UserHasDeleted
		res.Data.UserID = "教师用户不存在"
		c.JSON(http.StatusOK, res)
		return
	}
	res.Code = types.WrongPassword
	res.Data.UserID = "密码错误"
	c.JSON(http.StatusOK, res)
	return

}

/*
   登出时清除cookie
*/
func Logout(c *gin.Context) {
	c.SetCookie("camp-session", "1", -1, "/", "localhost", false, true)
	c.SetCookie("username", "username", -1, "/", "localhost", false, true)
	c.SetCookie("password", "password", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{
		"清除cookie": "camp-session",
	})
}

/*
	- 登录后访问个人信息页可以查看自己的信息，包括用户ID、用户名称、用户昵称。
*/
func Whoami(c *gin.Context) {
	val, _ := c.Cookie("camp-session")
	username, _ := c.Cookie("username")
	password, _ := c.Cookie("password")

	var res response.WhoAmIResponse
	var admin types.TAdmin
	var student types.TStudent
	var teacher types.TTeacher
	if val == "1" {
		utils.Db.Where("user_name = ? and password = ?", username, password).Find(&admin)
		res.Code = types.OK
		res.Data = admin.TMember
	} else if val == "2" {
		utils.Db.Where("user_name = ? and password = ?", username, password).Find(&student)
		res.Code = types.OK
		res.Data = admin.TMember
	} else if val == "3" {
		utils.Db.Where("user_name = ? and password = ?", username, password).Find(&teacher)
		res.Code = types.OK
		res.Data = admin.TMember
	}
	c.JSON(http.StatusOK, res)
}
