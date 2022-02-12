package router

import (
	"Register/controller"
	"Register/types"
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {



	g := r.Group("/api/v1")

	gob.Register(types.Teacher)
	gob.Register(types.Admin)
	gob.Register(types.Student)

	g.Use(sessions.Sessions(types.SessionName,types.Store))

	// 成员管理
	g.POST("/member/create", controller.MemberCreatePost)
	g.GET("/member", controller.MemberGetOne)
	g.GET("/member/list",controller.MemberGetList)
	g.POST("/member/update",controller.MemberUpdate)
	g.POST("/member/delete",controller.MemberDelete)

	// 登录
	g.POST("/auth/login",controller.MemberLogIn)
	g.POST("/auth/logout",controller.MemberLogOut)
	g.GET("/auth/whoami",controller.WhoAmI)

	// 排课
	g.POST("/course/create",controller.RangeCourseCreate)
	g.GET("/course/get",controller.RangeCourseGet)

	g.POST("/teacher/bind_course",controller.RangeBindCourse)
	g.POST("/teacher/unbind_course",controller.RangeUnbindCourse)
	//这个写错了,这个应该是teacherCourse
	g.GET("/teacher/get_course",controller.RangeGetTeacherCourse)
	g.POST("/course/schedule",controller.RangeSchedule)

	// 抢课
	g.POST("/student/book_course",controller.SelectCourseBookCourse)
	g.GET("/student/course",controller.SelectCourseGetCourse)

}
