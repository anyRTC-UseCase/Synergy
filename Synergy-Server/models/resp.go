package models

type ApiJson struct {
	// 响应码
	Code int `json:"code"`
	// 消息
	Msg string `json:"msg"`
	// 响应数据
	Data interface{} `json:"data"`
}

/**
 * 无message的json响应数据格式,S是simple的简写
 */
type ApiSJson struct {
	// 响应码
	Code int `json:"code"`
	// 响应数据
	Data interface{} `json:"data"`
}

func ApiJsonResp(code int, msg string, cont interface{}) (apiJson *ApiJson) {
	return &ApiJson{code, msg, cont}
}

/**
 * 录制接口响应信息
 */
type ApiVodJson struct {
	// 响应码
	Code int `json:"code"`
	// 响应数据
	Body interface{} `json:"body"`
}
