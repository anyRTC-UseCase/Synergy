package utils

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"github.com/go-redis/redis/v8"
)

//根据uid删除前缀是此uid的所有token
func DelKeysByPrefix(str string) bool {
	keys := global.GRedisClient.Keys(global.Context, str+consts.StrAsterisk).Val()
	global.GLogger.Info("keys:", keys)
	if len(keys) > 0 {
		for i, _ := range keys {
			global.GRedisClient.Del(global.Context, keys[i])
		}
		afterKeys := global.GRedisClient.Keys(global.Context, str+consts.StrAsterisk).Val()
		if len(afterKeys) == 0 {
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

//获取Redis里面进入频道的主播uid的集合的key值
func GetRdsRoomUIdKey(roomId string) string {
	return "teamview:uids:" + roomId
}

//进入频道的主播uid保存至redis中
func ZAddRdsRoomUId(roomId, uid string) {
	rdsKey := GetRdsRoomUIdKey(roomId)
	//global.GRedisClient.ZAdd(rdsKey, redis.Z{float64(FormatNowUnix()), uid})
	global.GRedisClient.ZAdd(global.Context, rdsKey, &redis.Z{Score: float64(FormatNowUnix()), Member: uid})
}

//判断redis中key是否存在
func QueryKeyIsNotExists(roomId string) bool {
	rdsKey := GetRdsRoomUIdKey(roomId)
	val := global.GRedisClient.Exists(global.Context, rdsKey).Val()
	var flag = true
	if val == 0 {
		return false
	}
	return flag
}

//离开频道的主播uid从redis中删除,并返回redis集合中的uid个数
func ZRemRdsRoomUId(roomId, uid string) int64 {
	rdsKey := GetRdsRoomUIdKey(roomId)
	global.GRedisClient.ZRem(global.Context, rdsKey, uid)
	uidNums := GetRdsRoomUIdNums(roomId)
	if uidNums == 0 {
		DelRdsRoomUIdKey(roomId)
	}
	return uidNums
}

//查询频道的主播uid数量
func GetRdsRoomUIdNums(roomId string) int64 {
	rdsKey := GetRdsRoomUIdKey(roomId)
	return global.GRedisClient.ZCard(global.Context, rdsKey).Val()
}

//删除rdskey
func DelRdsRoomUIdKey(roomId string) {
	rdsKey := GetRdsRoomUIdKey(roomId)
	global.GRedisClient.Del(global.Context, rdsKey)
}
