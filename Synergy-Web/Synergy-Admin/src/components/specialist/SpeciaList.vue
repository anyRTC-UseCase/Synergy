<template>
  <div class="">
    <!-- 列表 (进行中的任务)-->

    <div class="flex">
      <div class="flex-1 mr-2 flex flex-col overflow-hidden">
        <!-- 标题 -->
        <div
          class="bg-gray-100 flex text-sm text-gray-500 text-center h-12 items-center rounded"
        >
          <div class="w-1/3">房间ID</div>
          <div class="w-1/3">房间名称</div>
          <div class="w-1/3">开始时间</div>
        </div>
        <!-- 内容 -->
        <el-scrollbar
          v-if="DataLists.SpeciaDatas.length > 0"
          class="text-center text-sm"
          :height="
            DataLists.SpeciaDatas.length * 46 > setHeight
              ? setHeight + 'px'
              : ''
          "
        >
          <div
            v-for="list in DataLists.SpeciaDatas"
            :key="list"
            :class="[
              'list',
              DataLists.SelectSpeciaDatas.roomId == list.roomId
                ? 'list_active'
                : '',
            ]"
            @click="selectDataLists(list)"
          >
            <div class="w-1/3 truncate">{{ list.roomId }}</div>
            <div class="w-1/3 truncate">{{ list.roomName }}</div>
            <div class="w-1/3 truncate">
              {{ TimeTransitionFormatTwo(list.roomTs) }}
            </div>
          </div>
        </el-scrollbar>
        <ul class="text-center text-sm" v-else>
          <li class="h-11 border flex items-center my-1.5 justify-center">
            —— 暂无房间 ——
          </li>
        </ul>
      </div>

      <div
        :class="[isFullscreen ? 'full_screen' : 'isfull_screen']"
        id="screenfull"
      >
        <h2 v-show="DataLists.SelectSpeciaDatas.roomName" class="rounded z-50">
          <!-- 相关操作 -->
          <div v-if="controlDisplay" class="flex">
            <!-- 放大 -->
            <img
              @click="screenFn"
              v-show="!isFullscreen"
              src="@/assets/img/full_screen.svg"
              draggable="false"
              alt=""
              title="全屏"
              class="mr-4 cursor-pointer"
            />
            <img
              @click="screenFn"
              v-show="isFullscreen"
              src="@/assets/img/isfull_screen.svg"
              draggable="false"
              alt=""
              title="取消"
              class="mr-4 cursor-pointer"
            />
            <!-- 声音播放 -->
            <img
              v-show="audio"
              @click="audioOperation"
              src="@/assets/img/audio.svg"
              draggable="false"
              alt=""
              title=""
              class="cursor-pointer"
            />
            <!-- 声音关闭 -->
            <img
              v-show="!audio"
              @click="audioOperation"
              src="@/assets/img/audio_close.svg"
              draggable="false"
              alt=""
              title="静音"
              class="cursor-pointer"
            />
          </div>
          <div v-else></div>
          <div class="truncate text-white">
            当前房间:
            {{ DataLists.SelectSpeciaDatas.roomName }}
          </div>
          <el-button
            v-show="isSynergy != 3"
            type="info"
            @click="goSynergy"
            size="mini"
          >
            进入协同<i class="el-icon-arrow-right el-icon--right"></i>
          </el-button>
          <div v-show="isSynergy == 3"></div>
        </h2>
        <div class="video">
          <!-- <div class="text-synergy-gray_200 select-none flex-1"> -->
          <!-- 加载中 -->
          <div
            v-show="
              $store.state.videoLoding && DataLists.SelectSpeciaDatas.roomName
            "
            class="flex flex-col items-center justify-center"
          >
            <img
              v-show="controlDisplay"
              class="animate-spin"
              src="@/assets/img/loding.svg"
              draggable="false"
              alt=""
            />
            <p class="mt-2">
              {{ controlDisplay ? "努力加载中" : "加入房间失败" }}
            </p>
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
          <!-- 无房间 -->
          <div v-show="!DataLists.SelectSpeciaDatas.roomName">
            <p>无房间</p>
          </div>
          <div
            v-show="
              !$store.state.videoLoding && DataLists.SelectSpeciaDatas.roomName
            "
            id="monitorvideo"
            class="flex-1 h-full"
          ></div>
          <!-- </div> -->
        </div>
      </div>
    </div>
    <!-- 分页 -->
    <div class="mt-4">
      <el-pagination
        @current-change="handleCurrentChange"
        @size-change="handleSizeChange"
        background
        layout="sizes,prev, pager, next"
        :page-sizes="[10, 20, 30, 40]"
        :page-size="Pagination.pageSize"
        :total="Pagination.totalNum"
        :hide-on-single-page="true"
      >
      </el-pagination>
    </div>
  </div>
</template>

