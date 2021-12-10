<template>
  <!-- 列表 (已结束任务) -->
  <div class="min-h-96">
    <!-- 标题 -->
    <div
      class="bg-gray-100 flex text-sm text-gray-500 text-center h-12 items-center rounded"
    >
      <div class="w-1/6">房间名称</div>
      <div class="w-1/6">录制时长</div>
      <div class="w-1/6">开始时间</div>
      <div class="w-1/6">结束时间</div>
      <div class="w-1/6">状态</div>
      <div class="w-1/6">操作</div>
    </div>
    <!-- 内容 -->
    <ul v-if="AdminDatas.length > 0" class="text-center text-sm">
      <li v-for="list in AdminDatas" :key="list" class="list">
        <div class="w-1/6 truncate">{{ list.roomName }}</div>
        <div class="w-1/6 truncate">
          {{
            TimeTransitionFormatOne((list.roomStopTs - list.roomStartTs) / 1000)
          }}
        </div>
        <div class="w-1/6 truncate">
          {{ TimeTransitionFormatTwo(list.roomStartTs) }}
        </div>
        <div class="w-1/6 truncate">
          {{ TimeTransitionFormatTwo(list.roomStopTs) }}
        </div>
        <div class="w-1/6 truncate">
          <span v-if="list.roomState == 2"> 进行中 </span>
          <span v-else-if="list.roomState == 3" class="text-synergy-yellow_400">
            转码中
          </span>
          <span v-else class="text-synergy-green_400">已完成</span>
        </div>
        <div class="w-1/6 truncate flex justify-center">
          <!-- 播放 -->
          <!-- <el-tooltip content="播放" placement="top" effect="light"> -->
          <!--  v-show="list.roomState == 1 
            " -->
          <div
            v-if="list.roomState == 1 && list.roomFileUrl"
            class="play"
            title="播放"
            @click="playFn(list)"
          ></div>
          <!-- </el-tooltip> -->

          <!-- 播放转码中 -->
          <!--   v-show="
              list.roomState != 1
            " -->
          <img
            v-else
            class="opacity-30 cursor-not-allowed"
            src="@/assets/img/play.svg"
            draggable="false"
            alt=""
          />
          <!-- 下载 -->
          <div
            v-if="list.roomState == 1 && list.roomFileUrl"
            @click="downloadVideo(list)"
            class="download"
            :title="'下载' + list.roomFileUrl"
          ></div>
          <img
            v-else
            class="opacity-30 cursor-not-allowed ml-6"
            src="@/assets/img/download.svg"
            draggable="false"
            alt=""
          />
        </div>
      </li>
    </ul>
    <div v-else class="flex items-center justify-center min-h-96">
      —— 暂无结束房间 ——
    </div>
  </div>
  <div class="mt-5 flex justify-center">
    <el-pagination
      :page-sizes="[10, 20, 30, 40]"
      :page-size="Pagination.pageSize"
      :total="Pagination.totalNum"
      :hide-on-single-page="true"
      background
      layout="sizes, prev, pager, next"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    >
    </el-pagination>
  </div>
</template>

<script>
import {
  getRoomList, // 获取房间列表
} from "@/api/home.js";
// RTC 相关
import { LeaveRTCChannel } from "@/anyrtc/rtc.js";
// 时间转换
import {
  TimeTransitionFormatTwo,
  TimeTransitionFormatOne,
} from "@/assets/untils/time.js";
import { downloadUrl } from "@/assets/untils/download.js";
// 页面提示 ElMessage
import { ElMessageBox } from "element-plus";
// 路由跳转
import { useRouter } from "vue-router";
// vuex 刷新
import { useStore } from "vuex";
import { defineComponent, ref, reactive, onMounted, watch } from "vue";
export default defineComponent({
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
      getAdminDatas();
    };
    // 设置单页显示条数
    const handleSizeChange = (val) => {
      Pagination.pageSize = val;
      getAdminDatas();
    };

    // 获取列表数据
    const AdminDatas = ref([]);
    const getAdminDatas = async () => {
      const { code, data } = await getRoomList(
        Object.assign(Pagination, {
          roomState: 1, // 房间状态 结束
        })
      );
      if (code == 0) {
        if (data.list.length > 0) {
          AdminDatas.value = data.list;

          Pagination.totalNum = data.totalNum;
        } else {
          AdminDatas.value = [];
          Pagination.totalNum = 0;
        }
      }
    };

    // 播放视频
    const playFn = (data) => {
       oRoute.push({
        path: "/videoplay",
        query: {
          roomName: data.roomName,
          url: data.roomFileUrl,
        },
      });
      // const newHref = oRoute.resolve({
      //   path: "/videoplay",
      //   query: {
      //     roomName: data.roomName,
      //     url: data.roomFileUrl,
      //   },
      // });
      // window.open(newHref.href, "_blank");
    };

    // 下载视频
    const downloadVideo = (row) => {
      console.log(row);
      // 二次确定
      ElMessageBox.confirm(`确认要下载${row.roomName}回放吗？`, "提示", {
        distinguishCancelAndClose: true,
        confirmButtonText: "确认",
        cancelButtonText: "取消",
        // type: "warning",
        // showClose: false,
        // center: true,
      })
        .then(() => {
          downloadUrl(row.roomFileUrl, "视频");
        })
        .catch(() => {});
    };

    // 刷新列表
    watch(
      () => store.state.indexRefresh,
      (newValue) => {
        if (newValue) {
          getAdminDatas();
          store.commit("upDataRefresh", false);
        }
      }
    );

    onMounted(() => {
      LeaveRTCChannel();
      getAdminDatas();
    });
    return {
      AdminDatas, // 列表数据

      Pagination, // 分页相关数据
      handleCurrentChange,
      handleSizeChange,
      TimeTransitionFormatTwo, // 时间格式转换
      TimeTransitionFormatOne,

      playFn, // 播放

      downloadVideo, // 下载
    };
  },
});
</script>

<style lang="scss" scoped>
.list {
  @apply h-11 border flex items-center my-1.5 cursor-pointer;
  // &:hover {
  //   @apply bg-purple-50 text-synergy-blue_400;
  // }
}

.play {
  background-image: url("../../assets/img/play.svg");
  @apply w-5 h-5 bg-cover;
  &:hover {
    background-image: url("../../assets/img/play_hover.svg");
  }
}
.download {
  background-image: url("../../assets/img/download.svg");
  @apply w-5 h-5 bg-cover ml-6;
  &:hover {
    background-image: url("../../assets/img/download_hover.svg");
  }
}
</style>
