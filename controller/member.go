package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mergeVersion1/request"
	"mergeVersion1/response"
	"mergeVersion1/types"
	"mergeVersion1/utils"
	"net/http"
	"strconv"
)


func checkParams(params request.CreateMemberRequest) bool{
	if len(params.Nickname) < 4 || len(params.Nickname) > 20 {
		fmt.Println("昵称长度不正确")
		return false
	}
	if len(params.Username) < 8 || len(params.Username) > 20 {
		fmt.Println("用户名长度不正确")
		return false
	}
	//检查username是否只有大小写
	for _,val := range params.Username {
		if (val >= 'a' && val <= 'z') || (val >= 'A' && val <='Z'){
			continue
		}
		fmt.Println("用户名只能包含大小写")
		return false
	}
	if len(params.Password) < 8 || len(params.Password) > 20 {
		fmt.Println("密码长度不正确")
		return false
	}
	if params.UserType != types.Admin && params.UserType != types.Student && params.UserType != types.Teacher {
		fmt.Println("用户类型不正确，并不是管理员学生教师三种类型之一")
		return false
	}
	//检验密码是否同时包括大小写、数字
	lowCase, highCase, digit := false, false, false
	for _, value := range params.Password {
		if value >= '0' && value <= '9' {
			digit = true
		} else if value >= 'a' && value <= 'z' {
			lowCase = true
		} else if value >= 'A' && value <= 'Z' {
			highCase = true
		}
		if digit && lowCase && highCase {
			break
		}
	}
	if !digit || !lowCase || !highCase {
		fmt.Println("用户密码大小写不正确")
		return false
	}

	return true
}

func checkPermission(c *gin.Context) bool {
	cookie, err := c.Cookie("camp-session")
	var member types.TMember
	json.Unmarshal([]byte(cookie), &member)
	if err != nil || cookie == "" || member.UserType!=types.Admin{
		return false
	}
	return true
}

//创建成员
func MemberCreatePost(c *gin.Context) {
	var res response.CreateMemberResponse
	var createMember request.CreateMemberRequest

	if !checkPermission(c){
		res.Code=types.PermDenied
		c.JSON(http.StatusOK,res)
		return
	}

	jsons := utils.GetParams(c, createMember)
	json.Unmarshal(jsons,&createMember)

	if !checkParams(createMember){
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var generateId types.GenerateId
	generateId=types.GenerateId{IsDel: 0,UserType: createMember.UserType,UserName: createMember.Username,Nickname: createMember.Nickname}
	create := utils.Db.Create(&generateId)
	if create.Error!=nil{
		res.Code=types.UserHasExisted
		c.JSON(http.StatusOK,res)
		return
	}
	member:=types.TMember{Nickname:createMember.Nickname ,
		Username: createMember.Username,
		UserID: strconv.FormatUint(generateId.ID,10),
		UserType: createMember.UserType}
	if generateId.UserType ==types.Admin{
		admin:=types.TAdmin{
			TMember: member,
			Password: createMember.Password,
		}
		utils.Db.Create(&admin)
	}else if generateId.UserType == types.Student{
		student:=types.TStudent{
			TMember: member,
			Password: createMember.Password,
			IsRange: 0,
		}
		utils.Db.Create(&student)
	}else if generateId.UserType == types.Teacher{
		teacher:=types.TTeacher{
			TMember: member,
			Password: createMember.Password,
			IsRange: 0,
		}
		utils.Db.Create(&teacher)
	}
	res.Code=types.OK
	res.Data.UserID=strconv.FormatUint(generateId.ID,10)
	c.JSON(http.StatusOK,res)
	return
}

//查询成员单个成员
func MemberGetOne(c *gin.Context)  {
	var params request.GetMemberRequest
	var res response.GetMemberResponse
	var data types.TMember
	jsons := utils.GetParams(c,params)
	json.Unmarshal(jsons,&params)
	id, err := strconv.Atoi(params.UserID)
	if err !=nil || len(params.UserID)<1{
		res.Code=types.ParamInvalid
		return
	}
	generateId:=types.GenerateId{ID: uint64(id)}
	first := utils.Db.First(&generateId)

	if first.Error!=nil{
		res.Code=types.UserNotExisted
		res.Data=data
		c.JSON(http.StatusOK,res)
		return
	}

	if generateId.IsDel == 1{
		res.Code=types.UserHasDeleted
		res.Data=data
		c.JSON(http.StatusOK,res)
		return
	}

	if generateId.UserType == types.Admin{
		var admin types.TAdmin
		utils.Db.Where("user_id = ?",generateId.ID).Find(&admin)
		data=admin.TMember
	}else if generateId.UserType == types.Student{
		var student types.TStudent
		utils.Db.Where("user_id = ?",generateId.ID).Find(&student)
		data=student.TMember
	}else if generateId.UserType == types.Teacher{
		var teacher types.TTeacher
		utils.Db.Where("user_id = ?",generateId.ID).Find(&teacher)
		data=teacher.TMember
	}
	res.Code=types.OK
	res.Data=data
	c.JSON(http.StatusOK,res)
	return
}

//查询许多成员
func MemberGetList(c *gin.Context)  {
	//防止不传值,所以给了默认值
	var requestParams=request.GetMemberListRequest{Limit: -1,Offset: -1}
	jsons := utils.GetParams(c, requestParams)
	json.Unmarshal(jsons, &requestParams)
	limit:=requestParams.Limit
	offset:=requestParams.Offset
	var res response.GetMemberListResponse
	var data []types.TMember
	fmt.Println(limit)
	fmt.Println(offset)
	if limit<0 || offset<0 { //offset和limit小于0会出错
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var list []types.GenerateId
	utils.Db.Where("is_del=?", 0).Limit(limit).Offset(offset).Find(&list)
	for _,v := range list {
		var member types.TMember
		member.UserID=strconv.FormatUint(v.ID,10)
		member.UserType=v.UserType
		member.Username=v.UserName
		member.Nickname=v.Nickname
		data=append(data,member)
	}
	if len(data)>0 {
		res.Code = types.OK
	}else{
		res.Code=types.UserNotExisted
	}
	res.Data.MemberList=data
	c.JSON(http.StatusOK,res)
	return
}

//更新成员

func MemberUpdate(c *gin.Context){
	var paramRequest request.UpdateMemberRequest
	jsons:=utils.GetParams(c,paramRequest)
	json.Unmarshal(jsons,&paramRequest)
	var res response.UpdateMemberResponse
	//update要记住nickname的规定,在这里进行判断
	if _,err:=strconv.Atoi(paramRequest.UserID);  err != nil || len(paramRequest.UserID)<1 || len(paramRequest.Nickname)<4 || len(paramRequest.Nickname)>20{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var generateId types.GenerateId
	first := utils.Db.Where("id = ?", paramRequest.UserID).First(&generateId)
	if first.Error !=nil{
		res.Code=types.UserNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if generateId.IsDel == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res)
		return
	}

	//最好用事务去写
	var err error
	if generateId.UserType == types.Admin{
		err = utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForAdm := utils.Db.Model(&types.TAdmin{}).Where("user_id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
			updateForGen := utils.Db.Model(&types.GenerateId{}).Where("id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
			if updateForAdm.Error!=nil || updateForGen.Error!=nil || updateForGen.RowsAffected<1 || updateForAdm.RowsAffected<1{
				return errors.New("update error")
			}
			return nil
		})
	}else if generateId.UserType == types.Student{
		err = utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForStu :=utils.Db.Model(&types.TStudent{}).Where("user_id = ?",generateId.ID).Update("nick_name",paramRequest.Nickname)
			updateForGen :=utils.Db.Model(&types.GenerateId{}).Where("id = ?",generateId.ID).Update("nick_name",paramRequest.Nickname)
			if updateForGen.Error !=nil || updateForStu.Error !=nil || updateForStu.RowsAffected<1 || updateForGen.RowsAffected<1{
				return errors.New("update error")
			}
			return nil
		})
	}else if generateId.UserType == types.Teacher{
		err = utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForTea := utils.Db.Model(&types.TTeacher{}).Where("user_id = ?",generateId.ID).Update("nick_name",paramRequest.Nickname)
			updateForGen := utils.Db.Model(&types.GenerateId{}).Where("id = ?",generateId.ID).Update("nick_name",paramRequest.Nickname)
			if updateForGen.Error !=nil || updateForTea.Error !=nil || updateForGen.RowsAffected<1 || updateForTea.RowsAffected<1{
				return errors.New("update error")
			}
			return nil
		})
	}
	if err != nil{
		res.Code=types.UnknownError
	}else {
		res.Code = types.OK
	}
	c.JSON(http.StatusOK,res)
	return
}

