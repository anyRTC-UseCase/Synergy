// file: services/teamview_svr.go

package services

import (
	rtcTokenBuilder "ARTeamViewService/Tools/DynamicKey/ARDynamicKey/RtcTokenBuilder"
	rtmtokenbuilder "ARTeamViewService/Tools/DynamicKey/ARDynamicKey/RtmTokenBuilder"
	"ARTeamViewService/consts"
	"ARTeamViewService/database"
	"ARTeamViewService/global"
	"ARTeamViewService/models"
	"ARTeamViewService/utils"
	"ARTeamViewService/web/middleware"
	"github.com/kataras/iris/v12"
	json "github.com/wulinlw/jsonAmbiguity"
	"path"
	"sync"
	"time"
)

type UserSvr interface {
}

// NewUserSvr returns the default user svr.
func NewUserSvr() UserSvr {
	return &userSvr{}
}

var insUserSvr *userSvr
var onceUserSvr sync.Once

type userSvr struct {
}

func GetInsUserSvr() *userSvr {
	onceUserSvr.Do(func() {
		insUserSvr = &userSvr{}
	})
	return insUserSvr
}

/**
 *　用户登录
 */
func (u *userSvr) SignInSvr(userName, uid, workName string, userType int, ctx iris.Context) (int, string, interface{}) {
	global.GLogger.Info("SignInSvr called")

	//1.通过userName,userType,workName获取用户信息
	code, msg, arrUserInfo := u.GetUserInfoByUserInfo(uid, userName, userType, workName)

	if code == consts.ErrCodeOk {
		var userTmp models.UserInfo
		//2.有则返回,没有则插入在返回
		if len(arrUserInfo) == 0 {
			//3.入库
			affectRows, err := database.InsertUserInfo(userName, uid, workName, userType)
			if err != nil {
				global.GLogger.Error(err)
				return consts.ErrCodeDbError, consts.ErrDbError, nil
			}
			if affectRows > 0 {
				userTmp = models.UserInfo{
					UId:      uid,
					UserName: userName,
					WorkName: workName,
					UserType: userType,
					UserTs:   int(time.Now().Unix()),
				}
			} else {
				return consts.ErrCodeFailed, consts.ErrFailed, nil
			}
		} else {
			userTmp = arrUserInfo[0]
			if userTmp.UserType != userType {
				return consts.ErrCodeUserTypeDifferent, consts.ErrUserTypeDifferent, nil
			}
		}

		//4.生成userToken
		userToken, err := middleware.GenerateToken(userTmp.UId)
		if err != nil {
			global.GLogger.Error(err)
		}

		//将token存入redis
		key := utils.RedisUserKey(userTmp.UId)
		global.GRedisClient.Set(global.Context, key, userToken, global.GConfig.TokenExpire*time.Second)

		var user = models.RespSignIn{
			UId:      userTmp.UId,
			UserName: userTmp.UserName,
			WorkName: userTmp.WorkName,
			UserType: userTmp.UserType,
			UserTs:   userTmp.UserTs,
		}

		//5.生成rtmToken
		code, msg, rtmToken := u.GetRtmToken(utils.FormatNowUnix(), userTmp.UId)
		if code == consts.ErrCodeOk {
			user.RtmToken = rtmToken
		} else {
			return code, msg, nil
		}

		user.AppId = global.GConfig.AppId

		ctx.ResponseWriter().Header().Set(consts.RTCTokenHeader, userToken)
		ctx.ResponseWriter().Header().Set("Access-Control-Expose-Headers", consts.RTCTokenHeader)
		return consts.ErrCodeOk, consts.ErrOk, user
	} else {
		return code, msg, nil
	}
}

/**
 *　生成rtcToken
 */
func (u *userSvr) GetRtcToken(currentTs int64, channelName, userAccount string) (int, string, string) {

	//1.2 生成rtctoken
	var appId = global.GConfig.AppId
	var appToken = global.GConfig.AppToken

	//从配置文件获取token有效时间
	rtcTokenExpire := global.GConfig.RtcTokenExpire
	//计算超时时间戳
	expireTimeInSeconds := uint32(rtcTokenExpire)
	currentTimestamp := uint32(currentTs)
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	//生成AccessToken
	rtcToken, err := rtcTokenBuilder.BuildTokenWithUserAccount(appId, appToken, channelName, userAccount, rtcTokenBuilder.RoleAttendee, expireTimestamp)

	if err != nil {
		return consts.ErrCodeBuildTokenErr, consts.ErrBuildTokenErr, ""
	}

	global.GLogger.Info("rtcToken: ", rtcToken, ",expireTimestamp: ", expireTimestamp)
	return consts.ErrCodeOk, consts.ErrOk, rtcToken
}

