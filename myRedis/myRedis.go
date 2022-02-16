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
	RedisTimeOutForKeep time.Duration = time.Duration(49)*time.Hour
	RedisTimeOutForTemplate time.Duration=time.Duration(5)*time.Second
	RedisTimeOutForReadStore time.Duration=time.Duration(1)*time.Second
	ErrorForGetCourseAgain string = "重复抢课"
	ErrorForUpdateStore string ="库存扣减失败"
	ErrorForUpdateRecord string ="学生重复抢客"
)


