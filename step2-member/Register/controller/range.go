package controller

import (
	"Register/myRedis"
	"Register/request"
	"Register/response"
	"Register/types"
	"Register/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

//创建课程
func RangeCourseCreate(c *gin.Context)  {
	var requestParams request.CreateCourseRequest
	var res response.CreateCourseResponse
	jsons:=utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)

	if len(requestParams.Name)<1 || requestParams.Cap<0{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var course = types.TCourseCwc{CourseName: requestParams.Name,Capacity: requestParams.Cap}
	utils.Db.Create(&course)
	res.Data.CourseID=strconv.FormatUint(course.ID,10)
	res.Code=types.OK
	//设置 redis key = couresId , val = cap
	set := myRedis.RedisService.Set(res.Data.CourseID, requestParams.Cap, myRedis.RedisTimeOut)
	fmt.Println(set.Err())
	c.JSON(http.StatusOK,res)
	return
}


//获取课程
func RangeCourseGet(c *gin.Context)  {
	var requestParams request.GetCourseRequest
	jsons := utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	var res response.GetCourseResponse
	var tCourse types.TCourseCwc
	id,err:=strconv.Atoi(requestParams.CourseID)
	if err!=nil ||len(requestParams.CourseID)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	tCourse.ID=(uint64)(id)
	first := utils.Db.First(&tCourse)
	if first.Error!=nil{
		res.Code=types.CourseNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	var data types.TCourse
	data.CourseID=requestParams.CourseID
	data.Name=tCourse.CourseName
	data.TeacherID=tCourse.TeacherID
	res.Code=types.OK
	res.Data=data
	c.JSON(http.StatusOK,res)
	return
}

//绑定课程

func RangeBindCourse(c *gin.Context)  {
	var requestParams request.BindCourseRequest
	var res response.BindCourseResponse
	jsons:=utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	var tCourse types.TCourseCwc
	id,errid:=strconv.ParseUint(requestParams.CourseID,10,64)
	_,errtid:=strconv.ParseUint(requestParams.TeacherID,10,64)
	if errid!=nil || errtid!=nil || len(requestParams.CourseID)<1 || len(requestParams.TeacherID)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	tCourse.ID=id
	first := utils.Db.First(&tCourse)
	if first.Error != nil{
		res.Code=types.CourseNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if len(tCourse.TeacherID)>0{
		res.Code=types.CourseHasBound
		c.JSON(http.StatusOK,res)
		return
	}
	//teacher 不需要做已落库校验
	tCourse.TeacherID=requestParams.TeacherID
	utils.Db.Model(&tCourse).Update("teacher_id",requestParams.TeacherID)
	var teacher types.TTeacher
	utils.Db.Where("user_id = ?",requestParams.TeacherID).First(&teacher)
	utils.Db.Model(&teacher).Update("teach_record_id",teacher.TeachRecordId+"`"+requestParams.CourseID+"`;")
	res.Code=types.OK
	c.JSON(http.StatusOK,res)
	return
}

//解绑课程

func RangeUnbindCourse(c *gin.Context)  {
	var requestParams request.UnbindCourseRequest
	var res response.UnbindCourseResponse
	jsons:=utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	cid,errCid:=strconv.ParseUint(requestParams.CourseID,10,64)
	_,errTid:=strconv.ParseUint(requestParams.TeacherID,10,64)
	if errCid != nil || errTid !=nil || len(requestParams.TeacherID)<1 || len(requestParams.CourseID)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var tCourse types.TCourseCwc
	tCourse.ID=cid
	first := utils.Db.First(&tCourse)
	if first.Error != nil{
		res.Code=types.CourseNotExisted
		c.JSON(http.StatusOK,res)
		return
	}
	if len(tCourse.TeacherID)<1{
		res.Code=types.CourseNotBind
		c.JSON(http.StatusOK,res)
		return
	}

	var teacher types.TTeacher
	update := utils.Db.Model(&tCourse).Update("teacher_id", "")
	if update.Error!=nil{
		fmt.Println(update.Error)
		return
	}
	utils.Db.Where("user_id = ?",requestParams.TeacherID).First(&teacher)
	teacher.TeachRecordId= strings.ReplaceAll(teacher.TeachRecordId, "`"+requestParams.CourseID+"`;", "")
	utils.Db.Model(&teacher).Update("teach_record_id",teacher.TeachRecordId)
	res.Code=types.OK
	c.JSON(http.StatusOK,res)
	return
}

//老师获取课程
func RangeGetTeacherCourse(c *gin.Context)  {
	var requestParams request.GetTeacherCourseRequest
	var res response.GetTeacherCourseResponse
	jsons := utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	if _,err:=strconv.Atoi(requestParams.TeacherID); err!=nil || len(requestParams.TeacherID)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	var teacher types.TTeacher
	utils.Db.Where("user_id = ?",requestParams.TeacherID).First(&teacher)
	//明天写
	var data [] *types.TCourse
	split := strings.Split(teacher.TeachRecordId, ";")
	for _,courseId := range split{
		cid,err := strconv.ParseUint(strings.Trim(courseId, "`"),10,64)
		if err!=nil{
			continue
		}
		var tCourse types.TCourseCwc
		var course types.TCourse
		tCourse.ID=cid
		utils.Db.First(&tCourse)
		course.Name=tCourse.CourseName
		course.TeacherID=tCourse.TeacherID
		course.CourseID=strconv.FormatUint(cid,10)
		data=append(data,&course)
	}
	res.Data.CourseList=data
	c.JSON(http.StatusOK,res)
	return
}

//排课表
func RangeSchedule(c *gin.Context)  {
	var requestParams request.ScheduleCourseRequest
	var res response.ScheduleCourseResponse
	jsons := utils.GetParams(c,requestParams)
	json.Unmarshal(jsons,&requestParams)
	if requestParams.TeacherCourseRelationShip == nil || len(requestParams.TeacherCourseRelationShip)<1{
		res.Code=types.ParamInvalid
		c.JSON(http.StatusOK,res)
		return
	}
	haveInterest:=make(map[string]int)
	haveCourse:=make(map[string]string)
	ans:=make(map[string]string)
	for k,v := range requestParams.TeacherCourseRelationShip {
		for _,courseId := range v {
			key:=k+"->"+courseId
			haveInterest[key]= 1
			haveCourse[courseId]=k
		}
	}
	var sum=0
	for k,_ := range requestParams.TeacherCourseRelationShip {
		used:=make(map[string]int)
		if find(k,haveInterest,haveCourse,used,ans){
			sum++
		}
	}
	teacherTeachCourse:=make(map[string]string)
	for courseId,teacherId := range ans {
		teacherTeachCourse[teacherId]=courseId
		delete(haveCourse,courseId)
		delete(requestParams.TeacherCourseRelationShip,teacherId)
	}
	for teacherId,_ := range requestParams.TeacherCourseRelationShip{
		for courseId,_ :=range haveCourse{
			teacherTeachCourse[teacherId]=courseId
			delete(haveCourse,courseId)
			break
		}
	}
	res.Data=teacherTeachCourse
	c.JSON(http.StatusOK,res)
	return
}

func find(teacherId string, interest map[string]int, course map[string]string, used map[string]int, ans map[string]string) bool {
	for courseId,_ := range course {
		_,ok1:=interest[teacherId+"->"+courseId]
		_,ok2:=used[courseId]
		if ok1 && !ok2{
			used[courseId]=1
			if _,ok:=ans[courseId];!ok || find(ans[courseId],interest,course,used,ans){
				ans[courseId]=teacherId
				return true
			}
		}
	}
	return false
}





