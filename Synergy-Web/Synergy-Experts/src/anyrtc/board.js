// 公共方法
import { getServeInfo } from "./common";
// 加载中
import { ElLoading } from "element-plus";
import ArWhiteBoard from "ar-whiteboard";
console.log("ArWhiteBoard", ArWhiteBoard);
// 白板
var Store = {
  // 白板客户端
  board: null,
};
// 白板初始化
export const InitBoard = async (info, boardId) => {
  // 获取相关信息
  Store = await Object.assign(Store, getServeInfo(), info);

  Store.board = new ArWhiteBoard({
    id: boardId,
    appId: Store.appId,
    userId: Store.uid,
    channel: Store.roomId,
    token: Store.rtmToken,
    baseParams: {
      ratio: "4:3",
    },
    // 设置私有云
    // serverParams:{
    //   ConfPriCloudAddr: {
    //     ServerAdd: '',
    //     Port: ,
    //     Wss: true,
    //   }
    // }
  });

  console.log("白板当前版本", Store.board.getVersion());

  const oLoading = ElLoading.service({ text: "加载中..." });
  // 设置初始白板
  Store.board.on("data-sync-completed", () => {
    oLoading.close();
    Store.board.setBrushType(0);
    Store.board.setGlobalBackgroundColor("transparent");
  });

  // 网络状态回调
  Store.board.on("connection-state-change", (authState, reason) => {
    console.log("网络状态回调", authState, reason);
  });

  window.onresize = function () {
    Store.board && Store.board.resize();
  };
};

// 激光笔操作
export const LaserOperation = async (fase = true) => {
  if (Store.board) {
    if (fase) {
      // 启用激光笔
      await Store.board.setBrushType(4);
      return true;
    } else {
      // 关闭激光笔
      Store.board.setBrushType(0);
      return true;
    }
  } else {
    return false;
  }
};

// 白板离开
export const LeaveBoard = () => {
  if (Store.board) {
    // 销毁白板实例
    Store.board.destroy();
    Store.board = null;
  }
};
