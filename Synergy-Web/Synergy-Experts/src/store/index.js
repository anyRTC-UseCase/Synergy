import { createStore } from "vuex";

export default createStore({
  state: {
    // 首页列表刷新
    indexRefresh: false,
    // 本地网络状态
    localnetwork: true,
    // 视频渲染状态
    videoLoding: false,
    // 判断是否加入房间
    isJoinChannel: false,
    // 终端离线
    terminaloffline: false,
    // 邀请中
    invitationIng: [],
  },
  mutations: {
    // 更新首页列表刷新
    upDataRefresh(state, data) {
      state.indexRefresh = data;
    },
    // 更新本地网络状态
    upDataNetwork(state, data) {
      state.localnetwork = data;
    },
    // 更新视频渲染状态
    upDataVideoLoding(state, data) {
      state.videoLoding = data;
    },
    // 更新加入房间状态
    upDataIsJoinChannel(state, data) {
      state.isJoinChannel = data;
    },
    // 设置终端状态
    setTerminalOffline(state, data) {
      state.terminaloffline = data;
    },
    // 更新邀请中状态
    upDataInvitationIng(state, data) {
      if (state.invitationIng.length > 0) {
        const oM = state.invitationIng.filter((item) => {
          return item.uid == data.uid;
        });
        if (oM.length == 0) {
          state.invitationIng.push(data);
        } else {
          state.invitationIng.map((item) => {
            if (item.uid == data.uid) {
              item = Object.assign(item, data);
            }
          });
        }
      } else {
        state.invitationIng.push(data);
      }
    },
    // 清除邀请中状态
    clearInvitationIng(state, data) {
      const oM = [];
      state.invitationIng.map((item) => {
        if (item.uid != data.uid) {
          oM.push(item);
        }
      });
      state.invitationIng = oM;
    },
    // 清空本地邀请状态
    emptyInvitationIng(state, data) {
      state.invitationIng = data;
    },
  },
  actions: {},
  modules: {},
});
