package main

import (
	"bytedance-classmanager/src/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

//数据库字段和requeset的请求中（types.go）中相对应
type course struct {
	ID        string `json:"CourseID"`
	NAME      string `json:"Name"`
	CAP       int    `json:"Cap"`
	TeacherId string `json:"TeacherID"`
}

var err error
var db *gorm.DB //数据库连接
//var st []bool
var st map[string]string
var TeacherCourseRelationShip map[string][]string
var res map[string]string

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
	c.ShouldBindJSON(&u) //gin框架参数绑定，传入的指针
	db.Create(&u)        //创建记录
	db.Last(&u)          //db主键id是自增的,因此查询最新插入的记录,从而获取id值
	c.JSON(http.StatusOK, gin.H{ //返回值
		"Code": types.OK,
		"Data": gin.H{
			"CourseID": u.ID,
		},
	})
}

//获取课程信息
func course_get(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u)        //别用错方法，找了半天bug
	id, _ := strconv.Atoi(u.ID) //将id转为整型，因为数据库中id字段为整型
	db.First(&u, id)            //查询
	if u.NAME == "" {           //如果没查询到就返回Errno,CourseNotExisted
		c.JSON(http.StatusOK, gin.H{
			"Code": types.CourseNotExisted,
			"Data": gin.H{
				"CourseID":  u.ID,
				"Name":      u.NAME,
				"TeacherID": u.TeacherId,
			},
		})
		return
	} //已查询到返回OK
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
	c.ShouldBindJSON(&u)               //gin，参数绑定
	course_id, _ := strconv.Atoi(u.ID) //id转为数字，理由同上
	teacher_id := u.TeacherId
	db.First(&u, course_id) //先查询对应课程的记录
	if u.TeacherId == "" {  //该课程还未绑定老师，直接更新记录进行绑定
		db.Debug().Model(&u).Update("TeacherId", teacher_id) //仅修改部分
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
func UnBind_course(c *gin.Context) {
	var u course
	c.ShouldBindJSON(&u)
	course_id, _ := strconv.Atoi(u.ID)
	db.First(&u, course_id)                      //查询对应课程
	db.Debug().Model(&u).Update("TeacherId", "") //仅修改部分,teacher_id字段修改为空值
	c.JSON(http.StatusOK, gin.H{
		"Code": types.OK,
	})
}

//查询老师下的所有课程（未完成)，CourseList []*TCourse什么意思？
func teacher_course_get(c *gin.Context) {
	var u course                                     //接收request中的数据
	var res []course                                 //根据teacher_id得到记录的结果集
	c.ShouldBindJSON(&u)                             //参数绑定
	db.Where("teacher_id=?", u.TeacherId).Find(&res) //查询
	resp := new(types.GetTeacherCourseResponse)      //response结构体
	resp.Code = types.OK
	arrlen := len(res)                                //记录的个数
	ans := make([]*types.TCourse, arrlen, arrlen+100) //cap开大一点
	for i := 0; i < arrlen; i++ {                     //将u中的三个字段转移到ans中
		ans[i] = new(types.TCourse) //指针要先初始化！！！
		ans[i].TeacherID = res[i].TeacherId
		ans[i].CourseID = res[i].ID
		ans[i].Name = res[i].NAME
	}
	resp.Data.CourseList = ans
	c.JSON(http.StatusOK, resp)
}
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

//排课求解器
func schedule(c *gin.Context) { //匈牙利算法
	var u types.ScheduleCourseRequest
	c.ShouldBindJSON(&u)
	TeacherCourseRelationShip = u.TeacherCourseRelationShip
	res = make(map[string]string)
	for k, _ := range TeacherCourseRelationShip {
		for k, _ := range res {
			delete(res, k)
		}
		find(k)
	}
	resp := new(types.ScheduleCourseResponse)
	resp.Code = types.OK
	resp.Data = res
	c.JSON(http.StatusOK, resp)

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
	g.POST("/course/schedule", schedule)
	r.Run()
}
