package main

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"mergeVersion1/myRedis"
	"mergeVersion1/router"
	"mergeVersion1/utils"
)

func main() {

	//启动redis
	myRedis.RedisService = redis.NewClient(&redis.Options{Network: myRedis.RedisNetWork, Addr: myRedis.RedisAddr + ":" + myRedis.RedisPort, IdleTimeout: myRedis.RedisTimeOutForKeep})

	//定义一个路由
	r := gin.Default()

	//添加日志文件
	r.Use(utils.LoggerToFile())

	//导入模板
	r.LoadHTMLGlob("view/*")

	//导入静态资源
	r.Static("/static", "./static")

	//注册路由
	router.RegisterRouter(r)
	//连接字节的数据库
	utils.InitMysql()

	//可以取消表名的复数形式，使得表名和结构体名称一致
	utils.Db.SingularTable(true)

	r.Run(":8088")
}