/**
 *　生成rtmToken
 */
func (u *userSvr) GetRtmToken(currentTs int64, userAccount string) (int, string, string) {

	//1.2 生成rtctoken
	var appId = global.GConfig.AppId
	var appToken = global.GConfig.AppToken
	//从配置文件获取rtmToken有效时间
	interval := global.GConfig.RtmTokenExpire
	//计算超时时间戳
	expireTimeInSeconds := uint32(interval)
	currentTimestamp := uint32(currentTs)
	expireTimestamp := currentTimestamp + expireTimeInSeconds

	//生成rtmtoken
	rtmToken, err := rtmtokenbuilder.BuildToken(appId, appToken, userAccount, rtmtokenbuilder.RoleRtmUser, expireTimestamp)
	if err != nil {
		return consts.ErrCodeBuildTokenErr, consts.ErrBuildTokenErr, ""
	}
	global.GLogger.Info("rtmToken: ", rtmToken)
	return consts.ErrCodeOk, consts.ErrOk, rtmToken
}

/**
 *　消息回调
 */
func (u *userSvr) TeamViewRtcNotifySvr(roomId, uid string, eventType int, byteParam string) (int, string, interface{}) {
	global.GLogger.Info("===>TeamViewRtcNotifySvr called")
	//频道销毁  ARtcChanDestroy = 102
	//直播场景下主播加入频道  ARtcBroadcasterJoin = 103
	//直播场景下主播离开频道  ARtcBroadCasterLeave = 104
	//直播场景下观众进入频道  ARtcAudienceJoin = 105
	//直播场景下观众离开频道  ARtcAudienceLeave = 106

	if eventType == consts.ARtcChanDestroy {
		//频道销毁  ARtcChanDestroy = 102
		global.GLogger.Info("===>rtcNotify eventType: ", "ARtcChanDestroy")

		//更新所有用户离开时间
		affectRows, err := database.UpdateAllUserLeaveTs(roomId)
		if err != nil {
			global.GLogger.Error(err)
		}
		if affectRows < 0 {
			global.GLogger.Error("UpdateAllUserLeaveTs err")
		}
	} else if eventType == consts.ARtcBroadCasterLeave {
		//直播场景下主播离开频道  ARtcBroadCasterLeave = 104
		global.GLogger.Info("===>rtcNotify eventType: ", "ARtcBroadCasterLeave")

		//判断key是否存在
		exists := utils.QueryKeyIsNotExists(roomId)
		if exists {
			uidNums := utils.ZRemRdsRoomUId(roomId, uid)
			if uidNums == 0 {
				//房间状态(1:结束,2:进行中,3:转码中)
				//房间改为转码
				code, msg, _ := u.UpdateRoomStateMixed(roomId)
				if code != consts.ErrCodeOk {
					global.GLogger.Error(msg)
				}
				//4.停止录像
				u.StopVodSvr(roomId)
			}
		}
	}
	//else if eventType == consts.ARtcBroadCasterLeave {
	//	//直播场景下主播离开频道  ARtcBroadCasterLeave = 104
	//	global.GLogger.Info("----------->notify byteParam: ", byteParam)
	//	//先查此时该用户是否在房间,不在房间将leave_time=0全部更新
	//	//在的话将除最后一条的leave_time=0全部更新
	//	//先查用户是否在线
	//	totalNum, err := database.QueryUserOnlineNum(uid)
	//	if err != nil {
	//		global.GLogger.Error(err)
	//	}
	//	if totalNum > 0 {
	//		//更新除最后一条的leave_time=0的离开时间
	//		affectRows, err := database.UpdateUserLeaveTsExceptLast(roomId, uid)
	//		if err != nil {
	//			global.GLogger.Error(err)
	//			return consts.ErrCodeDbError, consts.ErrDbError, nil
	//		}
	//		if affectRows < 0 {
	//			return consts.ErrCodeFailed, consts.ErrFailed, nil
	//		}
	//	} else {
	//		affectRows, err := database.UpdateUserLeaveTs(roomId, uid)
	//		if err != nil {
	//			global.GLogger.Error(err)
	//			return consts.ErrCodeDbError, consts.ErrDbError, nil
	//		}
	//		if affectRows < 0 {
	//			return consts.ErrCodeFailed, consts.ErrFailed, nil
	//		}
	//	}
	//}

	return consts.ErrCodeOk, consts.ErrOk, nil
}

