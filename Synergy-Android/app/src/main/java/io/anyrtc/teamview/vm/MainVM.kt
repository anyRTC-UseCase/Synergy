package io.anyrtc.teamview.vm

import android.content.Context
import android.net.ConnectivityManager
import android.net.NetworkCapabilities
import android.os.Build
import android.util.Log
import android.view.TextureView
import androidx.lifecycle.MutableLiveData
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import io.anyrtc.teamview.App
import io.anyrtc.teamview.BuildConfig
import io.anyrtc.teamview.utils.*
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.GlobalScope
import kotlinx.coroutines.launch
import org.ar.rtm.ErrorInfo
import org.ar.rtm.LocalInvitation
import org.json.JSONObject
import rxhttp.*
import rxhttp.wrapper.param.RxHttp
import rxhttp.wrapper.param.RxHttpJsonParam
import java.io.*
import java.net.ConnectException
import java.net.UnknownHostException

class MainVM : ViewModel() {

    val loginResponse = MutableLiveData<MyConstants.Response<MyConstants.LoginResult>>()
    val roomtListResponse = MutableLiveData<MyConstants.Response<MyConstants.RoomListResult>>()
    val specialListResponse = MutableLiveData<MyConstants.Response<MyConstants.SpecialListResult>>()
	val filterSpecialListResponse = MutableLiveData<MyConstants.Response<MyConstants.SpecialListResult>>()

    val userInfoResponse = MutableLiveData<MyConstants.Response<MyConstants.UserInfoResult>>()
    val createRoomResponse = MutableLiveData<MyConstants.Response<MyConstants.CreateRoomResult>>()
    val joinRoomResponse = MutableLiveData<MyConstants.Response<MyConstants.JoinRoomResult>>()

    fun login(userName: str, userType: int, workName: str) {
        val uid = SpUtil.get().getString(MyConstants.UID, (0 until 9).fold("") { acc, _ ->
            "$acc${(0..9).random()}"
        })
        viewModelScope.launch {
            try {
                val response = RxHttp.postJson(MyConstants.LOGIN)
                    .add("uid", uid)
                    .add("userName", userName)
                    .add("userType", userType)
                    .add("workName", workName).awaitOkResponse()
                val responseHeaders = response.headers
                /*for (header in responseHeaders) {
                    val key = header.first
                    val value = header.second
                }*/
                val authorization = responseHeaders["Artc-Token"] ?: throw UnknownHostException()
                val json = response.body?.string() ?: throw NullPointerException()

                val jsonObj = JSONObject(json)
                val dataObj = jsonObj.getJSONObject("data")
                val msg = jsonObj.getString("msg")
                val code = jsonObj.getInt("code")

                if (code == 0) {
                    val appId = dataObj.getString("appId")
                    val rtmToken = dataObj.getString("rtmToken")
                    val sUid = dataObj.getString("uid")
                    val sUserName = dataObj.getString("userName")
                    val userTs = dataObj.getInt("userTs")
                    val sWorkName = dataObj.getString("workName")
                    val data = MyConstants.LoginResult(
                        appId, rtmToken, sUid, sUserName, userTs, userType, sWorkName
                    )

                    App.token = authorization
                    SpUtil.edit {
                        it.putString(MyConstants.APP_ID, appId)
                        //it.putString(MyConstants.HTTP_TOKEN, authorization)
                        it.putString(MyConstants.UID, sUid)
                        it.putString(MyConstants.USER_NAME, sUserName)
                        it.putString(MyConstants.WORK_NAME, sWorkName)
                    }

					loginRtm(rtmToken, sUid) { failed, err ->
						if (failed)
					        loginResponse.postValue(MyConstants.Response("服务异常，请稍后重试", null))
						else {
							loginResponse.postValue(MyConstants.Response(msg, data))
						}
					}
                    return@launch
                }
                loginResponse.value = MyConstants.Response("服务异常，请稍后重试", null)
            } catch (e: Exception) {
				if (e is ConnectException) {
					loginResponse.value = MyConstants.Response("网络出错，请稍后重试", null)
					return@launch
				}
                val tipsStr = checkNetworkConn()
                loginResponse.value = MyConstants.Response(tipsStr, null)
            }
        }
    }

