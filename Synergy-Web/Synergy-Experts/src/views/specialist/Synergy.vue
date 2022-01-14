<template>
  <!-- 协同页面 -->
  <el-container class="h-full overflow-hidden">
    <el-header
      class="flex justify-between items-center bg-synergy-gray_400 text-gray-200"
    >
      <span>{{ info.userName }} 的协同</span>
      <div class="flex items-center">
        <!-- 通话时长 -->
        <div class="mr-8">{{ Telephonometry.telephonometry }}</div>
        <!-- 本地音频开关 -->
        <el-button
          type="info"
          @click="audioSwitchFn"
          round
          size="small"
          class="cursor-pointer"
        >
          音频：{{ AudioSwitch ? "开" : "关" }}
        </el-button>
        <!-- 挂断 -->
        <el-button
          type="danger"
          @click="endcallFn"
          round
          size="small"
          class="cursor-pointer"
        >
          挂断
        </el-button>
      </div>
    </el-header>
    <el-main class="flex p-0 select-none overflow-hidden">
      <!-- 视频 -->
      <div class="flex-1 bg-black relative">
        <!-- 终端视频 -->
        <div
          class="w-full h-full relative flex justify-center items-center"
          id="TerminalContainer"
        >
          <!-- 添加视频 -->
          <div id="terminalVideo" class="w-full h-full"></div>
          <!-- 白板画笔 -->
          <div class="absolute z-30 w-full h-full top-0 left-0">
            <div id="whiteboardBrush" class="w-full h-full"></div>
          </div>
          <!-- 白板相关操作 @click="boardOperate" -->
          <div class="absolute bottom-9 z-30">
            <div
              :class="['board', boardType ? 'board_active' : '']"
              @click="boardOperate"
            >
              <div class="board_open"></div>
              <span class="ml-2">荧光笔</span>
            </div>
          </div>
        </div>
        <!-- 终端视频加载中 -->
        <div
          v-show="$store.state.videoLoding"
          class="absolute top-0 left-0 z-50 bg-black w-full h-full flex justify-center items-center text-synergy-gray_200 select-none"
        >
          <div class="flex flex-col items-center justify-center">
            <img
              class="animate-spin"
              src="@/assets/img/loding.svg"
              draggable="false"
              alt=""
            />
            <p class="mt-2">努力加载中</p>
          </div>
        </div>
        <!-- 终端视频离线 -->
        <div
          v-show="!$store.state.videoLoding && $store.state.terminaloffline"
          class="absolute top-0 left-0 z-50 bg-black w-full h-full flex justify-center items-center"
        >
          <div class="flex flex-col items-center">
            <img src="@/assets/img/off_line.svg" draggable="false" alt="" />
            <span class="text-synergy-gray_200 text-sm mt-2">
              终端画面离线了
            </span>
          </div>
        </div>
      </div>
      <!-- 侧边栏 -->
      <synergy-sidebar class="w-96 border-t border-gray-500" />
    </el-main>
  </el-container>
</template>

<script>
// 侧边栏
import SynergySidebar from "@/components/specialist/SynergySidebar";
// 计时
import { beginS, beginSclear } from "@/assets/untils/telephonometry.js";
// RTC 相关
import {
  JoinRTCChannel,
  CollectionAudio,
  CloseLocalAudio,
  LeaveRTCChannel,
  SetClientRole,
  // closeChannelAudio,
} from "@/anyrtc/rtc.js";
// RTM
import { ClearLocalInvitation } from "@/anyrtc/rtm.js";
// 白板相关
import { InitBoard, LaserOperation, LeaveBoard } from "@/anyrtc/board.js";
// 设置容器比例 4：3
import { containerRatio } from "@/assets/untils/until";
// 页面提示
import { ElMessageBox, ElMessage } from "element-plus";
// 路由跳转
import { useRouter, useRoute } from "vue-router";
// vuex 刷新
import { useStore } from "vuex";