//删除成员
func MemberDelete(c *gin.Context){
	var paramsRequest request.DeleteMemberRequest
	jsons:=utils.GetParams(c,paramsRequest)
	json.Unmarshal(jsons,&paramsRequest)
	var res response.DeleteMemberResponse
	_, err := strconv.Atoi(paramsRequest.UserID)

	if err!=nil || len(paramsRequest.UserID) < 1 {
		res.Code = types.ParamInvalid
		c.JSON(http.StatusOK, res)
		return
	}
	var generateId types.GenerateId
	first := utils.Db.Where("id = ?", paramsRequest.UserID).First(&generateId)
	if first.Error != nil{
		res.Code=types.UserNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if generateId.IsDel == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res);
		return
	}

	var errForDel error
	if generateId.UserType == types.Admin{
		errForDel =utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForGen := utils.Db.Model(&generateId).Update("is_del",1)
			updateForAdm := utils.Db.Model(&types.TAdmin{}).Where("user_id = ?",generateId.ID).Update("is_del",1)
			if updateForGen.Error != nil || updateForAdm.Error!=nil || updateForGen.RowsAffected<1 || updateForAdm.RowsAffected<1{
				return errors.New("update err")
			}
			return nil
		})
	}else if generateId.UserType == types.Student{
		errForDel = utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForGen := utils.Db.Model(&generateId).Update("is_del",1)
			updateForStu := utils.Db.Model(&types.TStudent{}).Where("user_id = ?",generateId.ID).Update("is_del",1)
			if updateForGen.Error!=nil || updateForStu.Error!=nil || updateForGen.RowsAffected<1 || updateForStu.RowsAffected<1{
				return errors.New("update err")
			}
			return nil
		})
	}else if generateId.UserType == types.Teacher{
		errForDel = utils.Db.Transaction(func(tx *gorm.DB) error {
			updateForGen := utils.Db.Model(&generateId).Update("is_del",1)
			updateForTea := utils.Db.Model(&types.TTeacher{}).Where("user_id = ?",generateId.ID).Update("is_del",1)
			if updateForGen.Error!=nil || updateForTea.Error!=nil || updateForGen.RowsAffected<1 || updateForTea.RowsAffected<1{
				return errors.New("update err")
			}
			return nil
		})
	}

	if errForDel != nil{
		res.Code=types.UnknownError
	}else{
		res.Code=types.OK
	}
	c.JSON(http.StatusOK,res)
}




