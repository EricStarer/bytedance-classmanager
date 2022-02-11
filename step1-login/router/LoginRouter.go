package router

import (
	"bytedance-classmanager/step1-login/handler"
	"github.com/gin-gonic/gin"
)

func LoginRouter(e *gin.Engine) {
	g := e.Group("/api/v1")
	{
		g.POST("/auth/login", handler.LoginHandler)
		g.POST("/auth/logout", handler.LogoutHandler)
		g.GET("/auth/whoami", handler.WhoamiHandler)
	}
}
