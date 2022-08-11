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
      <div class="flex-1 bg-black">
        <!-- 终端视频 -->
        <div
          v-show="!$store.state.terminaloffline && !$store.state.videoLoding"
          class="w-full h-full relative"
        >
          <!-- 添加视频 -->
          <div id="terminalVideo" class="w-full h-full"></div>
          <!-- 画笔标注 -->
          <!-- <brush-marks /> -->
        </div>
        <!-- 终端视频加载中 -->
        <div
          v-show="$store.state.videoLoding"
          class="w-full h-full flex justify-center items-center text-synergy-gray_200 select-none"
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
          class="w-full h-full flex justify-center items-center"
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
// 涂鸦
// import BrushMarks from "@/components/specialist/BrushMarks";
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
          ElMessage({
            message: "取消挂断" + JSON.stringify(err),
            showClose: true,
            type: "error",
          });
        });
    };

    onMounted(async () => {
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
</style>