import { defineComponent, reactive, ref, onMounted, onUnmounted } from "vue";
export default defineComponent({
  components: { SynergySidebar },
  setup() {
    // 操作路由
    const oRoute = useRouter();
    // 获取路由信息
    const oGetRoute = useRoute();
    // vuex
    const store = useStore();
    // 计时
    const Telephonometry = reactive({
      // 通话计时
      telephonometry: "00:00:00",
      // 通话计时定时器
      telephonometryTime: null,
    });
    // 通话计时开始/结束
    const telephonometryFn = (fase) => {
      if (fase) {
        // 计时开始
        Telephonometry.telephonometryTime = setInterval(() => {
          Telephonometry.telephonometry = beginS();
        }, 1000);
      } else {
        // 计时结束
        Telephonometry.telephonometryTime &&
          clearInterval(Telephonometry.telephonometryTime);
        beginSclear();
      }
    };

    // 音频开关
    const AudioSwitch = ref(true);
    const audioSwitchFn = () => {
      AudioSwitch.value = !AudioSwitch.value;
      CloseLocalAudio(AudioSwitch.value);
    };
    // 挂断
    const endcallFn = () => {
      ElMessageBox.confirm("是否要挂断通话，离开页面", "提示", {
        confirmButtonText: "确认",
        cancelButtonText: "取消",
        distinguishCancelAndClose: true,
      })
        .then(async () => {
          // 清除邀请
          await ClearLocalInvitation();
          await LeaveRTCChannel();
          // 清空计时
          await telephonometryFn(false);
          store.commit("emptyInvitationIng", []);
          // 返回
          oRoute.replace("/");
        })
        .catch((err) => {
          console.log(err);
          ElMessage.error("取消挂断" + JSON.stringify(err));
        });
    };

    // 激光笔
    const boardType = ref(false);
    const boardOperate = async () => {
      boardType.value = !boardType.value;
      await LaserOperation(boardType.value);
      ElMessage.success(boardType.value ? "已开启荧光笔" : "已关闭荧光笔");
    };
    onMounted(async () => {
      // 设置视频容器比例
      containerRatio("TerminalContainer", "terminalVideo");
      // 白板初始化
      await InitBoard(oGetRoute.query, "whiteboardBrush");
      // 设置角色 观众:只能订阅，不能发布音视频轨道
      await SetClientRole("主播");
      // 采集音频
      await CollectionAudio();
      // 加入房间
      await JoinRTCChannel(oGetRoute.query, "terminalVideo");
      await telephonometryFn(true);
    });
    // 页面销毁
    onUnmounted(() => {
      // 清空计时
      telephonometryFn(false);
      LeaveBoard();
      store.commit("emptyInvitationIng", []);
    });
    return {
      info: oGetRoute.query,
      // 计时
      Telephonometry,
      // 挂断
      endcallFn,
      // 音频开关
      AudioSwitch,
      audioSwitchFn,

      // 白板操作
      boardType,
      boardOperate,
    };
  },
});
</script>

<style lang="scss">
// .scrollbar-list-color {
//   .el-scrollbar__thumb {
//     @apply bg-red-500;
//   }
// }
.dit {
  @apply w-2 h-2 block rounded-full mr-3;
}
.board {
  //
  @apply py-4 cursor-pointer px-10 rounded-full bg-synergy-black_400 bg-opacity-80 font-bold text-synergy-gray_200 flex items-center justify-center;
  .board_open {
    background-image: url("../../assets/img/open.png");
    @apply bg-cover w-5 h-5;
  }
  &:hover {
    @apply bg-white text-synergy-blue_400;
    .board_open {
      background-image: url("../../assets/img/open_active.png");
      @apply bg-cover w-5 h-5;
    }
  }
}
.board_active {
  @apply text-synergy-blue_400 bg-white ring-synergy-blue_400;
  .board_open {
    background-image: url("../../assets/img/open_active.png");
    @apply bg-cover w-5 h-5;
  }
}
</style>
