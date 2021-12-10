package utils

import (
	"ARTeamViewService/global"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	json "github.com/wulinlw/jsonAmbiguity"
	"math/rand"
	"path"
	"strconv"
	"time"
)

/**
 * 获取当前日期时间戳(秒)
 */
func FormatNowUnix() int64 {
	return time.Now().Unix()
}

/**
 * 转换时间戳为日期
 */
func FormatFullDate(timestamp int64) string {
	t := time.Unix(timestamp, 0)
	dataStr := t.Format("2006-01-02 15:04:05")
	return dataStr
}

/*
 H function for MD5 algorithm (returns a lower-case hex MD5 digest)
*/
func MD5(data string) string {
	digest := md5.New()
	digest.Write([]byte(data))
	return hex.EncodeToString(digest.Sum(nil))
}

/**
 * 转化时间戳
 */
func FormatDateId() string {
	t := time.Now()
	dataStr := t.Format("20060102150405")
	return dataStr
}

/**
 * 随机首字母非0数字长度
 */
func RandomNoZeroNumberString(length int) string {
	str := "123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	str_extra := "0123456789"
	bytes_extra := []byte(str_extra)
	r_extra := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length-1; i++ {
		result = append(result, bytes_extra[r_extra.Intn(len(bytes_extra))])
	}
	return string(result)
}

/**
 * 获取仅仅文件后缀
 */
func GetFileSuffix(fullFileName string) string {
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fullFileName) //获取文件名带后缀
	global.GLogger.Info("filenameWithSuffix: ", filenameWithSuffix)
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀
	global.GLogger.Info("fileSuffix: ", fileSuffix)

	return fileSuffix
}

/**
 * HmacSHA1加密算法
 */
func HmacSHA1(key, data string) string {
	mac := hmac.New(sha1.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

/**
 * 随机字符串
 */
func RandomCharString(length int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}

	return string(result)
}

// StrVal 获取变量的字符串值
// 浮点型 3.0将会转换成字符串3, "3"
// 非数值或字符类型的变量将会被转换成JSON格式字符串
func StrVal(value interface{}) string {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {
	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, err := json.Marshal(value)
		if err != nil {
			global.GLogger.Info("json Marshal error: ", err)
		}
		key = string(newValue)
	}

	return key
}
