package global

import (
	"ARTeamViewService/conf"
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

var GUseConf string
var GConfig conf.ServerConfig

var GDb *sql.DB
var GDbCm *sql.DB

// GRedisClient redis客户端
var GRedisClient *redis.Client

var GLogger *logrus.Logger

var Context context.Context
