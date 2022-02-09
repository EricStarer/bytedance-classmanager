package router

import (
	"Register/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", controller.MemberCreatePost)
	g.GET("/member", controller.MemberGetOne)
	g.GET("/member/list", controller.MemberGetList)
	g.POST("/member/update", controller.MemberUpdate)
	g.POST("/member/delete", controller.MemberDelete)

	// 登录

	g.POST("/auth/login", controller.Login)
	g.POST("/auth/logout", controller.Logout)
	g.GET("/auth/whoami", controller.Whoami)

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
