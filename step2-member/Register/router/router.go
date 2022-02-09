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

	g.Use(sessions.Sessions(types.SessionName, types.Store))

	// 成员管理

	g.POST("/member/create", controller.MemberCreatePost)
	g.GET("/member", controller.MemberGetOne)
	g.GET("/member/list", controller.MemberGetList)
	g.POST("/member/update", controller.MemberUpdate)
	g.POST("/member/delete", controller.MemberDelete)

	// 登录

	g.POST("/auth/login", controller.MemberLogIn)
	g.POST("/auth/logout", controller.MemberLogOut)
	g.GET("/auth/whoami", controller.WhoAmI)

	// 排课
	g.POST("/course/create")
	g.GET("/course/get")

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
