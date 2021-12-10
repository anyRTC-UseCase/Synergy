package utils

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"ARTeamViewService/models"
	"github.com/kataras/iris/v12"
	json "github.com/wulinlw/jsonAmbiguity"
)

// ARIrisReadJson iris读取json数据
func ARIrisReadJson(ctx iris.Context, params interface{}) error {
	err := ctx.ReadJSON(&params)
	if err != nil {
		global.GLogger.Error(err)
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		ctx.StopExecution()
		return err
	}
	return nil
}

//ARUnmarshalJson　解析数据json数据
func ARUnmarshalJson(data string, params interface{}) error {
	err := json.Unmarshal([]byte(data), &params)
	if err != nil {
		global.GLogger.Error("ARUnmarshalJson data: ", data, ", json Unmarshal error: ", err)
		return err
	}
	return err
}

//ARParseJson　Interface 解析数据json数据
func ARParseJson(data interface{}, jsonObject interface{}) (int, string, interface{}) {
	byteData, marshalErr := json.Marshal(data)
	if marshalErr != nil {
		global.GLogger.Error("json Marshal error: ", marshalErr)
		return consts.ErrCodeInternal, consts.ErrJsonMarshal, nil
	} else {
		unmarshalErr := json.Unmarshal(byteData, &jsonObject)
		if unmarshalErr != nil {
			global.GLogger.Error("ARParseJson data: ", string(byteData), ", json Unmarshal error: ", unmarshalErr)
			return consts.ErrCodeInternal, consts.ErrJsonUnmarshal, nil
		}
	}
	return consts.ErrCodeOk, consts.ErrOk, nil
}
