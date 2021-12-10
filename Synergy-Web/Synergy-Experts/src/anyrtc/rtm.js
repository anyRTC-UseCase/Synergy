// RTM 封装
import ArRTM from "ar-rtm-sdk";
// 公共方法
import { getServeInfo } from "./common";
// 页面提示
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
// import {useRouter} from "vue-router";
import router from "@/router/index.js";
// vuex
import store from "@/store";
// RTC
import { LeaveRTCChannel, DestroyRTC } from "./rtc.js";

// RTM 私有云配置
const configuration = {
  // //配置私有云网关
  // confPriCloudAddr: {
  //   ServerAdd: "",
  //   Port: ,
  //   Wss: false,
  // },
};

var Store = {
  // 是否登录RTM
  loginRTM: false,
  // RTM 客户端
  rtmClient: null,
  // 呼叫邀请实例是否主动取消
  initiativeCancel: false,
  // 呼叫邀请实例定时器
  localInvitationTime: null,
  // 记录呼叫邀请实例
  localInvitationLists: [],
  // 记录收到的呼叫邀请
  recordRemoteInviationLists: [],
  // 被呼叫者
  peerid: "",
};

// 清除收到的呼叫邀请
const clearRecordRemoteInviationLists = (record) => {
  if (Store.recordRemoteInviationLists.length > 0) {
    // 清除对应实例
    if (record) {
      Store.recordRemoteInviationLists =
        Store.recordRemoteInviationLists.filter((item) => {
          return item.callerId != record.callerId;
        });
      if (Store.recordRemoteInviationLists.length == 0) {
        ElMessageBox.close(false); //重要
      }
    } else {
      Store.recordRemoteInviationLists = [];
      ElMessageBox.close(false); //重要
    }
  }
};

// RTM 相关回调
const RTMCallback = {
  // 被叫：收到来自主叫的呼叫邀请
  RemoteInvitationReceived(remoteInvitation) {
    console.log("收到来自主叫的呼叫邀请", remoteInvitation);
    Store.recordRemoteInviationLists.push(remoteInvitation);
    // 解析附带信息
    const invitationContent = JSON.parse(remoteInvitation.content);
    // 弹出通知
    ElMessageBox.confirm(
      invitationContent.userName +
        " 邀请您进入 " +
        invitationContent.roomName +
        "房间 协同",
      "邀请通知",
      {
        confirmButtonText: "同意",
        cancelButtonText: "拒绝",
        distinguishCancelAndClose: true,
        showClose: false,
      }
    )
      .then(async () => {
        ElMessageBox.close(false); //重要
        // 同意进入协同
        remoteInvitation.accept();
      })
      .catch(() => {
        // 拒绝进入协同
        remoteInvitation.refuse();
      });

    // 监听接受呼叫邀请
    remoteInvitation.on("RemoteInvitationAccepted", async () => {
      console.log("监听接受呼叫邀请");
      // 跳转至协同
      await LeaveRTCChannel();

      router.replace("/about");
      setTimeout(() => {
        router.replace({
          path: "/synergy",
          query: {
            roomId: invitationContent.roomId,
            roomName: invitationContent.roomName,
            roomTs: invitationContent.roomTs,
            userName: invitationContent.roomUserName
              ? invitationContent.roomUserName
              : invitationContent.userName,
          },
        });
      }, 200);

      // 拒绝其他邀请
      if (Store.recordRemoteInviationLists.length > 1) {
        Store.recordRemoteInviationLists.map((item) => {
          if (item._callId != remoteInvitation._callId) {
            item.refuse();
          }
        });
        clearRecordRemoteInviationLists();
      }
    });
    // 监听拒绝呼叫邀请
    remoteInvitation.on("RemoteInvitationRefused", () => {
      console.log("监听拒绝呼叫邀请", invitationContent);
      ElMessage.error("您已拒绝来自" + invitationContent.userName + "协同邀请");
    });
    // 监听主叫取消呼叫邀请
    remoteInvitation.on("RemoteInvitationCanceled", (reson) => {
      console.log(
        "🚀 ~ file: rtm.js ~ line 119 ~ remoteInvitation.on ~ reson",
        reson
      );
      console.log("监听主叫取消呼叫邀请", remoteInvitation);
      // 60秒无操作自动取消
      ElMessage.error("60s未接受" + invitationContent.userName + "的邀请");
      clearRecordRemoteInviationLists(remoteInvitation);
    });
    // 监听呼叫邀请进程失败
    remoteInvitation.on("RemoteInvitationFailure", () => {
      console.log("监听呼叫邀请进程失败");
      // ElMessageBox.close(false); //重要
    });
  },
  // 收到来自对端的点对点消息
  MessageFromPeer() {},
  // 连接状态发生了改变
  async ConnectionStateChanged(status, reason) {
    // 用户在其他地方登录，当前返回登录页面
    if (status == "DISCONNECTED" && reason == "REMOTE_LOGIN") {
      Store.loginRTM = false;
      ElMessage.error("账号在别处登录");
      await LeaveRTCChannel();
      await DestroyRTC();
      await LogoutRTM();

      await sessionStorage.clear();
      // // 返回登录
      await router.replace("/signin");
    }
  },
};

