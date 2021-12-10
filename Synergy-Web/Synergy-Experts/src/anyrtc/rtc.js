/**
 * RTC 封装
 */
import ArRTC from "ar-rtc-sdk";
// 页面提示
import { ElMessage, ElNotification } from "element-plus";
// 公共方法
import { getServeInfo } from "./common";
// vuex
import store from "@/store";
// import Route from "@/router";
import { LogoutRTM } from "./rtm";
// 相关接口
import { joinRoom, leaveRoom } from "@/api/rtc_rtm.js";
import { getUserInfo } from "@/api/home.js";

// RTC 私有云配置
const configuration = {
  // //配置私有云网关
  //   ConfPriCloudAddr: {
  //     ServerAdd: "",
  //     Port: ,
  //     Wss: false,
  //   },
};

var Store = {
  // 视频容器id
  VideoContainerID: "",
  // 远端用户存放
  RemoteUserList: [],
  // 判断是否加入房间
  isJoinChannel: false,
  // 频道内的声音
  closeChannelAudio: false,
  // 角色
  Role: "观众",
  // RTC客户端
  rtcClient: null,
  // 本地音频
  localAudioTrack: null,
  // 本地音频发布
  localAudioTrackPublishStatus: true,
  // 记录终端用户的uid
  RecordTerminalUid: "",
  // 当前RTC状态
  connectionState: "DISCONNECTED",

  // 提示用户加入(仅提示用户加入频道以后再加入的用户)
  hintUserJoin: false,

  // 设置重连时间
  reconnectionTime: 60 * 1000,
  // 重连定时器
  reconnectionTimer: null,
  // 重连退出
  reconnectionQuit: false,

  // 终端加载时间(加入房间后一定时间内无终端)
  terminalTime: 10 * 1000,
  // 终端加载定时器
  terminalTimer: null,
};
// RTC 相关回调封装
const RTCCallback = {
  // RTC SDK 监听用户发布
  userPublished: async (user, mediaType) => {
    // 订阅用户发布的音视频
    await Store.rtcClient.subscribe(user, mediaType);
    if (mediaType === "video") {
      Store.terminalTimer && clearTimeout(Store.terminalTimer);
      await store.commit("upDataVideoLoding", false);
      await store.commit("setTerminalOffline", false);
      // 存放终端用户
      Store.RecordTerminalUid = user.uid;
      // 渲染用户发布的视频
      user.videoTrack &&
        user.videoTrack.play(Store.VideoContainerID, {
          fit: "contain",
        });
    } else {
      // user.audioTrack && (await user.audioTrack.play());
      if (Store.Role === "主播") {
        user.audioTrack && (await user.audioTrack.play());
      } else {
        if (Store.closeChannelAudio) {
          user.audioTrack && (await user.audioTrack.play());
        } else {
          user.audioTrack && user.audioTrack.stop();
        }
        Store.RemoteUserList.push(user);
      }
    }
  },
  // RTC SDK 监听用户取消发布
  userUnpublished(user, mediaType) {
    console.log(" RTC SDK 监听用户取消发布", user, mediaType);
    if (mediaType == "video") {
      Store.terminalTimer && clearTimeout(Store.terminalTimer);
      store.commit("setTerminalOffline", true);
    }
    Store.RemoteUserList = Store.RemoteUserList.filter((item) => {
      return item.uid != user.uid;
    });
  },
  // RTC SDK 监听用户加入频道成功
  async userJoined(user) {
    if (Store.Role === "主播" && !Store.reconnectionQuit) {
      const { code, data } = await getUserInfo({ uid: user.uid });
      if (Store.hintUserJoin && code == 0) {
        ElNotification({
          // title: "Success",
          message: data.userName + "成功进入协同",
          showClose: true,
          type: "success",
          // position: 'bottom-right',
          // duration: 2000,
        });
      }
      store.commit(
        "upDataInvitationIng",
        Object.assign(data, { userState: 1 })
      );
    }
  },
  // RTC SDK 监听用户离开频道
  async userLeft(user, reason) {
    console.log("用户离开频道", user, reason);
    if (!Store.reconnectionQuit) {
      // 终端离线
      if (Store.RecordTerminalUid == user.uid) {
        store.commit("setTerminalOffline", true);
      } else {
        if (Store.Role === "主播") {
          const { code, data } = await getUserInfo({ uid: user.uid });
          if (code == 0) {
            ElNotification({
              // title: "Warning",
              message: data.userName + "离开协同",
              type: "warning",
              showClose: true,
              // position: 'bottom-right',
              // duration: 2000,
            });
          }
          store.commit(
            "clearInvitationIng",
            Object.assign(data, { userState: 2 })
          );
        }
      }
    }
  },
  // RTC SDK 连接状态
  async connectionStateChange(curState, status, reason) {
    Store.connectionState = curState;
    console.log("RTC SDK 连接状态", curState);
    switch (curState) {
      case "RECONNECTING":
        ElMessage.warning("RTC 正在重连中");
        if (Store.Role == "主播") {
          Store.reconnectionTimer = setTimeout(() => {
            ElMessage.error("连接已断开");
            Store.reconnectionQuit = true;
            // 返回首页
            const Router = require("../router/index");
            Router.default.replace("/");
          }, Store.reconnectionTime);
        }
        break;
      case "CONNECTED":
        // ElMessage.warning("RTC 已连接");
        Store.reconnectionTimer && clearTimeout(Store.reconnectionTimer);
        if (Store.reconnectionQuit) {
          setTimeout(() => {
            // 离开房间
            LeaveRTCChannel();
          }, 200);
        }
        break;
      case "CONNECTING":
        ElMessage.warning("RTC 连接中");
        Store.reconnectionTimer && clearTimeout(Store.reconnectionTimer);
        break;
      case "DISCONNECTING":
        ElMessage.warning("离开房间中");
        Store.reconnectionTimer && clearTimeout(Store.reconnectionTimer);
        break;
      case "DISCONNECTED":
        Store.reconnectionTimer && clearTimeout(Store.reconnectionTimer);
        // ElMessage.warning("RTC 连接断开");
        // if (reason == "LEAVE") {
        // ElMessage.info("离开房间");
        // }
        if (reason == "UID_BANNED") {
          ElMessage.error("账号在别处登录");
          // 清空存储返回登录页面
          // setTimeout(async () => {
          await LeaveRTCChannel();
          await DestroyRTC();
          await LogoutRTM();

          sessionStorage.clear();
          // 返回首页
          const Router = require("../router/index");
          // Router.default.replace("/");
          Router.default.replace("/signin");
          // }, 2000);
        }
        break;
      default:
        break;
    }
  },
  // Token 即将过期
  async tokenPrivilegeWillExpire() {
    if (Store.Role == "主播") {
      ElMessage.warning("体验时间即将结束");
    } else {
      // 更新token
      const { code, data } = await joinRoom({
        roomId: Store.roomId,
        userRole: Store.Role === "主播" ? 1 : 2,
      });
      if (code !== 0) {
        ElMessage.error("加入房间失败");
        Store.isJoinChannel = false;
        store.commit("upDataIsJoinChannel", false);
      } else {
        console.log("更新 token", data.rtcToken);
        Store.rtcClient.renewToken(data.rtcToken).catch((err) => {
          console.log(
            "🚀 ~ file: rtc.js ~ line 223 ~ Store.rtcClient.renewToken ~ err",
            err
          );
        });
      }
    }
  },
  // Token 已过期
  async tokenPrivilegeDidExpire() {
    if (Store.Role == "主播") {
      ElMessage.error("体验时间已到");
      // 离开房间
      await LeaveRTCChannel();
      // 返回首页
      const Router = require("../router/index");
      Router.default.replace("/");
    }
  },
};

