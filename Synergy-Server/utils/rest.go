package utils

import (
	"ARTeamViewService/global"
	"io/ioutil"
	"net/http"
	"strings"
)

/**
 * HTTP Post 请求
 * 添加了 HTTP 基本认证
 */
func PostStatusRequest(url string, contentType string, body, base64Credentials string) (statusCode int, ret string) {
	global.GLogger.Info("postRequest url: ", url)
	global.GLogger.Info("postRequest body: ", body)
	//1.http请求客户端
	client := &http.Client{}

	//2.获取request
	reqest, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		global.GLogger.Error("http.NewRequest err: ", err)
		return
	}

	//3.header增加Basic HTTP认证,设置Content-Type
	reqest.Header.Add("Authorization", "Basic "+base64Credentials)
	reqest.Header.Set("Content-Type", contentType)
	//7.发送请求
	resp, err := client.Do(reqest)
	if err != nil {
		global.GLogger.Error("client.Do err: ", err)
		return
	}
	defer resp.Body.Close()

	//4.读取body
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		global.GLogger.Error("ioutil.ReadAll err: ", err)
		return
	}

	return resp.StatusCode, string(response)
}
