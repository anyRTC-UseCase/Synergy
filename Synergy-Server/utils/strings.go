package utils

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"strings"
)

// Equals 字符串是否相等
func Equals(firstStr, secStr string) bool {
	if strings.Compare(firstStr, secStr) == 0 {
		return true
	} else {
		return false
	}
}

// StrEmpty 是否是空字符串
func StrEmpty(str string) bool {
	return Equals(str, consts.StrEmpty)
}

// RedisUserKey token存入redis的key
func RedisUserKey(uid string) string {
	return consts.ARAnyRtcStr + consts.StrColon + consts.ARUIdStr + consts.StrColon + uid
}

// StrContains 字符串是否包含另一个字符串
func StrContains(firstStr, secStr string) bool {
	contains := strings.Contains(firstStr, secStr)
	return contains
}

func GetDeleteFilePrefix(url string) string {
	//https://teameeting.oss-cn-shanghai.aliyuncs.com/ar/synergy/X2177364628_60285144.mp4
	global.GLogger.Info("getDeleteFilePrefix url: ", url)

	split := strings.Split(url, "/")
	tmpPrefix := ""
	if len(split) == 6 {
		tmpPrefix = split[3] + "/" + split[4] + "/" + split[5]
	} else if len(split) == 5 {
		tmpPrefix = split[3] + "/" + split[4]
	}
	prefix := strings.Split(tmpPrefix, ".")
	global.GLogger.Info("getDeleteFilePrefix prefix[0]: ", prefix[0])
	return prefix[0]
}