    private fun checkNetworkConn(): String {
        val cm = App.app.applicationContext.getSystemService(Context.CONNECTIVITY_SERVICE) as ConnectivityManager
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            val activeNetwork = cm.activeNetwork ?: return "网络出错，请稍后再试"
            val capabilityNetwork = cm.getNetworkCapabilities(activeNetwork) ?: return "网络出错，请稍后再试"
            return if (capabilityNetwork.hasTransport(NetworkCapabilities.TRANSPORT_WIFI)
                || capabilityNetwork.hasTransport(NetworkCapabilities.TRANSPORT_CELLULAR)
                || capabilityNetwork.hasTransport(NetworkCapabilities.TRANSPORT_ETHERNET)
                || capabilityNetwork.hasTransport(NetworkCapabilities.TRANSPORT_BLUETOOTH)
            ) "请求超时，请稍后重试" else "网络出错，请稍后重试"
        } else {
            return if (cm.activeNetworkInfo?.isConnected == true) "请求超时，请稍后重试" else "网络出错，请稍后重试"
        }
    }

    fun getRoomList(pageNum: int, pageSize: int, roomState: int) {
        reqPost(
            RxHttp.postJson(MyConstants.ROOM_LIST)
                .add("pageNum", pageNum)
                .add("pageSize", pageSize)
                .add("roomState", roomState), roomtListResponse
        )
    }

    fun getSpecialList(pageNum: int, pageSize: int) = reqPost(
        RxHttp.postJson(MyConstants.SPECIAL_LIST)
            .add("pageNum", pageNum)
            .add("pageSize", pageSize), specialListResponse
    )

	fun filterSpecialList(listSize: int) = reqPost(
		RxHttp.postJson(MyConstants.SPECIAL_LIST)
            .add("pageNum", 1)
            .add("pageSize", listSize), filterSpecialListResponse
	)

    fun getUserInfo(uid: str) = reqPost(
        RxHttp.postJson(MyConstants.USER_INFO).add("uid", uid), userInfoResponse
    )

    fun createRoom(withJoin: bool) {
        Log.e("::", "---------------------- createRoom POST ----------------------")
        reqPost(
            RxHttp.postJson(MyConstants.CREATE_ROOM).add("isJoin", if (withJoin) 1 else 2),
            createRoomResponse
        )
    }

    fun insertUserOnlineInfo(optTs: int, uid: str) = reqPost<Any>(
        RxHttp.postJson(MyConstants.ONLINE_INFO).add("optTs", optTs).add("uid", uid), null
    )

    fun joinRoom(roomId: str, userRole: int) = reqPost(
        RxHttp.postJson(MyConstants.JOIN_ROOM).add("roomId", roomId).add("userRole", userRole),
        joinRoomResponse
    )

    fun leaveRoom(roomId: str, uid: str) = reqPost<Any>(
        RxHttp.postJson(MyConstants.LEAVE_ROOM).add("roomId", roomId).add("uid", uid),
        null,
        scope = GlobalScope
    )

    fun startVod() = viewModelScope.launch {
        RxHttp.get(MyConstants.VOD)
    }

    fun sdkJoinRoom(
        context: Context,
        rtcToken: str,
        roomId: str,
        userId: str,
        joinCallback: (failed: bool) -> Unit
    ) {
        RtcManager.INSTANCE.let {
            it.initRtc(context)
            it.joinChannel(
                rtcToken = rtcToken, roomId = roomId, userId = userId, joinCallback = joinCallback
            )
            it.enableVideo()
        }
    }

    fun sendCall(
        uid: str,
        roomId: str,
        roomName: str,
        roomTs: int,
        nickname: str,
        sent: ((success: bool, err: ErrorInfo?, localInvitation: LocalInvitation?) -> Unit) = { _, _, _ -> }
    ) {
        val content = "{\"roomId\": $roomId, \"roomName\": \"$roomName\", \"roomTs\": $roomTs, \"userName\": \"$nickname\"}"
        RtcManager.INSTANCE.sendCall(uid, content, sent)
    }

    fun cancelCall(
        localInvitation: LocalInvitation,
        callback: ((success: bool, err: ErrorInfo?) -> Unit) = { _, _ -> }
    ) {
        RtcManager.INSTANCE.cancelCall(localInvitation, callback)
    }

    fun switchCamera() = RtcManager.INSTANCE.switchCamera()

    private fun loginRtm(
        token: str, uid: str,
        loginCallback: (failed: Boolean, loginFailed: ErrorInfo?) -> Unit
    ) = RtcManager.INSTANCE.run {
        initRtm()
        loginRtm(token, uid, loginCallback)
    }

	fun addRtcHandler(handler: RtcManager.RtcHandler) {
		RtcManager.INSTANCE.addRtcHandler(handler)
	}

	fun removeRtcHandler(handler: RtcManager.RtcHandler) {
		RtcManager.INSTANCE.removeRtcHandler(handler)
	}

    private fun setUpTextureView(textureView: TextureView, uid: String) {
        RtcManager.INSTANCE.changeRoleToBroadcaster()
        RtcManager.INSTANCE.enableAudio(true)
        RtcManager.INSTANCE.setUpTextureView(true, textureView, uid)
    }

    fun createRendererView(context: Context, uid: str): TextureView {
        val textureView = RtcManager.INSTANCE.createRendererView(context)
        textureView.post {
            setUpTextureView(textureView, uid)
        }
        return textureView
    }

    private inline fun <reified T> reqPost(
        http: RxHttpJsonParam,
        callback: MutableLiveData<MyConstants.Response<T>>?,
        scope: CoroutineScope = viewModelScope,
        crossinline success: ((T) -> Unit) = {}
    ) = scope.launch {
        val iAwait = http.toClass<MyConstants.Result<T>>()
        try {
            val ret = iAwait.await()
            if (ret.code == 0) {
                callback?.value = MyConstants.Response(ret.msg, ret.data)
                success.invoke(ret.data)
                return@launch
            }
            //callback?.value = MyConstants.Response(ret.msg, null)
            callback?.value = MyConstants.Response("服务异常，请稍后重试", null)
        } catch (e: Exception) {
            Log.e("::", "Req ERROR: $e")
			if (e is ConnectException) {
				callback?.value = MyConstants.Response("网络出错，请稍后重试", null)
				return@launch
			}
            //callback?.value = MyConstants.Response(it.toString(), null)
            val tipsStr = checkNetworkConn()
            callback?.value = MyConstants.Response(tipsStr, null)
            if (BuildConfig.DEBUG) Log.e(":: Req Error", e.toString())
        }
        /*http.toClass<MyConstants.Result<T>>().awaitResult { ret ->
        }.onFailure {
        }*/
    }

    fun releaseRtc() {
        RtcManager.INSTANCE.clearSelf()
    }

	fun logoutRtm() {
		RtcManager.INSTANCE.logoutRtm()
	}
}