/**
 *　录像回调
 */
func (u *userSvr) TeamViewVodNotifySvr(roomId, uid string, eventType int, byteParam string, sendTs int64) (int, string, interface{}) {
	global.GLogger.Info("===>TeamViewVodNotifySvr called")
	// 录制文件已上传至指定的第三方云存储
	//ARtcVodUploaded = 31
	// 录制服务已启动
	//ARtcVodRecorderStarted = 40
	// 录制组件已退出
	//ARtcVodRecorderLeave = 41
	if eventType == consts.ARtcVodUploaded {
		global.GLogger.Info("===>vodNotify eventType: ", "ARtcVodUploaded")
		//房间状态(1:结束,2:进行中,3:转码中)
		//房间改为结束
		code, msg, _ := u.UpdateRoomStateClosed(roomId)
		if code != consts.ErrCodeOk {
			global.GLogger.Error(msg)
		}
	} else if eventType == consts.ARtcVodRecorderStarted {
		global.GLogger.Info("===>vodNotify eventType: ", "ARtcVodRecorderStarted")
		//更新房间录像开始时间
		affectRows, err := database.UpdateRoomStarTs(roomId, sendTs)
		if err != nil {
			global.GLogger.Error(err)
			return consts.ErrCodeDbError, consts.ErrDbError, nil
		}
		if affectRows < 0 {
			return consts.ErrCodeFailed, consts.ErrFailed, nil
		}
	} else if eventType == consts.ARtcVodRecorderLeave {
		global.GLogger.Info("===>vodNotify eventType: ", "ARtcVodRecorderLeave")

		//更新房间录像结束时间
		affectRows, err := database.UpdateRoomStopTs(roomId, sendTs)
		if err != nil {
			global.GLogger.Error(err)
		}
		if affectRows < 0 {
			global.GLogger.Error("UpdateRoomStopTs err")
		}
	}

	return consts.ErrCodeOk, consts.ErrOk, nil
}

/**
 *　通过uid获取用户信息
 */
