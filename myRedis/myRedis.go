package myRedis

import (
	"github.com/go-redis/redis"
	"sync"
	"time"
)

var RedisService *redis.Client

var CourseCapacityMap sync.Map

const (
	RedisNetWork string = "tcp"
	RedisAddr string = "180.184.74.238"
	RedisPort string = "6379"
	RedisPrefix string = ""
	RedisTimeOutForKeep time.Duration = time.Duration(25)*time.Hour
	RedisTimeOutForTemplate time.Duration=time.Duration(3)*time.Second
	ErrorForUpdateStore string ="库存扣减失败"
	ErrorForUpdateRecord string ="学生重复抢客"
)