// 初始化
export const InitRTM = async () => {
  if (Store.rtmClient) {
    return;
  }
  console.log("RTM 版本", ArRTM.VERSION);
  Store = await Object.assign(Store, getServeInfo());
  Store.rtmClient = await ArRTM.createInstance(Store.appId);
  // RTM 私有云
  configuration && Store.rtmClient.setParameters(configuration);
  // 登录
  Store.rtmClient
    .login({ token: Store.rtmToken, uid: Store.uid })
    .then((res) => {
      console.log("登录", res);
      Store.loginRTM = true;
    })
    .catch((err) => {
      console.log("登录失败", err);
      ElMessage.error("登录失败" + err.message);
    });
  // 监听收到来自主叫的呼叫邀请
  Store.rtmClient.on(
    "RemoteInvitationReceived",
    RTMCallback.RemoteInvitationReceived
  );
  // 监听收到来自对端的点对点消息
  Store.rtmClient.on("MessageFromPeer", RTMCallback.MessageFromPeer);
  // 通知 SDK 与 RTM 系统的连接状态发生了改变
  Store.rtmClient.on(
    "ConnectionStateChanged",
    RTMCallback.ConnectionStateChanged
  );
};

// 发起呼叫
export const SendCall = async (info) => {
  Store.peerid = info.uid;
  const localInvitation = await Store.rtmClient.createLocalInvitation(info.uid);
  // 设置呼叫发送的信息
  localInvitation.content = await JSON.stringify({
    // 发起呼叫的用户名称
    userName: Store.userName,
    // 房间名称
    roomName: info.roomName,
    // 房间ID
    roomId: info.roomId,
    roomTs: info.roomTs,
    roomUserName: info.roomUserName,
  });
  // 主叫：被叫已收到呼叫邀请
  localInvitation.on("LocalInvitationReceivedByPeer", () => {
    console.log("主叫：被叫已收到呼叫邀请");
    ElNotification({
      message: info.userName + " 已收到邀请",
      showClose: true,
      type: "info",
      // position: 'bottom-right',
    });
  });
  // 主叫：呼叫被叫端失败
  localInvitation.on("LocalInvitationFailure", (response) => {
    console.log("主叫：呼叫被叫端失败", response);
    if (response != "PEER_NO_RESPONSE") {
      ElNotification({
        message: info.userName + "呼叫失败",
        showClose: true,
        type: "error",
        // position: 'bottom-right',
      });
    }
    ClearLocalInvitation(localInvitation);
    // 可重新呼叫
    store.commit("upDataInvitationIng", Object.assign(info, { userState: 2 }));
  });
  // 主叫：呼叫邀请已取消
  localInvitation.on("LocalInvitationCanceled", () => {
    console.log("主叫：呼叫邀请已取消");
    if (Store.initiativeCancel) {
      ElNotification({
        message: info.userName + " 60s无操作自动取消邀请",
        showClose: true,
        type: "success",
        // position: 'bottom-right',
      });
    }

    ClearLocalInvitation(localInvitation);
    // 可重新呼叫
    store.commit("upDataInvitationIng", Object.assign(info, { userState: 2 }));
  });
  // 主叫：被叫已接受呼叫邀请
  localInvitation.on("LocalInvitationAccepted", (response) => {
    console.log("主叫：被叫已接受呼叫邀请", response);
    ElNotification({
      message: info.userName + " 已接受邀请",
      showClose: true,
      type: "success",
      // position: 'bottom-right',
    });
    ClearLocalInvitation(localInvitation);
    // 可重新呼叫
    store.commit("upDataInvitationIng", Object.assign(info, { userState: 1 }));
  });
  // 主叫：被叫已拒绝呼叫邀请
  localInvitation.on("LocalInvitationRefused", (response) => {
    console.log("主叫：被叫已拒绝呼叫邀请", response);
    ElNotification({
      message: info.userName + " 已拒绝邀请",
      showClose: true,
      type: "warning",
      // position: 'bottom-right',
    });
    ClearLocalInvitation(localInvitation);
    // 可重新呼叫
    store.commit("upDataInvitationIng", Object.assign(info, { userState: 2 }));
  });
  // 发送
  await localInvitation.send();
  Store.localInvitationTime = setTimeout(() => {
    Store.initiativeCancel = true;
  }, 58 * 1000);
  Store.localInvitationLists.push(localInvitation);
};
// 登出 RTM
export const LogoutRTM = async () => {
  if (Store.loginRTM) {
    console.log("登出 RTM", Store.rtmClient);
    await Store.rtmClient.logout().catch((err) => {
      console.log(
        "🚀 ~ file: rtm.js ~ line 249 ~ awaitStore.rtmClient.logout ~ err",
        err
      );
    });
  }
  Store = Object.assign(Store, {
    loginRTM: false,
    // RTM 客户端
    rtmClient: null,
    // 记录收到的呼叫邀请
    recordRemoteInviationLists: [],
    // 被呼叫者
    peerid: "",
    // 记录呼叫邀请实例
    localInvitationLists: [],
    // 呼叫邀请实例是否主动取消
    initiativeCancel: false,
    // 呼叫邀请实例定时器
    localInvitationTime: null,
  });
};
// 清除邀请实例(协同页面邀请后退出协同)
export const ClearLocalInvitation = (invitation) => {
  if (invitation) {
    Store.localInvitationLists = Store.localInvitationLists.filter((item) => {
      return item.calleeId !== invitation.calleeId;
    });
  } else {
    clearTimeout(Store.localInvitationTime);
    if (Store.localInvitationLists.length > 0) {
      Store.localInvitationLists.map((item) => {
        item.cancel();
      });
    }
  }
};
