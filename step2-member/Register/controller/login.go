package controller

import (
	"Register/request"
	"Register/response"
	"Register/types"
	"Register/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func MemberLogIn(c *gin.Context)  {

	var paramsRequest request.LoginRequest
	jsons:=utils.GetParams(c,paramsRequest)
	json.Unmarshal(jsons,&paramsRequest)
	var res response.LoginResponse
	if len(paramsRequest.Password)<1 || len(paramsRequest.Username)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var generateId types.GenerateId
	first := utils.Db.Where("user_name = ?", paramsRequest.Username).First(&generateId)
	if first.Error != nil{
		res.Code=types.UserNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if generateId.IsDel == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res)
		return
	}
	var password string
	if generateId.UserType == types.Admin{
		var admin types.TAdmin
		utils.Db.Where("user_id = ?",generateId.ID).First(&admin)
		password=admin.Password
	}else if generateId.UserType == types.Student{
		var student types.TStudent
		utils.Db.Where("user_id = ?",generateId.ID).First(&student)
		password=student.Password
	}else if generateId.UserType == types.Teacher{
		var teacher types.TTeacher
		utils.Db.Where("user_id = ?",generateId.ID).First(&teacher)
		password=teacher.Password
	}
	if password !=paramsRequest.Password{
		res.Code=types.WrongPassword
		c.JSON(http.StatusOK,res)
		return
	}
	res.Code=types.OK
	res.Data.UserID=strconv.FormatUint(generateId.ID,10)

	//后面是session
	session :=sessions.Default(c)
	option := sessions.Options{MaxAge: utils.SessionAgeForLive,Path: utils.SessionPath,Domain: utils.SessionDomain}
	session.Options(option)
	session.Set("user_id",res.Data.UserID)
	session.Set("user_type",generateId.UserType)
	err := session.Save()
	if err!=nil{
		fmt.Println("save 出错了")
	}
	c.JSON(http.StatusOK,res)
	return
}

//登出模块

func MemberLogOut(c *gin.Context)  {
	session := sessions.Default(c)
	var res response.LogoutResponse
	userId := session.Get("user_id")
	userType := session.Get("user_type")
	if userId == nil || userType == nil{
		res.Code=types.LoginRequired
		c.JSON(http.StatusOK,res)
		return
	}
	option := sessions.Options{Path: utils.SessionPath,MaxAge: utils.SessionAgeForDelete,Domain: utils.SessionDomain}
	session.Options(option)
	res.Code=types.OK
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK,res)
	return
}


//Who am i
func WhoAmI(c *gin.Context)  {
	var res response.WhoAmIResponse
	session := sessions.Default(c)
	id := session.Get("user_id")
	userType := session.Get("user_type")
	if id == nil || userType ==nil{
		res.Code=types.LoginRequired
		c.JSON(http.StatusOK,res)
		return
	}
	var member types.TMember
	if userType == types.Admin{
		var admin types.TAdmin
		utils.Db.Where("user_id = ?",id).First(&admin)
		member=admin.TMember
	}else if userType == types.Student{
		var student types.TStudent
		utils.Db.Where("user_id = ?",id).First(&student)
		member=student.TMember
	}else if userType == types.Teacher{
		var teacher types.TTeacher
		utils.Db.Where("user_id = ?",id).First(&teacher)
		member=teacher.TMember
	}
	res.Code=types.OK
	res.Data=member
	c.JSON(http.StatusOK,res)
	return
}