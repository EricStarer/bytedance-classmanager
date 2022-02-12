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
	RedisTimeOut time.Duration = time.Duration(72)*time.Hour
)


