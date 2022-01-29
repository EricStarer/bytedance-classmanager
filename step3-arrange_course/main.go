package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"project02/types"
	"strconv"
)

type course struct {
	ID        string
	NAME      string
	CAP       int
	TeacherId string
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
	name := c.PostForm("Name")
	cap, e := strconv.Atoi(c.PostForm("Cap"))
	if e != nil {
		panic(e)
	}
	//db.AutoMigrate(&types.CreateCourseRequest{})
	t1 := course{NAME: name, CAP: cap}
	db.Create(&t1)
	var t2 course
	db.Last(&t2)
	fmt.Println("%#v\n", t2)
  //返回值，有疑问，types.go中CreateCourseResponse中Data
	c.JSON(http.StatusOK, gin.H{
		"Code":     types.OK,
		"courseId": t2.ID,
	})
}
//获取课程信息
func course_get(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u) //别用错方法，找了半天bug
	//fmt.Printf("u=%#v\n", u)
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
	g.GET("/course/get",course_get)

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")
	r.Run()
}
