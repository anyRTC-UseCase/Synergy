package utils

import (
	"ARTeamViewService/consts"
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
