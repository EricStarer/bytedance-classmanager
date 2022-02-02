package main

import (
	"bytedance-classmanager/src/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

type course struct {
	ID        string `json:"CourseID"`
	NAME      string `json:"Name"`
	CAP       int    `json:"Cap"`
	TeacherId string `json:"TeacherID"`
}

var err error
var db *gorm.DB //数据库连接

//连接数据库函数
func connect() error {
	db, err = gorm.Open("mysql", "root:bytedancecamp@(180.184.74.238)/test_syp?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		return nil
	}
	return db.DB().Ping()
}

//创建课程函数
func course_create(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u)
	db.Create(&u)
	db.Last(&u)
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
		"Data": gin.H{
			"CourseID": u.ID,
		},
	})
}

//获取课程信息
func course_get(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u) //别用错方法，找了半天bug
	id, _ := strconv.Atoi(u.ID)
	db.First(&u, id)
	if u.NAME == "" {
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseNotExisted,
			"Data": gin.H{
				"CourseID":  u.ID,
				"Name":      u.NAME,
				"TeacherID": u.TeacherId,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
		"Data": gin.H{
			"CourseID":  u.ID,
			"Name":      u.NAME,
			"TeacherID": u.TeacherId,
		},
	})
}

//绑定课程
func Bind_Course(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u)
	course_id, _ := strconv.Atoi(u.ID)
	teacher_id := u.TeacherId
	db.First(&u, course_id)
	if u.TeacherId == "" {
		db.Debug().Model(&u).Update("TeacherId", teacher_id) //仅修改部分
		c.JSON(http.StatusOK, gin.H{
			"Code": types.OK,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseHasBound,
		})
	}
}

//解绑，有疑问:异常情况
func UnBind_course(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u)
	course_id, _ := strconv.Atoi(u.ID)
	db.First(&u, course_id)
	db.Debug().Model(&u).Update("TeacherId", "") //仅修改部分
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
	})
}

//查询老师下的所有课程（未完成)，CourseList []*TCourse什么意思？
func teacher_course_get(c *gin.Context) {
	var u course
	var res course[]
	c.ShouldBindJSON(&u)
	db.Where("teacher_id=?", u.TeacherId).Find(&res)
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
	})
}
func main() {
	e := connect()
	//连接失败
	if e != nil {
		panic(e)
	}
	defer db.Close()
	r := gin.Default()
	g := r.Group("/api/v1")
	// 排课
	g.POST("/course/create", course_create)
	g.GET("/course/get", course_get)

	g.POST("/teacher/bind_course", Bind_Course)
	g.POST("/teacher/unbind_course", UnBind_course)
	g.GET("/teacher/get_course", teacher_course_get)
	g.POST("/course/schedule")
	r.Run()
}
