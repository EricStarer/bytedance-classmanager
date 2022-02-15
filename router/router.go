package router

import (
	"mergeVersion1/controller"
	"mergeVersion1/types"
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
	g.POST("/auth/login",controller.LoginHandler)
	g.POST("/auth/logout",controller.LogoutHandler)
	g.GET("/auth/whoami",controller.WhoamiHandler)

	// 排课

	g.POST("/course/create", controller.Course_create)
	g.GET("/course/get", controller.Course_get)

	g.POST("/teacher/bind_course", controller.Bind_Course)
	g.POST("/teacher/unbind_course", controller.UnBind_course)
	g.GET("/teacher/get_course", controller.Teacher_course_get)
	g.POST("/course/schedule", controller.Schedule)

	// 抢课
	g.POST("/student/book_course",controller.SelectCourseBookCourse)
	g.GET("/student/course",controller.SelectCourseGetCourse)

}
