package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mergeVersion1/myRedis"
	"mergeVersion1/request"
	"mergeVersion1/response"
	"mergeVersion1/types"
	"mergeVersion1/utils"
	"net/http"
	"strconv"
)

var st map[string]string
var TeacherCourseRelationShip map[string][]string
var res map[string]string

//创建课程函数
func Course_create(c *gin.Context) {
	var u types.Course
	resp := new(response.CreateCourseResponse)
	c.ShouldBindJSON(&u)
	if u.NAME == "" || u.CAP == 0 { //参数不合法
		resp.Code = types.ParamInvalid
		resp.Data.CourseID = ""
		c.JSON(http.StatusOK, resp)
		return
	}
	utils.Db.Create(&u)
	u.CourseID = strconv.Itoa(u.ID)
	fmt.Printf("%#v\n", u)
	resp.Code = types.OK
	resp.Data.CourseID = u.CourseID
	//设置 redis key = couresId , val = cap
	fmt.Println(u.CAP) //整合测试要看看cap对不对
	set := myRedis.RedisService.Set(resp.Data.CourseID, u.CAP, myRedis.RedisTimeOutForKeep)
	fmt.Println(set.Err())
	c.JSON(http.StatusOK, resp)
}

//获取课程信息
func Course_get(c *gin.Context) {
	var u types.Course
	var id string
	c.ShouldBindJSON(&u) //别用错方法，找了半天bug
	if u.CourseID == "" {
		id = c.Query("CourseID")
	} else {
		id = u.CourseID
	}
	if id == "" { //参数不合法
		resp := new(response.GetCourseResponse)
		resp.Code = types.ParamInvalid
		resp.Data.CourseID = id
		resp.Data.Name = ""
		resp.Data.TeacherID = ""
		c.JSON(http.StatusOK, resp)
		return
	}
	utils.Db.First(&u, id) //查询
	if u.NAME == "" {      //如果没查询到就返回Errno,CourseNotExisted
		resp := new(response.GetCourseResponse)
		resp.Code = types.CourseNotExisted
		resp.Data.CourseID = id
		resp.Data.Name = u.NAME
		resp.Data.TeacherID = u.TeacherId
		c.JSON(http.StatusOK, resp)
		return
	}
	resp := new(response.GetCourseResponse)
	resp.Code = types.OK
	resp.Data.CourseID = id
	resp.Data.Name = u.NAME
	resp.Data.TeacherID = u.TeacherId
	c.JSON(http.StatusOK, resp)
}

//绑定课程
func Bind_Course(c *gin.Context) {
	var u, t types.Course
	c.ShouldBindJSON(&u) //gin，参数绑定
	course_id := u.CourseID
	teacher_id := u.TeacherId
	if teacher_id == "" || course_id == "" { //有空值，参数不合法
		c.JSON(http.StatusOK, gin.H{
			"Code": types.ParamInvalid,
		})
		return
	}
	utils.Db.First(&t, course_id) //先查询对应课程的记录
	if t.NAME == "" {             //课程不存在
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseNotExisted,
		})
		return
	}
	fmt.Println(t)
	if t.TeacherId == "" { //该课程还未绑定老师，直接更新记录进行绑定
		utils.Db.Debug().Model(&t).Update("TeacherId", teacher_id) //仅修改部分
		c.JSON(http.StatusOK, gin.H{
			"Code": types.OK,
		})
	} else { //老师字段已经有人了，返回课程已被绑定的ErrNO:CourseHasBound
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseHasBound,
		})
	}
}

//解绑，有疑问:异常情况
//已完成
func UnBind_course(c *gin.Context) {
	var u, t types.Course
	c.ShouldBindJSON(&u)
	course_id := u.CourseID
	if u.TeacherId == "" || course_id == "" { //有空值，参数不合法
		c.JSON(http.StatusOK, gin.H{
			"Code": types.ParamInvalid,
		})
		return
	}
	utils.Db.First(&t, course_id) //查询对应课程
	if t.NAME == "" {             //不存在这个课程
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseNotExisted,
		})
		return
	}
	if t.TeacherId == "" {
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseNotBind,
		})
		return
	}
	utils.Db.Debug().Model(&u).Update("TeacherId", "") //仅修改部分,teacher_id字段修改为空值
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
	})
}

//查询老师下的所有课程（未完成)，CourseList []*TCourse什么意思？
//已完成
func Teacher_course_get(c *gin.Context) {
	var u types.Course     //接收request中的数据
	var res []types.Course //根据teacher_id得到记录的结果集
	c.ShouldBindJSON(&u)   //参数绑定
	if u.TeacherId == "" { //get传参
		u.TeacherId = c.Query("TeacherID")
	}
	fmt.Println(u)
	if u.TeacherId == "" { //参数不合法
		resp := new(response.GetTeacherCourseResponse)
		resp.Code = types.ParamInvalid
		resp.Data.CourseList = nil
		c.JSON(http.StatusOK, resp)
		return
	}
	utils.Db.Where("teacher_id=?", u.TeacherId).Find(&res) //查询
	resp := new(response.GetTeacherCourseResponse)         //response结构体
	resp.Code = types.OK
	arrlen := len(res)                                //记录的个数
	ans := make([]*types.TCourse, arrlen, arrlen+100) //cap开大一点
	for i := 0; i < arrlen; i++ {                     //将u中的三个字段转移到ans中
		ans[i] = new(types.TCourse) //指针要先初始化！！！
		ans[i].TeacherID = res[i].TeacherId
		ans[i].CourseID = strconv.Itoa(res[i].ID)
		ans[i].Name = res[i].NAME
	}
	resp.Data.CourseList = ans
	c.JSON(http.StatusOK, resp)
}

/*
匈牙利算法find函数
*/
func find(x string) bool {
	t := TeacherCourseRelationShip[x]
	for _, s := range t {
		if st[s] == "" {
			st[s] = "1"
			if res[s] == "" || find(res[s]) {
				res[s] = x
				return true
			}
		}
	}
	return false
}

/*
排课求解器
已处理异常：参数为空值（参数不合法）
测试数据（json）：
{"TeacherCourseRelationShip":{"1":["1","2","3"],"2":["1"],"3":["3"]}}
*/
func Schedule(c *gin.Context) { //匈牙利算法
	var u request.ScheduleCourseRequest
	c.ShouldBindJSON(&u)
	TeacherCourseRelationShip = u.TeacherCourseRelationShip
	res, st = make(map[string]string), make(map[string]string)
	resp := new(response.ScheduleCourseResponse)
	resp.Code = types.OK
	resp.Data = res
	for k, _ := range TeacherCourseRelationShip {
		if k == "" { //如果teacherid是空值
			resp.Code = types.ParamInvalid
			res = make(map[string]string)
			break
		}
		tmp := TeacherCourseRelationShip[k]
		for _, t := range tmp { //如果课程为空值
			if t == "" {
				resp.Code = types.ParamInvalid
				res = make(map[string]string) //重置结果为空
				break
			}
		}
		for k, _ := range res {
			delete(st, k)
		}
		find(k)
	}
	c.JSON(http.StatusOK, resp)
}
