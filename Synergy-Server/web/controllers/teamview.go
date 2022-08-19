package controllers

import (
	"ARTeamViewService/consts"
	"ARTeamViewService/global"
	"ARTeamViewService/models"
	"ARTeamViewService/services"
	"ARTeamViewService/utils"
	"ARTeamViewService/web/middleware"
	"github.com/kataras/iris/v12"
	"sync"
)

type ApiUser interface {
}

var insUser *objApiUser
var onceUser sync.Once

type objApiUser struct {
	dbOpened bool
}

func GetInsApiUser() *objApiUser {
	onceUser.Do(func() {
		insUser = &objApiUser{
			dbOpened: false,
		}
	})
	return insUser
}

func NewApiUser() ApiUser {
	return &objApiUser{}
}

// SignIn godoc
// @summary 用户登录
// @description 用户登录
// @tags teamview
// @accept json
// @product json
// @param default body models.ReqSignIn true "用户登录参数"
// @success 200 {object} models.RespSignIn
// @router /teamview/signIn [post]
func (o *objApiUser) SignIn(ctx iris.Context) {
	global.GLogger.Info("SignIn called")
	var params models.ReqSignIn
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var userName = params.UserName
	var uid = params.UId
	var workName = params.WorkName
	var userType = params.UserType

	if utils.StrEmpty(userName) || utils.StrEmpty(uid) {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamMissing, consts.ErrParamMissing, nil))
		return
	}
	if userType < consts.UserTypeTerminal || userType > consts.UserTypeAdmin {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}

	code, msg, content := services.GetInsUserSvr().SignInSvr(userName, uid, workName, userType, ctx)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// TeamViewRtcNotify godoc
// @summary 消息回调
// @description 消息回调
// @tags teamview
// @accept json
// @product json
// @param default body models.ReqTeamViewRtcNotify true "消息回调"
// @success 200 {object} models.ApiJson
// @router /teamview/teamViewRtcNotify [post]
func (o *objApiUser) TeamViewRtcNotify(ctx iris.Context) {
	global.GLogger.Info("===>TeamViewRtcNotify called")

	var params models.ReqTeamViewRtcNotify
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}

	global.GLogger.Info("===>rtcNotify params: ", utils.StrVal(params))

	//2.接收到的签名
	signature := ctx.GetHeader("Ar-Signature")
	global.GLogger.Info("receive signature: ", signature)

	//3.读取的body
	buf := make([]byte, consts.IntOneZeroTwoFour)
	n, _ := ctx.Request().Body.Read(buf)

	//4.生成的签名
	secret := global.GConfig.DevpSecret
	sha1Tmp := utils.HmacSHA1(secret, string(buf[0:n]))
	global.GLogger.Info("old signature: ", sha1Tmp)

	if !utils.Equals(signature, sha1Tmp) {
		global.GLogger.Error("signature is invalid!")
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeInvalidRequest, consts.ErrInvalidRequest, nil))
		return
	} else {
		global.GLogger.Info("signature is valid!")
		var eventType = params.Payload.EventType
		var roomId = params.Payload.ChannelName
		var uid = params.Payload.UId

		code, msg, content := services.GetInsUserSvr().V2TeamViewRtcNotifySvr(roomId, uid, eventType, string(utils.StrVal(params)))
		ctx.JSON(models.ApiJsonResp(code, msg, content))
	}
}

// TeamViewVodNotify godoc
// @summary 录像回调
// @description 录像回调
// @tags teamview
// @accept json
// @product json
// @param default body models.ReqTeamViewVodNotify true "录像回调"
// @success 200 {object} models.ApiJson
// @router /teamview/teamViewVodNotify [post]
func (o *objApiUser) TeamViewVodNotify(ctx iris.Context) {
	global.GLogger.Info("===>TeamViewVodNotify called")

	var params models.ReqTeamViewVodNotify
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}

	global.GLogger.Info("===>vodNotify params: ", utils.StrVal(params))

	//2.接收到的签名
	signature := ctx.GetHeader("Ar-Signature")

	//3.读取的body
	buf := make([]byte, consts.IntOneZeroTwoFour)
	n, _ := ctx.Request().Body.Read(buf)

	//4.生成的签名
	secret := global.GConfig.DevpSecret
	sha1 := utils.HmacSHA1(secret, string(buf[0:n]))

	if !utils.Equals(signature, sha1) {
		global.GLogger.Error("signature is invalid!")
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeInvalidRequest, consts.ErrInvalidRequest, nil))
		return
	} else {
		global.GLogger.Info("signature is valid!")
		var eventType = params.EventType
		var roomId = params.Payload.CName
		var uid = params.Payload.UId
		var sendTs = params.Payload.SendTs
		var details = params.Payload.Details

		//code, msg, content := services.GetInsUserSvr().TeamViewVodNotifySvr(roomId, uid, eventType, utils.StrVal(params), sendTs)
		code, msg, content := services.GetInsUserSvr().V2TeamViewVodNotifySvr(details, roomId, uid, eventType, utils.StrVal(params), sendTs)
		ctx.JSON(models.ApiJsonResp(code, msg, content))
	}
}

