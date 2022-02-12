package controller

import (
	"Register/myRedis"
	"Register/request"
	"Register/response"
	"Register/types"
	"Register/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"strings"
)

//处理抢课功能的BookCourse请求
func SelectCourseBookCourse(c *gin.Context)  {
	var requestParams request.BookCourseRequest
	var res response.BookCourseResponse
	jsons := utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	_,errSid :=strconv.Atoi(requestParams.StudentID)
	cid,errCid :=strconv.Atoi(requestParams.CourseID)
	if len(requestParams.StudentID)<1 || len(requestParams.CourseID)<1 || errCid!=nil || errSid!=nil{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	//已经抢过该课程的学生直接返回 下面业务也已经处理这种情况,此处是针对抢客提速的特殊设计
	result, errSid := myRedis.RedisService.Exists(requestParams.StudentID + "->" + requestParams.CourseID).Result()
	if result == 1{
		res.Code=types.StudentHasCourse
		c.JSON(http.StatusOK,res)
		return
	}
	var student types.TStudent
	firstForStudent := utils.Db.Where("user_id = ?", requestParams.StudentID).First(&student)
	if firstForStudent.Error!=nil{
		res.Code=types.StudentNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if student.IsDel == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res)
		return
	}
	//业务上禁止学生多次选同样课程
	if contains := strings.Contains(student.CourseRecordId, "`"+requestParams.CourseID+"`");contains{
		res.Code=types.StudentHasCourse
		c.JSON(http.StatusOK,res)
		return
	}
	var tCourse types.TCourseCwc
	tCourse.ID=uint64(cid)
	firstForCourse := utils.Db.First(&tCourse)
	if firstForCourse.Error!=nil{
		res.Code=types.CourseNotExisted
		c.JSON(http.StatusOK,res)
		return
	}

	//如果该课程已经被抢完,没必要走下面流程
	if _,ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID); ok{
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
		return
	}
	//这里可以读redis,如果redis记录该门课程为0,则返回课程已选
	_, errSid = myRedis.RedisService.Get(requestParams.CourseID).Result()
	if errSid != nil{
		myRedis.RedisService.Set(strconv.FormatUint(tCourse.ID,10),tCourse.Capacity,myRedis.RedisTimeOut)
	}
	decr := myRedis.RedisService.Decr(requestParams.CourseID)
	if decr.Val()<0{
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
		return
	}

	//此处用事务写
	errUpdate := utils.Db.Transaction(func(tx *gorm.DB) error {
		errForCourse := tx.Model(&tCourse).Where("capacity > ?", 0).UpdateColumn("capacity", gorm.Expr("capacity - ?", 1))
		errForStudent := tx.Model(&student).Update("course_record_id",student.CourseRecordId+"`"+requestParams.CourseID+"`;")
		if errForCourse.Error !=nil || errForStudent.Error!=nil{
			return errForCourse.Error
		}
		return nil
	})

	//更新失败
	if errUpdate !=nil{
		myRedis.RedisService.Incr(requestParams.CourseID)
		myRedis.CourseCapacityMap.Delete(requestParams.CourseID)
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
		return
	}

	if tCourse.Capacity<=0{
		myRedis.RedisService.Incr(requestParams.CourseID)
		myRedis.CourseCapacityMap.Store(requestParams.CourseID,true)
	}
	myRedis.RedisService.Set(requestParams.StudentID+"->"+requestParams.CourseID,1,myRedis.RedisTimeOut)
	res.Code=types.OK
	c.JSON(http.StatusOK,res)
	return
}

//处理抢客功能的course请求
func SelectCourseGetCourse(c *gin.Context)  {
	var requestParams request.GetStudentCourseRequest
	var res response.GetStudentCourseResponse
	jsons :=utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)

	if _,err:=strconv.Atoi(requestParams.StudentID); err != nil || len(requestParams.StudentID)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var student types.TStudent
	first := utils.Db.Where("user_id = ?", requestParams.StudentID).First(&student)

	if first.Error != nil{
		res.Code=types.StudentNotExisted
		c.JSON(http.StatusOK,res)
		return
	}

	if student.IsDel == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res)
		return
	}

	if len(student.CourseRecordId)<1{
		res.Code=types.StudentHasNoCourse
		c.JSON(http.StatusOK,res)
		return
	}
	var data []types.TCourse
	for _,val:= range strings.Split(student.CourseRecordId,";"){
		cid,err:=strconv.ParseUint(strings.Trim(val,"`"),10,64)
		if err !=nil{
			continue
		}
		var tCourse types.TCourseCwc
		tCourse.ID=cid
		utils.Db.First(&tCourse)
		var course =types.TCourse{CourseID: strconv.FormatUint(tCourse.ID,10),TeacherID: tCourse.TeacherID,Name: tCourse.CourseName}
		data=append(data,course)
	}
	res.Data.CourseList=data
	c.JSON(http.StatusOK,res)
	return
}