<script>
// 全屏相关
import screenfull from "screenfull";
// RTC 相关
import {
  JoinRTCChannel,
  LeaveRTCChannel,
  SetClientRole,
  closeChannelAudio,
} from "@/anyrtc/rtc.js";
import {
  getRoomList, // 获取房间列表
} from "@/api/home.js";
// 时间转换
import { TimeTransitionFormatTwo } from "@/assets/untils/time.js";
// 节流
import { throttle } from "@/assets/untils/until.js";
// 页面提示
import { ElMessage, ElLoading } from "element-plus";
// 路由跳转
import { useRouter } from "vue-router";
// vuex 刷新
import { useStore } from "vuex";
import {
  defineComponent,
  reactive,
  ref,
  onMounted,
  onUnmounted,
  watch,
} from "vue";
export default defineComponent({
  props: {
    setHeight: Number,
  },
  setup() {
    // 路由
    const oRoute = useRouter();
    // vuex
    const store = useStore();

    // 定义分页
    const Pagination = reactive({
      totalNum: 10, // 总数
      pageNum: 1, // 页码
      pageSize: 10, // 单页显示条数
    });
    const handleCurrentChange = (val) => {
      Pagination.pageNum = val;
      getSpeciaDatas();
    };
    // 设置单页显示条数
    const handleSizeChange = (val) => {
      Pagination.pageSize = val;
      getSpeciaDatas();
    };

    // RTC 加入房间
    const RTCJoinChannel = async () => {
      // 离开房间
      await LeaveRTCChannel();
      // 设置角色 观众:只能订阅，不能发布音视频轨道
      await SetClientRole("观众");
      await JoinRTCChannel(DataLists.SelectSpeciaDatas, "monitorvideo");
    };

    // 定义列表相关数据
    const DataLists = reactive({
      SelectSpeciaDatas: {},
      SpeciaDatas: [],
    });
    // 请求列表方法
    const getSpeciaDatas = async () => {
      const { code, data } = await getRoomList(
        Object.assign(Pagination, {
          roomState: 2, // 房间状态 进行中
        })
      );
      if (code == 0) {
        if (data) {
          // 分页总条数
          Pagination.totalNum = data.totalNum;
          // 展示列表
          DataLists.SpeciaDatas = data.list;
          if (data.list.length > 0) {
            DataLists.SelectSpeciaDatas = data.list[0];
            RTCJoinChannel();
          } else {
            DataLists.SelectSpeciaDatas = {};
          }
        } else {
          Pagination.totalNum = 0;
          // 展示列表
          DataLists.SpeciaDatas = [];
          DataLists.SelectSpeciaDatas = {};
        }
      }
    };
    // 刷新列表
    watch(
      () => store.state.indexRefresh,
      (newValue) => {
        if (newValue) {
          getSpeciaDatas();
          store.commit("upDataRefresh", false);
        }
      }
    );
    const controlDisplay = ref(true);
    // 控制功能按钮显示
    watch(
      () => store.state.isJoinChannel,
      (newValue) => {
        controlDisplay.value = newValue;
        audio.value = false;
      }
    );
    // 选择列表
    const selectDataLists = throttle(function (data) {
      DataLists.SelectSpeciaDatas = data;
      RTCJoinChannel();
    }, 500);

    // 是否全屏
    const isFullscreen = ref(false);
    // 放大缩小操作
    const screenFn = () => {
      if (!screenfull.isEnabled) {
        ElMessage({
          message: "Your browser does not support!",
          type: "warning",
        });
        return false;
      }
      const oScreenfull = document.getElementById("screenfull");
      screenfull.toggle(oScreenfull);
    };

    //  音频操作
    const audio = ref(false);
    // 音频播放/关闭
    const audioOperation = () => {
      audio.value = !audio.value;
      closeChannelAudio(audio.value);
    };

    // 协同权限
    const isSynergy = ref(
      JSON.parse(sessionStorage.getItem("SynergyLogin_UserInfo")).userType
    );
    // 跳转到协同
    const goSynergy = async () => {
      const loadingInstance = ElLoading.service({
        text: "跳转协同中",
      });
      // rtc 离开当前频道
      // 离开房间
      await LeaveRTCChannel();
      loadingInstance.close();
      oRoute.push({
        path: "/synergy",
        query: DataLists.SelectSpeciaDatas,
        replace: true,
      });
    };

    onMounted(async () => {
      // 初始请求列表
      await getSpeciaDatas();
      // 全屏变化监听
      screenfull.isEnabled &&
        screenfull.on("change", () => {
          isFullscreen.value = screenfull.isFullscreen;
        });
    });
    // 销毁
    onUnmounted(() => {
      screenfull.isEnabled &&
        screenfull.off("change", () => {
          isFullscreen.value = screenfull.isFullscreen;
        });
    });
    return {
      DataLists, // 列表相关数据
      selectDataLists, // 选中列表

      Pagination, // 分页相关数据
      handleCurrentChange,
      handleSizeChange,
      TimeTransitionFormatTwo, // 时间格式转换

      isFullscreen, // 是否全屏
      screenFn, // 全屏方法

      audio, // 音频操作
      audioOperation, // 音频操作方法

      isSynergy, // 协同权限
      goSynergy, // 跳转到协同

      controlDisplay, // 控制功能按钮显示
    };
  },
});
</script>

<style lang="scss" scoped>
.list {
  @apply h-11 border flex items-center my-1.5 cursor-pointer;
  &:hover {
    @apply bg-purple-50 text-synergy-blue_400;
  }
}
.list_active {
  @apply bg-purple-50 text-synergy-blue_400 border-synergy-blue_400 relative;
  &::before {
    content: "";
    border-color: transparent transparent transparent #3f72ff;
    @apply block border-8 absolute left-0;
  }
  &::after {
    content: "";
    border-color: transparent #3f72ff transparent transparent;
    @apply block border-8 absolute right-0;
  }
}

// 全屏
.full_screen {
  @apply flex flex-col items-center bg-synergy-gray_400 relative;
  > h2 {
    @apply flex justify-between items-center h-12 absolute bottom-10 px-8 w-6/12 bg-gray-900 bg-opacity-70;
  }
}
// 非全屏
.isfull_screen {
  @apply flex-1 flex flex-col min-h-96 bg-synergy-gray_400 p-5 rounded;
  > h2 {
    @apply flex justify-between items-center h-12 pb-4;
  }
}

// 房间
.channel {
  @apply flex flex-col justify-center items-center;
}

.video {
  @apply flex-1 w-full flex justify-center items-center rounded text-synergy-gray_200 select-none;
}
</style>