func (u *userSvr) GetUserInfoByUId(uid string) (int, string, []models.UserInfo) {

	//查询用户信息
	userList, err := database.QueryUserInfoByUId(uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	var arrUserInfo []models.UserInfo
	if code, msg, _ := utils.ARParseJson(userList, &arrUserInfo); code != consts.ErrCodeOk {
		return code, msg, []models.UserInfo{}
	}

	return consts.ErrCodeOk, consts.ErrOk, arrUserInfo
}

/**
 *　通过uid或userName,userType,workName获取用户信息
 */
func (u *userSvr) GetUserInfoByUserInfo(uid, userName string, userType int, workName string) (int, string, []models.UserInfo) {

	//查询用户信息
	userList, err := database.QueryUserInfoByUserInfo(uid, userName, userType, workName)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	var arrUserInfo []models.UserInfo
	if code, msg, _ := utils.ARParseJson(userList, &arrUserInfo); code != consts.ErrCodeOk {
		return code, msg, []models.UserInfo{}
	}

	return consts.ErrCodeOk, consts.ErrOk, arrUserInfo
}

/**
 *　创建房间
 */
func (u *userSvr) InsertRoomSvr(uid string, isJoin int) (int, string, interface{}) {
	global.GLogger.Info("InsertRoomSvr called")

	//1.根据uid查询用户信息
	code, msg, arrUserInfo := u.GetUserInfoByUId(uid)

	if code == consts.ErrCodeOk {
		if len(arrUserInfo) > 0 {
			user := arrUserInfo[0]
			//2.生成roomid
			roomId := database.GetRoomId()
			roomName := user.UserName + consts.RoomNameSuffix

			roomState := consts.RoomStateClosed

			//是否加入房间(1:是,2:否)
			if isJoin == consts.IsJoinRoom {
				roomState = consts.RoomStateOpen
			}

			//2 入库
			affectRows, err := database.InsertRoomInfo(roomId, roomName, uid, roomState)
			if err != nil {
				global.GLogger.Error(err)
				return consts.ErrCodeDbError, consts.ErrDbError, nil
			}

			if affectRows == 0 {
				return consts.ErrCodeDbError, consts.ErrDbError, nil
			}

			var roomInfo models.RespInsertRoom
			if isJoin == consts.IsJoinRoom {
				//rtcToken
				code, msg, rtcToken := u.GetRtcToken(utils.FormatNowUnix(), roomId, uid)
				if code != consts.ErrCodeOk {
					return code, msg, nil
				}
				roomInfo.RtcToken = rtcToken
				//插入user_join
				affectRows, err := database.InsertUserJoinInfo(roomId, uid, consts.UserRoleHost)
				if err != nil {
					global.GLogger.Error(err)
					return consts.ErrCodeDbError, consts.ErrDbError, nil
				}
				if affectRows == 0 {
					return consts.ErrCodeFailed, consts.ErrFailed, nil
				}
				//创建房间以主播身份进入频道调用录像
				go u.StartVod(roomId, uid)
			}
			roomInfo.RoomId = roomId
			roomInfo.RoomHostId = uid
			roomInfo.RoomName = roomName
			roomInfo.RoomState = roomState
			roomInfo.RoomTs = int(utils.FormatNowUnix())

			return consts.ErrCodeOk, consts.ErrOk, roomInfo
		} else {
			return consts.ErrCodeUserNotExists, consts.ErrUserNotExists, nil
		}
	} else {
		return code, msg, nil
	}
}

//主播身份进入频道调用录像
func (u *userSvr) StartVod(roomId, uid string) {
	//1.主播uid存入redis中
	utils.ZAddRdsRoomUId(roomId, uid)
	//4.加入房间,录像
	//尝试3次
	var count = consts.IntZero
	var code = consts.IntZero
	var msg = consts.StrEmpty
	var vodUId = consts.ARAnyRtcStr + roomId
	var flag = true
	var vodRecord models.RespStartVodRecording
	for {
		count++
		code, msg, vodRecord = u.V2StartVodSvr(roomId, vodUId)
		if code != iris.StatusOK {
			global.GLogger.Error("StartVodSvr code: ", code, " msg: ", msg)
			flag = false
		} else {
			flag = true
			break
		}
		//录像restful api调用尝试3次,超过退出录制
		if count == consts.IntThree {
			global.GLogger.Error("vod start limit 3")
			flag = false
			return
		}
	}
	//flag=true说明录制成功
	if flag {
		//更新房间录像文件url,resourceId,sid信息
		var sid = vodRecord.Body.Sid
		var resourceId = vodRecord.Body.ResourceId

		affectRows, err := database.UpdateRoomInfo(roomId, resourceId, sid, vodUId)
		if err != nil {
			global.GLogger.Error(err)
		}
		if affectRows < 0 {
			global.GLogger.Error("update room vodinfo err")
		}
	}
}

/**
 *　获取用户信息
 */
func (u *userSvr) GetUserInfoSvr(uid string) (int, string, interface{}) {
	global.GLogger.Info("GetUserInfoSvr called")

	//查询用户信息
	userList, err := database.QueryUserInfoByUId(uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, err
	}

	var arrUserInfo []models.UserInfo
	if code, msg, data := utils.ARParseJson(userList, &arrUserInfo); code != consts.ErrCodeOk {
		return code, msg, data
	}
	if len(arrUserInfo) > 0 {
		return consts.ErrCodeOk, consts.ErrOk, arrUserInfo[0]
	} else {
		return consts.ErrCodeUserNotExists, consts.ErrUserNotExists, nil
	}
}

/**
 *　加入房间
 */
func (u *userSvr) JoinRoomSvr(roomId, uid string, userRole int) (int, string, interface{}) {
	global.GLogger.Info("JoinRoomSvr called")

	//1.获取房间信息
	code, msg, arrRoomInfo := u.GetRoomInfoById(roomId)
	if code == consts.ErrCodeOk {
		if len(arrRoomInfo) > 0 {
			roomInfo := arrRoomInfo[0]

			//2.插入user_join
			affectRows, err := database.InsertUserJoinInfo(roomId, uid, userRole)
			if err != nil {
				global.GLogger.Error(err)
				return consts.ErrCodeDbError, consts.ErrDbError, nil
			}
			if affectRows == 0 {
				return consts.ErrCodeFailed, consts.ErrFailed, nil
			}

			//3.生成rtcToken
			code, msg, rtcToken := u.GetRtcToken(utils.FormatNowUnix(), roomId, uid)
			if code != consts.ErrCodeOk {
				return code, msg, nil
			}
			var respJoinRoom models.RespJoinRoom
			respJoinRoom.RoomInfo = roomInfo
			respJoinRoom.RtcToken = rtcToken

			//4.如果是主播,uid存入redis中
			if userRole == consts.UserRoleHost {
				utils.ZAddRdsRoomUId(roomId, uid)
			}

			//5.查询房间里人数
			totalNum, err := database.QueryRoomUserNum(roomId)
			if err != nil {
				global.GLogger.Error(err)
			}
			if totalNum == consts.IntOne {
				//6.房间人数为1说明创建房间时没有加入房间,那么当第一个人加入房间时录像
				go u.StartVod(roomId, uid)
			}

			return consts.ErrCodeOk, consts.ErrOk, respJoinRoom
		} else {
			return consts.ErrCodeRoomNotExists, consts.ErrRoomNotExists, nil
		}
	} else {
		return code, msg, nil
	}
}

/**
 *　获取房间信息
 */
func (u *userSvr) GetRoomInfoById(roomId string) (int, string, []models.RoomInfo) {
	global.GLogger.Info("GetRoomInfoById called")

	roomList, err := database.QueryRoomInfoById(roomId)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	var arrRoomInfo []models.RoomInfo
	if code, msg, _ := utils.ARParseJson(roomList, &arrRoomInfo); code != consts.ErrCodeOk {
		return code, msg, []models.RoomInfo{}
	}
	return consts.ErrCodeOk, consts.ErrOk, arrRoomInfo
}

/**
 *　用户离开房间
 */
func (u *userSvr) LeaveRoomSvr(roomId, uid string) (int, string, interface{}) {
	global.GLogger.Info("LeaveRoomSvr called")

	//1.更新用户离开时间
	affectRows, err := database.UpdateUserLeaveTs(roomId, uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	if affectRows < 0 {
		return consts.ErrCodeFailed, consts.ErrFailed, nil
	}

	go func() {
		//2.判断此用户是以主播还是观众身份离开
		code, msg, arrUserJoinInfo := u.GetUserJoinInfo(roomId, uid)
		if code != consts.ErrCodeOk {
			global.GLogger.Error("u.GetUserJoinInfo err: ", msg)
		}
		if len(arrUserJoinInfo) > 0 {
			userJoinInfo := arrUserJoinInfo[0]
			if userJoinInfo.UserRole == consts.UserRoleHost {
				//判断key是否存在
				exists := utils.QueryKeyIsNotExists(roomId)
				if exists {
					//删除redis集合中的uid
					uidNums := utils.ZRemRdsRoomUId(roomId, uid)
					if uidNums == 0 {
						//房间状态(1:结束,2:进行中,3:转码中)
						//3.房间改为转码
						code, msg, _ := u.UpdateRoomStateMixed(roomId)
						if code != consts.ErrCodeOk {
							global.GLogger.Error(msg)
						}
						//4.停止录像
						u.StopVodSvr(roomId)
					}
				}
			}
		}
	}()
	return consts.ErrCodeOk, consts.ErrOk, nil
}

/**
 *　获取房间列表
 */
func (u *userSvr) GetRoomListSvr(roomState, pageNum, pageSize int) (int, string, interface{}) {
	global.GLogger.Info("GetRoomListSvr called")
	var respData interface{}
	//房间状态(1:结束,2:进行中)
	if roomState == consts.RoomStateClosed {
		//结束房间列表
		roomList, err := database.QueryFinishedRoomList(pageNum, pageSize)
		if err != nil {
			global.GLogger.Error(err)
			return consts.ErrCodeDbError, consts.ErrDbError, err
		}

		var arrFinishedRoomInfo []models.FinishedRoomInfo
		if code, msg, _ := utils.ARParseJson(roomList, &arrFinishedRoomInfo); code != consts.ErrCodeOk {
			return code, msg, []models.RoomInfo{}
		}
		//2.获取房间总数
		totalNum, err := database.QueryRoomNum(roomState)
		if err != nil {
			global.GLogger.Error(err)
			return consts.ErrCodeDbError, consts.ErrDbError, err
		}
		var respFinishedRoomList models.RespFinishedRoomList
		respFinishedRoomList.List = arrFinishedRoomInfo
		respFinishedRoomList.TotalNum = totalNum

		respData = respFinishedRoomList
	} else if roomState == consts.RoomStateOpen {
		//进行中房间列表
		roomList, err := database.QueryOnGoingRoomList(pageNum, pageSize, roomState)
		if err != nil {
			global.GLogger.Error(err)
			return consts.ErrCodeDbError, consts.ErrDbError, err
		}

		var arrOnGoingRoomInfo []models.OnGoingRoomInfo
		if code, msg, _ := utils.ARParseJson(roomList, &arrOnGoingRoomInfo); code != consts.ErrCodeOk {
			return code, msg, []models.RoomInfo{}
		}
		//2.获取房间总数
		totalNum, err := database.QueryRoomNum(roomState)
		if err != nil {
			global.GLogger.Error(err)
			return consts.ErrCodeDbError, consts.ErrDbError, err
		}
		var respOngoingRoomList models.RespOngoingRoomList
		respOngoingRoomList.List = arrOnGoingRoomInfo
		respOngoingRoomList.TotalNum = totalNum

		respData = respOngoingRoomList
	}
	return consts.ErrCodeOk, consts.ErrOk, respData
}

/**
 *　获取专家列表
 */
func (u *userSvr) GetSpecialistSvr(pageNum, pageSize int, uid string) (int, string, interface{}) {
	global.GLogger.Info("GetSpecialistSvr called")

	//1.获取专家列表
	userList, err := database.QuerySpecialist(pageNum, pageSize, uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, err
	}

	var arrSpecialistUserInfo []models.SpecialistUserInfo
	if code, msg, data := utils.ARParseJson(userList, &arrSpecialistUserInfo); code != consts.ErrCodeOk {
		return code, msg, data
	}
	//2.获取专家总数
	totalNum, err := database.QuerySpecialistNum(uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, err
	}
	var respGetSpecialist models.RespGetSpecialist
	respGetSpecialist.List = arrSpecialistUserInfo
	respGetSpecialist.TotalNum = totalNum
	return consts.ErrCodeOk, consts.ErrOk, respGetSpecialist

}

/**
 *　记录用户在线心跳包信息
 */
func (u *userSvr) InsertUserOnlineInfoSvr(uid string, optTs int, roomId string) (int, string, interface{}) {
	global.GLogger.Info("InsertUserOnlineInfoSvr called")
	//2 入库
	affectRows, err := database.InsertUserOnlineInfo(uid, optTs, roomId)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	if affectRows == 0 {
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	return consts.ErrCodeOk, consts.ErrOk, nil
}

/**
 *　调用resrful api开始录制
 */
func (u *userSvr) V2StartVodSvr(cname, uid string) (int, string, models.RespStartVodRecording) {
	global.GLogger.Info("V2StartVodSvr called")
	var respStartVodRecording models.RespStartVodRecording
	//1.获取resourceId
	acquireParam := models.ReqGetVodResourceId{
		Cname: cname,
		UID:   uid,
		ClientRequest: models.AcquireClientRequest{
			ResourceExpiredHour: consts.DftResExpire,
		},
	}

	byteAcquireParam, jsonErr := json.Marshal(acquireParam)
	if jsonErr != nil {
		global.GLogger.Error(jsonErr)
		return consts.ErrCodeInternal, consts.ErrJsonMarshal, respStartVodRecording
	}

	//http基本认证
	base64Credentials := utils.Base64Credentials()
	appId := global.GConfig.AppId
	if utils.StrEmpty(appId) {
		global.GLogger.Error("appId len is 0")
		return consts.ErrCodeInternal, "appId len is 0", respStartVodRecording
	}

	//1.获取resourceId
	//v1/apps/<appid>/cloud_recording/acquire
	//url,contentType
	acquireUrl := global.GConfig.HttpVodPrefix + "/v1/apps/" + appId + "/cloud_recording/acquire"

	//发送post请求
	code, body := utils.PostStatusRequest(acquireUrl, consts.DftVodContentType, string(byteAcquireParam), base64Credentials)
	global.GLogger.Info("acquire code: ", code, " body: ", body)

	if code != iris.StatusOK {
		global.GLogger.Error("vod acquire error")
		return code, "vod acquire error", respStartVodRecording
	}

	//解析数据得到resourceId
	var respGetVodResourceId models.RespGetVodResourceId
	utils.ARUnmarshalJson(body, &respGetVodResourceId)

	global.GLogger.Info("acquire resourceId: ", respGetVodResourceId.Body.ResourceId)
	resourceId := respGetVodResourceId.Body.ResourceId
	if utils.StrEmpty(resourceId) {
		global.GLogger.Error("vod acquire resourceId error")
		return code, "vod acquire resourceId error", respStartVodRecording
	}

	time.Sleep(consts.Dft5Sec)
	//2.开始录制
	//rtcToken
	code, msg, rtcToken := u.GetRtcToken(utils.FormatNowUnix(), cname, uid)
	if code != consts.ErrCodeOk {
		global.GLogger.Error(msg)
		return code, msg, respStartVodRecording
	}
	startParam := models.ReqStartVodRecording{
		Cname: cname,
		UID:   uid,
		ClientRequest: models.StartClientRequest{
			Token: rtcToken,
			RecordingConfig: models.RecordingConfig{
				ChannelType: consts.DftVodChannelType,
				MaxIdleTime: consts.DftVodMaxIdleTime,
				StreamTypes: consts.DftVodStreamType,
				TranscodingConfig: models.TranscodingConfig{
					Height:           consts.DftVodTransHeight,
					Width:            consts.DftVodTransWidth,
					Bitrate:          consts.DftVodTransBitrate,
					Fps:              consts.DftVodTransFps,
					MixedVideoLayout: consts.DftVodTransVideoLayout,
				},
				SubscribeUidGroup: consts.DftSubscribeUidGroup,
			},
			RecordingFileConfig: models.RecordingFileConfig{
				AvFileType: []string{consts.VodStreamHls, consts.VodStreamMp4},
			},
			StorageConfig: models.StorageConfig{
				AccessKey:      global.GConfig.StorageConfig.AccessKey,
				Region:         global.GConfig.StorageConfig.Region,
				Bucket:         global.GConfig.StorageConfig.Bucket,
				SecretKey:      global.GConfig.StorageConfig.SecretKey,
				Vendor:         global.GConfig.StorageConfig.Vendor,
				FileNamePrefix: global.GConfig.StorageConfig.FileNamePrefix,
			},
		},
	}

	byteStartParam, jsonErr := json.Marshal(startParam)
	if jsonErr != nil {
		global.GLogger.Error(jsonErr)
		return consts.ErrCodeInternal, consts.ErrJsonMarshal, respStartVodRecording
	}

	//1.开始录制
	///v1/apps/<yourappid>/cloud_recording/resourceid/<resourceid>/mode/<mode>/start
	//url,contentType
	startUrl := global.GConfig.HttpVodPrefix + "/v1/apps/" + appId + "/cloud_recording/resourceid/" + resourceId + "/mode/" + consts.DftVodModMix + "/start"

	//发送post请求
	code, content := utils.PostStatusRequest(startUrl, consts.DftVodContentType, string(byteStartParam), base64Credentials)
	global.GLogger.Info("start code: ", code, " body: ", body)

	if code != iris.StatusOK {
		global.GLogger.Error("vod start error")
		return code, "vod start error", respStartVodRecording
	}

	//2.解析数据得到resourceId

	utils.ARUnmarshalJson(content, &respStartVodRecording)

	global.GLogger.Info("start resourceId: ", respStartVodRecording.Body.ResourceId)
	global.GLogger.Info("start sid: ", respStartVodRecording.Body.Sid)

	return code, consts.ErrOk, respStartVodRecording
}

//每天0点50分执行任务删掉前一天的心跳包
func (u *userSvr) DayDeleteUserOnlineInfo() (int, string, interface{}) {
	global.GLogger.Info("DayDeleteUserOnlineInfo called")
	ts := utils.FormatNowUnix() - int64(consts.Dft24Hour/time.Second)
	cnt, err := database.DayDeleteUserOnlineInfo(ts)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	if cnt < 0 {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	return consts.ErrCodeOk, consts.ErrOk, nil
}

//更新房间状态
func (u *userSvr) UpdateRoomState(roomId string, roomState int) (int, string, interface{}) {
	global.GLogger.Info("UpdateRoomState ts: ", utils.FormatNowUnix(), ",roomId: ", roomId, ",roomState: ", roomState)

	affectRows, err := database.UpdateRoomState(roomId, roomState)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	if affectRows < 0 {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	return consts.ErrCodeOk, consts.ErrOk, nil
}

//更新房间状态为转码
func (u *userSvr) UpdateRoomStateMixed(roomId string) (int, string, interface{}) {

	affectRows, err := database.UpdateRoomStateMixed(roomId)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	if affectRows < 0 {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	global.GLogger.Info("UpdateRoomStateMixed ts: ", utils.FormatNowUnix(), ",affectRows: ", affectRows)
	return consts.ErrCodeOk, consts.ErrOk, nil
}

//更新房间状态为关闭
func (u *userSvr) UpdateRoomStateClosed(roomId string) (int, string, interface{}) {

	affectRows, err := database.UpdateRoomStateClosed(roomId)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	if affectRows < 0 {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}
	global.GLogger.Info("UpdateRoomStateClosed ts", utils.FormatNowUnix(), ",affectRows: ", affectRows)
	return consts.ErrCodeOk, consts.ErrOk, nil
}

/**
 *　根据roomId和uid查询用户最后一条加入房间的信息
 */
func (u *userSvr) GetUserJoinInfo(roomId, uid string) (int, string, []models.UserJoinInfo) {
	global.GLogger.Info("GetUserJoinInfo called")

	userList, err := database.QueryUserJoinInfo(roomId, uid)
	if err != nil {
		global.GLogger.Error(err)
		return consts.ErrCodeDbError, consts.ErrDbError, nil
	}

	var arrUserJoinInfo []models.UserJoinInfo
	if code, msg, _ := utils.ARParseJson(userList, &arrUserJoinInfo); code != consts.ErrCodeOk {
		return code, msg, nil
	}

	return consts.ErrCodeOk, consts.ErrOk, arrUserJoinInfo
}

/**
 *　调用resrful api停止录制
 */
func (u *userSvr) StopVodSvr(cname string) {
	global.GLogger.Info("StopVodSvr called")

	//房间信息
	code, msg, arrRoomInfo := u.GetRoomInfoById(cname)
	if code == consts.ErrCodeOk {
		if len(arrRoomInfo) > 0 {
			roomInfo := arrRoomInfo[0]
			uid := roomInfo.RoomVodUId
			resourceId := roomInfo.RoomVodResourceId
			sid := roomInfo.RoomVodSId
			if utils.StrEmpty(uid) || utils.StrEmpty(resourceId) || utils.StrEmpty(sid) {
				global.GLogger.Error("StopVodSvr param nil")
				global.GLogger.Error("StopVodSvr uid: ", uid, " resourceId: ", resourceId, " sid: ", sid)
				return
			}
			//停止云端录制请求参数
			stopParam := models.ReqStopVodRecording{
				Cname: cname,
				UId:   uid,
			}

			byteStopParam, jsonErr := json.Marshal(stopParam)
			if jsonErr != nil {
				global.GLogger.Error(jsonErr)
				return
			}

			//http基本认证
			base64Credentials := utils.Base64Credentials()
			appId := global.GConfig.AppId
			if utils.StrEmpty(appId) {
				global.GLogger.Error("appId len is 0")
				return
			}

			///v1/apps/<appid>/cloud_recording/resourceid/<resourceid>/sid/<sid>/mode/<mode>/stop
			//url,contentType
			stopUrl := global.GConfig.HttpVodPrefix + "/v1/apps/" + appId + "/cloud_recording/resourceid/" + resourceId +
				"/sid/" + sid + "/mode/" + consts.DftVodModMix + "/stop"
			//发送post请求
			code, body := utils.PostStatusRequest(stopUrl, consts.DftVodContentType, string(byteStopParam), base64Credentials)
			global.GLogger.Info("stop code: ", code, " body: ", body)

			if code != iris.StatusOK {
				global.GLogger.Error("vod stop error")
				return
			}

			//解析数据
			var respStopResponse models.RespStopResponse
			utils.ARUnmarshalJson(body, &respStopResponse)
			global.GLogger.Info("vod stop response respStopResponse: ", respStopResponse)

			//存入录像文件地址
			fileList := respStopResponse.Body.ServerResponse.FileList
			fileName := consts.StrEmpty

			if len(fileList) == 0 {
				global.GLogger.Error("vod notify fileList len 0")
				return
			}
			for i, _ := range fileList {
				if utils.Equals(path.Ext(fileList[i].FileName), consts.StrFileSuffix) {
					fileName = fileList[i].FileName
					break
				}
			}
			if utils.StrEmpty(fileName) {
				global.GLogger.Error("vod notify no mp4 file")
				return
			}
			fileUrl := global.GConfig.HttpVodFilePrefix + fileName
			global.GLogger.Info("vod fileUrl: ", fileUrl)

			affectRows, err := database.UpdateRoomVodUrl(cname, fileUrl)
			if err != nil {
				global.GLogger.Error(err)
				return
			}
			if affectRows < 0 {
				global.GLogger.Error("UpdateRoomVodUrl err")
				return
			}

		} else {
			global.GLogger.Error(consts.ErrRoomNotExists)
		}
	} else {
		global.GLogger.Error("u.GetRoomInfoById err: ", msg)
	}
}
