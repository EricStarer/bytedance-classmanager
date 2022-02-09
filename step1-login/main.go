package main

import (
	"bytedance-classmanager/step1-login/router"
	"bytedance-classmanager/step1-login/util"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	//注册路由
	router.LoginRouter(r)
	//初始化数据库
	util.InitMysql()
	r.Run()

}
