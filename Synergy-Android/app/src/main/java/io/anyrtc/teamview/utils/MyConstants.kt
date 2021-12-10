package io.anyrtc.teamview.utils

import rxhttp.wrapper.annotation.DefaultDomain

typealias str = String
typealias int = Int
typealias bool = Boolean

object MyConstants {

    @DefaultDomain
	const val BASE_URL = ""

    //const val HTTP_TOKEN = "http_token"

	// --- Server API ---
    const val LOGIN = "signIn"
	const val ROOM_LIST = "getRoomList"
	const val SPECIAL_LIST = "getSpecialist"
	const val USER_INFO = "getUserInfo"
	const val CREATE_ROOM = "insertRoom"
	const val ONLINE_INFO = "insertUserOnlineInfo"
    const val JOIN_ROOM = "joinRoom"
    const val LEAVE_ROOM = "leaveRoom"
	//const val RTC_NOTIFY = "teamViewRtcNotify"
    //const val VOD_NOTIFY = "teamViewVodNotify"
    const val VOD = "startVod"


	// --- SP KEY ---
	const val UID = "uid"
	const val APP_ID = "appId"
	const val USER_NAME = "userName"
    const val WORK_NAME = "workName"
	const val PRO_MODE = "pro_mode"
	// ---  END  ---

    data class Result<T>(
        val code: int,
        val msg: str,
        val data: T,
    )

    data class Response<T>(
        val errMsg: str,
        val data: T?,
    )

    data class LoginResult(
        val appId: str,
        val rtmToken: str,
        val uid: str,
        val userName: str,
        val userTs: int,
        val userType: int,
        val workName: str,
    )

	data class RoomListResult(
		val list: List<RoomListMenu>,
		val totalNum: int,
	)
	data class RoomListMenu(
		val roomHostId: str,
		val roomId: str,
		val roomName: str,
		val roomState: int,
		val roomTs: int,
	)

	data class SpecialListResult(
		val list: List<SpecialMenu>,
		val totalNum: int,
	)
	data class SpecialMenu(
		val uid: str,
		val userName: str,
		val userState: int,
		val userTs: int,
		val userType: int,
		val workName: str,
		val roomId: str,
	)

	data class UserInfoResult(
        val uid: str,
		val userName: str,
		val userTs: int,
		val userType: int,
		val workName: str,
	)

	data class CreateRoomResult(
        val roomHostId: str,
		val roomId: str,
		val roomName: str,
		val roomState: int,
        val roomTs: int,
		val rtcToken: str,
	) {
		override fun toString(): String {
			return "CreateRoomResult(roomHostId='$roomHostId', roomId='$roomId', roomName='$roomName', roomState=$roomState, roomTs=$roomTs, rtcToken='$rtcToken')"
		}
	}

	data class JoinRoomResult(
		val roomInfo: JoinRoomMenu,
        val rtcToken: str,
	)
	data class JoinRoomMenu(
		val roomHostId: str,
		val roomId: str,
		val roomName: str,
		val roomState: int,
		val roomTs: int,
	)
}
