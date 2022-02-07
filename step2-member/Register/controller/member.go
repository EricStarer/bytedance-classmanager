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
	"strconv"
)

////位于注册首页
//func Member(c *gin.Context) {
//
//	var test types.Test
//	jsons := utils.GetParams(c, test)
//	if jsons == nil {
//		return
//	}
//	json.Unmarshal(jsons, &test)
//	fmt.Println("输出相应参数：", test.UserType, test.Username, test.Password, test.NickName)
//	//model.Create(params["nick_name"].(string), params["type_user"].(string), params["username"].(string), params["password"].(string))
//}

/*
- 用户昵称，必填，不小于 4 位，不超过 20 位（字节）
- 用户名，必填，支持大小写，不小于 8 位 不超过 20 位（字节）
- 密码，必填，同时包括大小写、数字，不少于 8 位 不超过 20 位（字节）
- 用户类型，必填，枚举值
  - 1：管理员
  - 2：学生
  - 3：教师
- 参数不合法返回ParamInvalid状态码
*/
func checkParams(params request.CreateMemberRequest) bool {
	if len(params.Nickname) < 4 || len(params.Nickname) > 20 {
		fmt.Println("昵称长度不正确")
		return false
	}
	if len(params.Username) < 8 || len(params.Username) > 20 {
		fmt.Println("用户名长度不正确")
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

/*
  只有管理员才有操作权限，无权限返回PermDenied 状态码
*/
func checkPermission(c *gin.Context) bool {
	val, _ := c.Cookie("camp-session")
	/*
		用户类型:
		  - 1：管理员
		  - 2：学生
		  - 3：教师
	*/
	if val == "1" {
		return true
	}
	return false
}

//创建成员
func MemberCreatePost(c *gin.Context) {
	//undo 对传入的数据进行检查 以及权限检查
	var createMember request.CreateMemberRequest
	jsons := utils.GetParams(c, createMember)
	fmt.Println(string(jsons))
	json.Unmarshal(jsons, &createMember)
	var res response.CreateMemberResponse
	//权限检查
	if !checkPermission(c) {
		res.Code = types.PermDenied
		res.Data.UserID = "只有管理员才有权限"
		c.JSON(http.StatusOK, res)
		return
	}

	// 数据检查
	if !checkParams(createMember) {
		res.Code = types.ParamInvalid
		res.Data.UserID = "账号注册出错"
		c.JSON(http.StatusOK, res)
		return
	}

	fmt.Println(createMember.Password, createMember.UserType, createMember.Username, createMember.Nickname)
	var generateId types.GenerateId
	generateId = types.GenerateId{IsDel: 0, UserType: createMember.UserType}
	result := utils.Db.Create(&generateId)
	member := types.TMember{Nickname: createMember.Nickname,
		Username: createMember.Username,
		UserID:   strconv.FormatUint(generateId.ID, 10),
		UserType: createMember.UserType}
	if generateId.UserType == types.Admin {
		admin := types.TAdmin{
			TMember:  member,
			Password: createMember.Password,
		}
		utils.Db.Create(&admin)
	} else if generateId.UserType == types.Student {
		student := types.TStudent{
			TMember:  member,
			Password: createMember.Password,
		}
		utils.Db.Create(&student)
	} else if generateId.UserType == types.Teacher {
		teacher := types.TTeacher{
			TMember:  member,
			Password: createMember.Password,
		}
		utils.Db.Create(&teacher)
	}
	fmt.Println(result)
}

//查询成员单个成员
func MemberGetOne(c *gin.Context) {
	var params request.GetMemberRequest
	jsons := utils.GetParams(c, params)
	if jsons == nil {
		return
	}
	json.Unmarshal(jsons, &params)
	id, err := strconv.Atoi(params.UserID)
	if err != nil {
		return
	}
	generateId := types.GenerateId{ID: uint64(id)}
	first := utils.Db.First(&generateId)
	var res response.GetMemberResponse
	var data types.TMember
	if first.Error != nil {
		res.Code = types.UserNotExisted
		res.Data = data
		c.JSON(http.StatusOK, res)
		return
	}
	if generateId.IsDel == 1 {
		res.Code = types.UserHasDeleted
		res.Data = data
		c.JSON(http.StatusOK, res)
		return
	}
	if generateId.UserType == types.Admin {
		var admin types.TAdmin
		utils.Db.Where("user_id = ?", generateId.ID).Find(&admin)
		data = admin.TMember
	} else if generateId.UserType == types.Student {
		var student types.TStudent
		utils.Db.Where("user_id = ?", generateId.ID).Find(&student)
		data = student.TMember
	} else if generateId.UserType == types.Teacher {
		var teacher types.TTeacher
		utils.Db.Where("user_id = ?", generateId.ID).Find(&teacher)
		data = teacher.TMember
	}
	res.Code = types.UserHasExisted
	res.Data = data
	c.JSON(http.StatusOK, res)
	return
}

//查询许多成员
func MemberGetList(c *gin.Context) {
	var requestParams request.GetMemberListRequest
	jsons := utils.GetParams(c, requestParams)
	json.Unmarshal(jsons, &requestParams)
	limit := requestParams.Limit
	offset := requestParams.Offset
	var res response.GetMemberListResponse
	var data []types.TMember
	if limit < 0 || offset < 0 { //offset和limit小于0会出错
		res.Code = types.ParamInvalid
		c.JSON(http.StatusOK, res)
		return
	}
	var list []types.GenerateId
	utils.Db.Where("is_del=?", 0).Limit(limit).Offset(offset).Find(&list)
	for _, v := range list {
		var member types.TMember
		if v.UserType == types.Admin {
			var admin types.TAdmin
			first := utils.Db.Where("user_id = ?", v.ID).First(&admin)
			if first.Error != nil {
				continue
			}
			member = admin.TMember
		} else if v.UserType == types.Student {
			var student types.TStudent
			first := utils.Db.Where("user_id = ?", v.ID).First(&student)
			if first.Error != nil {
				continue
			}
			member = student.TMember
		} else if v.UserType == types.Teacher {
			var teacher types.TTeacher
			first := utils.Db.Where("user_id = ?", v.ID).First(&teacher)
			if first.Error != nil {
				continue
			}
			member = teacher.TMember
		}
		data = append(data, member)
	}
	if len(data) > 0 {
		res.Code = types.OK
	} else {
		res.Code = types.UserNotExisted
	}
	res.Data.MemberList = data
	c.JSON(http.StatusOK, res)
	return
}

//更新成员

func MemberUpdate(c *gin.Context) {
	var paramRequest request.UpdateMemberRequest
	jsons := utils.GetParams(c, paramRequest)
	json.Unmarshal(jsons, &paramRequest)
	var res response.UpdateMemberResponse
	//update要记住nickname的规定,在这里进行判断
	if id, err := strconv.Atoi(paramRequest.UserID); id < 0 || err != nil || len(paramRequest.UserID) < 1 || len(paramRequest.Nickname) < 1 {
		res.Code = types.ParamInvalid
		c.JSON(http.StatusOK, res)
		return
	}
	var generateId types.GenerateId
	first := utils.Db.Where("id = ?", paramRequest.UserID).First(&generateId)
	if first.Error != nil {
		res.Code = types.UserNotExisted
		c.JSON(http.StatusOK, res)
		return
	}
	if generateId.IsDel == 1 {
		res.Code = types.UserHasDeleted
		c.JSON(http.StatusOK, res)
		return
	}
	res.Code = types.OK
	//最好用事务去写
	if generateId.UserType == types.Admin {
		utils.Db.Model(&types.TAdmin{}).Where("user_id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
		utils.Db.Model(&types.GenerateId{}).Where("id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
	} else if generateId.UserType == types.Student {
		utils.Db.Model(&types.TStudent{}).Where("user_id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
		utils.Db.Where(&types.GenerateId{}).Where("id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
	} else if generateId.UserType == types.Teacher {
		utils.Db.Model(&types.TTeacher{}).Where("user_id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
		utils.Db.Model(&types.GenerateId{}).Where("id = ?", generateId.ID).Update("nick_name", paramRequest.Nickname)
	}
	c.JSON(http.StatusOK, res)
	return
}

//删除成员
func MemberDelete(c *gin.Context) {
	var paramsRequest request.DeleteMemberRequest
	jsons := utils.GetParams(c, paramsRequest)
	json.Unmarshal(jsons, &paramsRequest)
	fmt.Println(paramsRequest.UserID)
	var res response.DeleteMemberResponse
	_, err := strconv.Atoi(paramsRequest.UserID)

	if err != nil || len(paramsRequest.UserID) < 1 {
		res.Code = types.ParamInvalid
		c.JSON(http.StatusOK, res)
		return
	}
	var generateId types.GenerateId
	first := utils.Db.Where("id = ?", paramsRequest.UserID).First(&generateId)
	if first.Error != nil {
		res.Code = types.UserNotExisted
		c.JSON(http.StatusOK, res)
		return
	}
	if generateId.IsDel == 1 {
		res.Code = types.UserHasDeleted
		c.JSON(http.StatusOK, res)
		return
	}
	//待改进,用事务来处理
	if generateId.UserType == types.Admin {
		var admin types.TAdmin
		utils.Db.Model(&generateId).Update("is_del", 1)
		utils.Db.Where("user_id = ?", generateId.ID).First(&admin)
		utils.Db.Model(&admin).Update("is_del", 1)
	} else if generateId.UserType == types.Student {
		var student types.TStudent
		utils.Db.Model(&generateId).Update("is_del", 1)
		utils.Db.Where("user_id = ?", generateId.ID).First(&student)
		utils.Db.Model(&student).Update("is_del", 1)
	} else if generateId.UserType == types.Teacher {
		var teacher types.TTeacher
		utils.Db.Model(&generateId).Update("is_del", 1)
		utils.Db.Where("user_id = ?", generateId.ID).First(&teacher)
		utils.Db.Model(&teacher).Update("is_del", 1)
	}
	res.Code = types.OK
	c.JSON(http.StatusOK, res)
}

func MemberCreateGet(c *gin.Context) {

}