// RTC 初始化
export const InitRTC = async () => {
  if (Store.rtcClient) {
    return;
  }
  // RTC 版本信息
  console.log("RTC 版本", ArRTC.VERSION);
  // 获取相关信息
  Store = await Object.assign(Store, getServeInfo());
  // 客户端 配置: 通话场景、编码格式
  Store.rtcClient = await ArRTC.createClient({
    mode: "live",
    codec: "h264",
  });
  // 设置私有云
  configuration && Store.rtcClient.setParameters(configuration);

  // window.client = Store.rtcClient;
  // RTC SDK 监听用户发布
  Store.rtcClient.on("user-published", RTCCallback.userPublished);
  // RTC SDK 监听用户取消发布
  Store.rtcClient.on("user-unpublished", RTCCallback.userUnpublished);
  //  RTC SDK 监听用户加入频道成功
  Store.rtcClient.on("user-joined", RTCCallback.userJoined);
  // RTC SDK 监听用户离开频道
  Store.rtcClient.on("user-left", RTCCallback.userLeft);
  // RTC SDK 连接状态
  Store.rtcClient.on(
    "connection-state-change",
    RTCCallback.connectionStateChange
  );
  // Token 即将过期
  Store.rtcClient.on(
    "token-privilege-will-expire",
    RTCCallback.tokenPrivilegeWillExpire
  );
  // Token 已过期
  Store.rtcClient.on(
    "token-privilege-did-expire",
    RTCCallback.tokenPrivilegeDidExpire
  );
};