// InsertRoom godoc
// @summary 创建房间
// @description 创建房间
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqInsertRoom true "创建房间"
// @success 200 {object} models.RespInsertRoom
// @router /teamview/insertRoom [post]
func (o *objApiUser) InsertRoom(ctx iris.Context) {
	global.GLogger.Info("InsertRoom called")
	var params models.ReqInsertRoom
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var isJoin = params.IsJoin

	if isJoin < consts.IsJoinRoom || isJoin > consts.IsNotJoinRoom {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}

	var uid = middleware.JWTUserId(ctx)
	code, msg, content := services.GetInsUserSvr().V2InsertRoomSvr(uid, isJoin)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// GetUserInfo godoc
// @summary 获取用户信息
// @description 获取用户信息
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqUId true "获取用户信息"
// @success 200 {object} models.UserInfo
// @router /teamview/getUserInfo [post]
func (o *objApiUser) GetUserInfo(ctx iris.Context) {
	global.GLogger.Info("GetUserInfo called")
	var params models.ReqUId
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var uid = params.UId
	if utils.StrEmpty(uid) {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamMissing, consts.ErrParamMissing, nil))
		return
	}

	code, msg, content := services.GetInsUserSvr().GetUserInfoSvr(uid)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// JoinRoom godoc
// @summary 加入房间
// @description 加入房间
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqJoinRoom true "加入房间参数"
// @success 200 {object} models.RespJoinRoom
// @router /teamview/joinRoom [post]
func (o *objApiUser) JoinRoom(ctx iris.Context) {
	global.GLogger.Info("JoinRoom called")
	var params models.ReqJoinRoom
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var roomId = params.RoomId
	var userRole = params.UserRole

	if utils.StrEmpty(roomId) {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamMissing, consts.ErrParamMissing, nil))
		return
	}
	if userRole < consts.UserRoleHost || userRole > consts.UserRoleAudience {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}

	var uid = middleware.JWTUserId(ctx)
	code, msg, content := services.GetInsUserSvr().V2JoinRoomSvr(roomId, uid, userRole)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// LeaveRoom godoc
// @summary 用户离开房间
// @description 用户离开房间
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqLeaveRoom true "用户离开房间"
// @success 200 {object} models.ApiJson
// @router /teamview/leaveRoom [post]
func (o *objApiUser) LeaveRoom(ctx iris.Context) {
	global.GLogger.Info("LeaveRoom called")
	var params models.ReqLeaveRoom
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var roomId = params.RoomId
	var uid = params.UId

	if utils.StrEmpty(roomId) || utils.StrEmpty(uid) {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamMissing, consts.ErrParamMissing, nil))
		return
	}

	code, msg, content := services.GetInsUserSvr().V2LeaveRoomSvr(roomId, uid)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// GetRoomList godoc
// @summary 获取房间列表
// @description 获取房间列表
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqGetRoomList true "获取房间列表参数"
// @success 200 {object} models.RespOngoingRoomList "进行中房间列表"
// @success 2002 {object} models.RespFinishedRoomList "结束房间列表(为了区分文档，code使用2002)"
// @router /teamview/getRoomList [post]
func (o *objApiUser) GetRoomList(ctx iris.Context) {
	global.GLogger.Info("GetRoomList called")
	var params models.ReqGetRoomList
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var pageNum = params.PageNum
	var pageSize = params.PageSize
	var roomState = params.RoomState

	if pageSize < consts.IntOne || pageNum < consts.IntOne || roomState < consts.RoomStateClosed || roomState > consts.RoomStateMixed {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}

	code, msg, content := services.GetInsUserSvr().GetRoomListSvr(roomState, pageNum, pageSize)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// GetSpecialist godoc
// @summary 获取专家列表
// @description 获取专家列表
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ParamPageList true "获取专家列表参数"
// @success 200 {object} models.RespGetSpecialist
// @router /teamview/getSpecialist [post]
func (o *objApiUser) GetSpecialist(ctx iris.Context) {
	global.GLogger.Info("GetSpecialist called")
	var params models.ParamPageList
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var pageNum = params.PageNum
	var pageSize = params.PageSize

	if pageSize < consts.IntOne || pageNum < consts.IntOne {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}
	var uid = middleware.JWTUserId(ctx)

	code, msg, content := services.GetInsUserSvr().GetSpecialistSvr(pageNum, pageSize, uid)
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}

// InsertUserOnlineInfo godoc
// @summary 记录用户在线心跳包信息
// @description 记录用户在线心跳包信息
// @tags teamview
// @accept json
// @product json
// @param Authorization header string true "Bearer token"
// @param default body models.ReqInsertUserOnlineInfo true "记录用户在线心跳包信息参数"
// @success 200 {object} models.ApiJson
// @router /teamview/insertUserOnlineInfo [post]
func (o *objApiUser) InsertUserOnlineInfo(ctx iris.Context) {
	global.GLogger.Info("InsertUserOnlineInfo called")
	var params models.ReqInsertUserOnlineInfo
	if err := utils.ARIrisReadJson(ctx, &params); err != nil {
		return
	}
	global.GLogger.Info("parameters: ", utils.StrVal(params))

	var uid = params.UId
	var optTs = params.OptTs
	//var roomId = params.RoomId

	if optTs < consts.IntZero {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamInvalid, consts.ErrParamInvalid, nil))
		return
	}
	if utils.StrEmpty(uid) {
		ctx.JSON(models.ApiJsonResp(consts.ErrCodeParamMissing, consts.ErrParamMissing, nil))
		return
	}

	code, msg, content := services.GetInsUserSvr().InsertUserOnlineInfoSvr(uid, optTs, "")
	ctx.JSON(models.ApiJsonResp(code, msg, content))
}
