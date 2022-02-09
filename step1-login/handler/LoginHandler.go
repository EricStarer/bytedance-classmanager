package handler

import (
	"bytedance-classmanager/step1-login/types"
	"bytedance-classmanager/step1-login/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginHandler(c *gin.Context) {
	var loginRequest types.LoginRequest
	var loginResponse types.LoginResponse
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		loginResponse.Code = types.ParamInvalid
		c.JSON(http.StatusOK, loginResponse)
		return
	}
	member, err := util.LoginQuery(loginRequest)
	if err != nil {
		fmt.Println("query failed, ", err)
		loginResponse.Code = types.ParamInvalid
		c.JSON(http.StatusOK, loginResponse)
		return
	}
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
		"/", "180.184.74.238",
		false, true)
	c.JSON(http.StatusOK, loginResponse)
}

func LogoutHandler(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	var logoutResponse types.LogoutResponse
	if err != nil {
		logoutResponse.Code = types.UnknownError
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
		"/", "180.184.74.238",
		false, true)
	logoutResponse.Code = types.OK
	c.JSON(http.StatusOK, logoutResponse)
}

func WhoamiHandler(c *gin.Context) {
	cookie, err := c.Cookie("camp-session")
	var whoamiResponse types.WhoAmIResponse
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
