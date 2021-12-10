package utils

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"encoding/base64"
)

/**
 * 使用客户ID和客户密钥生成一个Credential（使用 Base64 算法编码）
 */
func Base64Credentials() string {
	var str = global.GConfig.CustomerId + consts.StrColon + global.GConfig.CustomerCert
	strbytes := []byte(str)
	encoded := base64.StdEncoding.EncodeToString(strbytes)
	global.GLogger.Info("base64Credentials: ", encoded)
	return encoded
}
