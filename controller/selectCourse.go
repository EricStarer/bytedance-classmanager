package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"mergeVersion1/myRedis"
	"mergeVersion1/request"
	"mergeVersion1/response"
	"mergeVersion1/types"
	"mergeVersion1/utils"
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

	//step1 顶级判断如果该课程已经被抢完,没必要走下面流程
	if _,ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID); ok{
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
		return
	}

	//step2 判断学生是否已抢到该课程
	hasSuccess, _ := myRedis.RedisService.Exists(requestParams.StudentID + "->" + requestParams.CourseID).Result()
	if hasSuccess == 1{
		res.Code=types.StudentHasCourse
		c.JSON(http.StatusOK,res)
		return
	}

	//step3 限制请求发起过于频繁,已经抢过该课程的学生直接返回 下面业务也已经处理这种情况,此处是针对抢客提速的特殊设计
	hasComeIn, _ := myRedis.RedisService.Exists(requestParams.StudentID + "*" + requestParams.CourseID).Result()
	if hasComeIn == 1{
		res.Code=types.RepetitiveSubmit
		c.JSON(http.StatusOK,res)
		return
	}
	myRedis.RedisService.Set(requestParams.StudentID+"*"+requestParams.CourseID,1,myRedis.RedisTimeOutForTemplate)

	//step4 判断学生是否合理(根据不同情况设置了黑名单,减少查表消耗性能)
	if hasDelStudent, _ := myRedis.RedisService.Exists("SD" + requestParams.StudentID).Result(); hasDelStudent == 1{
		res.Code=types.UserHasDeleted
		c.JSON(http.StatusOK,res)
		return
	}
	if hasHeiMingDan, _ :=myRedis.RedisService.Exists("SN"+requestParams.StudentID).Result(); hasHeiMingDan == 1{
		res.Code=types.UserNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if hasStudent, _ := myRedis.RedisService.Exists("S" + requestParams.StudentID).Result(); hasStudent != 1{
		var student types.TStudent
		firstForStudent := utils.Db.Where("user_id = ?", requestParams.StudentID).First(&student)
		if firstForStudent.Error != nil{
			myRedis.RedisService.Set("SN"+requestParams.StudentID,1,myRedis.RedisTimeOutForKeep)
			res.Code=types.StudentNotExisted
			c.JSON(http.StatusOK,res)
			return
		}

		if student.IsDel == 1{
			myRedis.RedisService.Set("SD"+requestParams.StudentID,1,myRedis.RedisTimeOutForKeep)
			res.Code=types.UserHasDeleted
			c.JSON(http.StatusOK,res)
			return
		}
		myRedis.RedisService.Set("S"+requestParams.StudentID,1,myRedis.RedisTimeOutForKeep)
	}

	//step5 判断课程是否存在(利用redis减轻压力)
	var tCourse types.Course
	tCourse.ID=cid
	_, errSid = myRedis.RedisService.Get(requestParams.CourseID).Result()
	if errSid != nil{
		firstForCourse := utils.Db.First(&tCourse)
		if firstForCourse.Error!=nil{
			res.Code=types.CourseNotExisted
			c.JSON(http.StatusOK,res)
			return
		}
		myRedis.RedisService.Set(strconv.Itoa(tCourse.ID),tCourse.CAP,myRedis.RedisTimeOutForKeep)
	}

	//step6 利用redis实时监控课程容量
	if myRedis.RedisService.Decr(requestParams.CourseID).Val()<0{
		myRedis.RedisService.Incr(requestParams.CourseID)
		if _, ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID); !ok{
			myRedis.CourseCapacityMap.Store(requestParams.CourseID, true)
		}
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
		return
	}

	//step7 进入抢客流程此处用事务写
	var updateRaw int64
	errUpdate := utils.Db.Transaction(func(tx *gorm.DB) error {
		recordLike := "%`"+requestParams.CourseID+"`%"
		record := "`"+requestParams.CourseID+"`;"
		errForCourse := tx.Model(&tCourse).Where("cap > ? ", 0).UpdateColumn("cap", gorm.Expr("cap - ?", 1))

		updateRaw=errForCourse.RowsAffected
		if errForCourse.RowsAffected<1{
			return nil
		}
		if errForCourse.Error !=nil {
			return errors.New(myRedis.ErrorForUpdateStore)
		}

		errForStudent := tx.Exec("UPDATE t_student SET course_record_id = CONCAT(course_record_id,?) WHERE user_id = ? and course_record_id not like ?",record,requestParams.StudentID,recordLike)
		//errForStudent := tx.Model(&student).Where("course_record_id not like ? ",recordLike).Update("course_record_id",student.CourseRecordId+"`"+requestParams.CourseID+"`;")

		if errForStudent.Error !=nil || errForStudent.RowsAffected<1{
			return errors.New(myRedis.ErrorForUpdateRecord)
		}
		return nil
	})

	//更新失败
	if errUpdate != nil && errUpdate.Error() == myRedis.ErrorForUpdateStore {
		myRedis.RedisService.Incr(requestParams.CourseID)
		if _, ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID);ok{
			myRedis.CourseCapacityMap.Delete(requestParams.CourseID)
		}
		res.Code=types.UnknownError
		c.JSON(http.StatusOK,res)
		return
	}

	//已抢到
	if errUpdate != nil && errUpdate.Error() == myRedis.ErrorForUpdateRecord{
		myRedis.RedisService.Incr(requestParams.CourseID)
		if _, ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID);ok{
			myRedis.CourseCapacityMap.Delete(requestParams.CourseID)
		}
		res.Code=types.StudentHasCourse
		c.JSON(http.StatusOK,res)
		return
	}

	//课已满
	if updateRaw<1 {
		if _,ok := myRedis.CourseCapacityMap.Load(requestParams.CourseID);!ok{
			myRedis.CourseCapacityMap.Store(requestParams.CourseID,true)
		}
		res.Code=types.CourseNotAvailable
		c.JSON(http.StatusOK,res)
	}

	myRedis.RedisService.Set(requestParams.StudentID+"->"+requestParams.CourseID,1,myRedis.RedisTimeOutForKeep)
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
		cid,err:=strconv.Atoi(strings.Trim(val,"`"))
		if err !=nil{
			continue
		}
		var tCourse types.Course
		tCourse.ID=cid
		utils.Db.First(&tCourse)
		var course =types.TCourse{CourseID: strconv.Itoa(tCourse.ID),TeacherID: tCourse.TeacherId,Name: tCourse.NAME}
		data=append(data,course)
	}
	res.Data.CourseList=data
	c.JSON(http.StatusOK,res)
	return
}