// 直播模式下设置角色
export const SetClientRole = (role = "主播") => {
  Store.Role = role;
  Store.rtcClient
    .setClientRole(role === "主播" ? "host" : "audience")
    .catch((err) => {
      console.log(
        "🚀 设置角色失败 ~ file: rtc.js ~ line 229 ~ SetClientRole ~ err",
        role,
        err
      );
    });
};
// 加入频道 videoId 视频容器id
export const JoinRTCChannel = async (info, videoId) => {
  if (!Store.isJoinChannel && Store.connectionState == "DISCONNECTED") {
    store.commit("upDataVideoLoding", true);
    Store = Object.assign(Store, info, {
      VideoContainerID: videoId,
      isJoinChannel: true,
    });
    const { code, data } = await joinRoom({
      roomId: Store.roomId,
      userRole: Store.Role === "主播" ? 1 : 2,
    });
    if (code !== 0) {
      ElMessage.error("加入房间失败");
      Store.isJoinChannel = false;
      store.commit("upDataIsJoinChannel", false);
      return;
    }
    Store.rtcClient
      .join(Store.appId, info.roomId, data.rtcToken, Store.uid)
      .then(() => {
        // 加入频道成功
        store.commit("upDataIsJoinChannel", true);
        // 终端离线
        Store.terminalTimer = setTimeout(() => {
          store.commit("upDataVideoLoding", false);
          store.commit("setTerminalOffline", true);
        }, Store.terminalTime);

        ElMessage.success(
          "加入房间ID为" + Store.roomId + "的" + Store.roomName + "成功"
        );
        if (Store.Role === "主播") {
          //  主播发布音频
          if (Store.localAudioTrack) {
            Store.rtcClient
              .publish(Store.localAudioTrack)
              .then(() => {
                console.log("主播发布音频");
                Store.hintUserJoin = true;
              })
              .catch((res) => {
                ElMessage.error("发布失败", JSON.stringify(res));
              });
          } else {
            ElMessage.error("无音频，无法发布");
          }
        }
      })
      .catch((err) => {
        console.log("🚀 ~ file: rtc.js ~ line 276 ~ JoinRTCChannel ~ err", err);
        // 加入频道失败
        let str = "";
        if (err.message.includes("CONNECTED")) {
          str = "正在离开房间中，无法加入，请稍后在加入";
        } else {
          str = err.message;
        }

        ElMessage.error("加入频道失败:" + str);
        Store.isJoinChannel = false;
        store.commit("upDataIsJoinChannel", false);
      });
  } else {
    ElMessage.error("操作过于频繁");
  }
};
// 采集音频
export const CollectionAudio = async () => {
  const microhones = await ArRTC.getMicrophones();
  if (microhones.length != 0) {
    Store.localAudioTrack = await ArRTC.createMicrophoneAudioTrack().catch(
      (err) => {
        console.log(err);
      }
    );
  } else {
    ElMessage.warning("SDK没有找到麦克风");
  }
};
// 关闭频道内所有声音
export const closeChannelAudio = (iswitch) => {
  Store.closeChannelAudio = iswitch;
  if (Store.RemoteUserList.length > 0) {
    Store.RemoteUserList.map((item) => {
      iswitch ? item.audioTrack.play() : item.audioTrack.stop();
    });
  }
};
// 关闭本地声音
export const CloseLocalAudio = (fase) => {
  console.log("关闭本地声音", Store.localAudioTrack);
  Store.localAudioTrackPublishStatus = fase;
  if (Store.localAudioTrack) {
    Store.localAudioTrack.setEnabled(fase);
  } else {
    ElMessage.warning("SDK没有找到麦克风，无效操作");
  }
};
// 离开频道
export const LeaveRTCChannel = async () => {
  if (Store.isJoinChannel && Store.connectionState == "CONNECTED") {
    await Store.rtcClient.leave();
  }
  store.commit("upDataIsJoinChannel", false);
  if (Store.localAudioTrack) {
    // 停止发布
    Store.localAudioTrackPublishStatus && (await Store.rtcClient.unpublish());
    // 释放音频
    await Store.localAudioTrack.close();
  }
  leaveRoom({
    roomId: Store.roomId,
    uid: Store.uid,
  });
  Store.terminalTimer && clearTimeout(Store.terminalTimer);
  Store = Object.assign(Store, {
    // 视频容器id
    VideoContainerID: "",
    // 远端用户存放
    RemoteUserList: [],
    // 判断是否加入房间
    isJoinChannel: false,
    // 频道内的声音
    closeChannelAudio: false,
    // 角色
    Role: "观众",
    // RTC客户端
    // rtcClient: null,
    // 本地音频
    localAudioTrack: null,
    // 本地音频发布
    localAudioTrackPublishStatus: true,
    // 记录终端用户的uid
    RecordTerminalUid: "",

    // 提示用户加入(仅提示用户加入频道以后再加入的用户)
    hintUserJoin: false,
    // 设置重连时间
    reconnectionTime: 60 * 1000,
    // 重连定时器
    reconnectionTimer: null,
    // 重连退出
    reconnectionQuit: false,

    // 终端加载时间(加入房间后一定时间内无终端)
    terminalTime: 10 * 1000,
    // 终端加载定时器
    terminalTimer: null,
  });
  store.commit("upDataVideoLoding", false);
  store.commit("setTerminalOffline", false);
};
// 销毁RTC
export const DestroyRTC = () => {
  Store = Object.assign(Store, {
    // 视频容器id
    VideoContainerID: "",
    // 远端用户存放
    RemoteUserList: [],
    // 判断是否加入房间
    isJoinChannel: false,
    // 频道内的声音
    closeChannelAudio: false,
    // 角色
    Role: "观众",
    // RTC客户端
    rtcClient: null,
    // 本地音频
    localAudioTrack: null,
    // 本地音频发布
    localAudioTrackPublishStatus: true,
    // 记录终端用户的uid
    RecordTerminalUid: "",
    // 当前RTC状态
    connectionState: "DISCONNECTED",

    // 重连定时器
    reconnectionTimer: null,
    // 重连退出
    reconnectionQuit: false,

    terminalTimer: null,
  });
};
