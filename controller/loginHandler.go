package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"mergeVersion1/request"
	"mergeVersion1/response"
	"mergeVersion1/types"
	"mergeVersion1/utils"
	"net/http"
	"strings"
)

func LoginHandler(c *gin.Context) {
	var loginRequest request.LoginRequest
	var loginResponse response.LoginResponse
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		loginResponse.Code = types.ParamInvalid
		c.JSON(http.StatusOK, loginResponse)
		return
	}
	member := LoginQuery(loginRequest)
	if member.UserID == "" {
		loginResponse.Code = types.WrongPassword
		c.JSON(http.StatusOK, loginResponse)
		return
	}
	loginResponse.Code = types.OK
	loginResponse.Data.UserID = member.UserID
	//将结构体格式化为json字符串
	cookieVal, err := json.Marshal(member)
	if err != nil {
		fmt.Println("Format json failed, ", err)
	}
	//设置cookie
	c.SetCookie("camp-session", string(cookieVal),
		7*24*60*60,
		"/", "",
		false, true)
	c.JSON(http.StatusOK, loginResponse)
}

func LogoutHandler(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	var logoutResponse response.LogoutResponse
	if err != nil {
		logoutResponse.Code = types.LoginRequired
		c.JSON(http.StatusUnauthorized, logoutResponse)
		return
	}
	if cookie == "" {
		logoutResponse.Code = types.LoginRequired
		c.JSON(http.StatusUnauthorized, logoutResponse)
		return
	}
	//清空cookie
	c.SetCookie("camp-session", "",
		-1,
		"/", "",
		false, true)
	logoutResponse.Code = types.OK
	c.JSON(http.StatusOK, logoutResponse)
}

func WhoamiHandler(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	var whoamiResponse response.WhoAmIResponse
	if err != nil {
		whoamiResponse.Code = types.LoginRequired
		c.JSON(http.StatusUnauthorized, whoamiResponse)
		return
	}
	if cookie == "" {
		whoamiResponse.Code = types.LoginRequired
		c.JSON(http.StatusUnauthorized, whoamiResponse)
		return
	}
	var cookieVal types.TMember
	json.Unmarshal([]byte(cookie), &cookieVal)
	//fmt.Println(cookieVal)
	whoamiResponse.Code = types.OK
	whoamiResponse.Data = cookieVal
	c.JSON(http.StatusOK, whoamiResponse)
}

//登录查询
func LoginQuery(loginRequest request.LoginRequest) types.TMember {
	username := loginRequest.Username
	password := loginRequest.Password
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return types.TMember{}
	}
	//按顺序查询三个表
	return QueryMember(username, password)
}

//按照用户名密码查询表
func QueryMember(username string, password string) types.TMember {
	var admin types.TAdmin
	utils.Db.Where("user_name = ? and password = ? and is_del = 0",
		username, password).First(&admin)
	if admin.TMember.UserID != "" {
		return admin.TMember
	}
	var student types.TStudent
	utils.Db.Where("user_name = ? and password = ? and is_del = 0",
		username, password).First(&student)
	if student.TMember.UserID != "" {
		return student.TMember
	}
	var teacher types.TTeacher
	utils.Db.Where("user_name = ? and password = ? and is_del = 0",
		username, password).First(&teacher)
	if teacher.TMember.UserID != "" {
		return teacher.TMember
	}
	return types.TMember{}
}
