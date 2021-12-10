<template>
  <!--协同列表人员 -->
  <div
    class="bg-synergy-gray_400 text-gray-200 flex flex-col overflow-hidden"
    v-loading="loadingList"
    element-loading-background="rgba(0, 0, 0, 0.8)"
  >
    <h2 class="flex justify-between items-center p-5">
      <div class="text-xl font-bold">调解员（{{ Pagination.totalNum }}）</div>
      <img
        :class="['cursor-pointer', refreshLoading ? 'animate-spin' : '']"
        src="@/assets/img/refresh.svg"
        title="刷新"
        draggable="false"
        alt=""
        @click="refreshFn"
      />
    </h2>
    <div
      class="px-5 bg-synergy-gray_400 overflow-hidden"
      ref="scrollbarWarp"
      style="height: 90%"
    >
      <el-scrollbar class="scrollbar-list-color" @scroll="PullRefresh">
        <!--  :infinite-scroll-disabled="disabledScroll"   :infinite-scroll-immediate="true" -->
        <ul ref="scrollbarList">
          <li v-for="(i, index) in ExpertsList" :key="index" class="py-1 pr-3">
            <div class="flex justify-between items-center text-sm h-8">
              <div class="flex items-center flex-nowrap">
                <span
                  :class="[
                    'dit',
                    i.userState === 2 ? 'bg-synergy-green_400' : '',
                    i.userState === 0 ? 'bg-synergy-blue_400' : '',
                    i.userState === 1 ? 'bg-synergy-red_400' : '',
                    i.userState === 3 ? 'bg-synergy-gray_200' : '',
                  ]"
                ></span>
                <!-- <el-tooltip
                :content="i.userName"
                placement="left-start"
                effect="light"
              > -->
                <div
                  :class="[
                    'truncate w-64',
                    i.userState === 3 ? 'text-synergy-gray_200' : '',
                  ]"
                  :title="i.userName"
                >
                  {{ i.userName }}
                </div>
                <!-- </el-tooltip> -->
              </div>
              <!-- 状态 -->
              <div class="">
                <!-- 在线 -->
                <el-button
                  v-if="i.userState === 2"
                  @click="sendCallFn(i)"
                  type="primary"
                  size="mini"
                >
                  邀请
                </el-button>
                <!-- 邀请中 -->
                <span
                  v-else-if="i.userState === 0"
                  class="text-synergy-blue_400 h-8"
                >
                  邀请中
                </span>
                <!-- 通话中 -->
                <span
                  v-else-if="i.userState === 1"
                  class="text-synergy-red_400"
                >
                  通话中
                </span>
                <!-- 不在线 -->
                <span
                  v-else-if="i.userState === 3"
                  class="text-synergy-gray_200 h-8"
                >
                  不在线
                </span>
                <span v-else>未知</span>
              </div>
            </div>
          </li>
          <li v-show="moreLoading" class="flex items-center justify-center">
            <i class="el-icon-loading text-2xl mr-2"></i>加载中...
          </li>
          <li v-show="noMore" class="text-center">
            {{ Pagination.totalNum > 0 ? "没有更多了" : "暂无调解员" }}
          </li>
        </ul>
      </el-scrollbar>
    </div>
  </div>
</template>

<script>
// 专家列表
import { getSpecialist } from "@/api/home.js";
import { SendCall } from "@/anyrtc/rtm.js";
// 路由跳转
import { useRoute } from "vue-router";
// vuex 刷新
import { useStore } from "vuex";
import {
  defineComponent,
  ref,
  reactive,
  computed,
  onMounted,
  watch,
} from "vue";
export default defineComponent({
  setup() {
    // 获取路由信息
    const oGetRoute = useRoute();
    // vuex
    const store = useStore();

    // 定义分页
    const Pagination = reactive({
      totalNum: 10, // 总数
      pageNum: 1, // 页码
      pageSize: 20, // 单页显示条数
    });
    // 专家列表
    const ExpertsList = ref([]);
    // 列表刷新
    const loadingList = ref(false);
    // 获取专家列表
    const getExpertsList = async () => {
      Pagination.pageNum = 1;
      loadingList.value = true;
      const { code, data } = await getSpecialist(Pagination);
      loadingList.value = false;
      if (code == 0) {
        if (store.state.invitationIng.length == 0) {
          store.commit("upDataInvitationIng", data.list[0]);
        }

        ExpertsList.value = await statusReplace(
          store.state.invitationIng,
          data.list
        );
        Pagination.totalNum = data.totalNum;
      }
    };

    // 刷新调解员列表
    const refreshLoading = ref(false);
    const refreshFn = async () => {
      refreshLoading.value = true;
      await getExpertsList();
      refreshLoading.value = false;
    };

    // 下拉加载中
    const moreLoading = ref(false);
    // 滚动包裹层
    const scrollbarWarp = ref(null);
    // 滚动层
    const scrollbarList = ref(null);

    // 下拉刷新
    const PullRefresh = async ({ scrollTop }) => {
      if (ExpertsList.value.length < Pagination.totalNum) {
        // 滚动包裹层 + scrollTop = 滚动层
        if (
          scrollbarWarp.value.offsetHeight + Math.ceil(scrollTop) >=
          scrollbarList.value.offsetHeight
        ) {
          // 添加数据
          Pagination.pageNum++;
          moreLoading.value = true;
          const { code, data } = await getSpecialist(Pagination);

          if (code == 0) {
            if (data.list.length > 0) {
              ExpertsList.value = ExpertsList.value.concat(data.list);
              moreLoading.value = false;
            }
          }
        }
      }
    };
    // 计算加载完成
    const noMore = computed(() => {
      return ExpertsList.value.length >= Pagination.totalNum;
    });

    // 发送邀请
    const sendCallFn = (info) => {
      const oInfo = Object.assign(info, {
        roomName: oGetRoute.query.roomName,
        roomId: oGetRoute.query.roomId,
        roomUserName: oGetRoute.query.userName,
        userState: 0,
      });
      // // 记录存储修改为邀请中
      store.commit("upDataInvitationIng", oInfo);

      SendCall(oInfo);
    };

    // 状态数据替换
    const statusReplace = (InvitationIng, DataLists) => {
      if (InvitationIng.length > 0 && DataLists.length > 0) {
        DataLists.map((port) => {
          // 接口数据 通话中
          if (port.userState == 1) {
            // 如果接口数据中的roomId与当前页面roomId不一致
            if (port.roomId != oGetRoute.query.roomId) {
              port.userState = 2;
            }
          }
          InvitationIng.map((local) => {
            // 接口数据与本地数据 相同
            if (port.uid == local.uid) {
              // 本地数据 邀请中
              if (port.userState == 2) {
                // 接口数据更改状态为 邀请中
                if (local.userState == 0) {
                  port.userState = 0;
                }
                if (
                  local.userState == 1 &&
                  port.roomId == oGetRoute.query.roomId
                ) {
                  port.userState = 1;
                }
              }
            }
          });
        });
      }
      // 排序
      const oSort = {
        // 通话中
        callIngArray: [],
        // 邀请
        inviteArray: [],
        // 邀请中
        inviteIngArray: [],
        // 离线
        offArray: [],
      };

      DataLists.map((i) => {
        switch (i.userState) {
          case 1:
            // 通话中
            oSort.callIngArray.push(i);
            break;
          case 2:
            // 邀请
            oSort.inviteArray.push(i);
            break;
          case 0:
            // 邀请中
            oSort.inviteIngArray.push(i);
            break;
          case 3:
            // 离线
            oSort.offArray.push(i);
            break;
          default:
            break;
        }
      });
      // 组装
      const assembleArray = [
        ...oSort.inviteArray,
        ...oSort.inviteIngArray,
        ...oSort.callIngArray,
        ...oSort.offArray,
      ];
      return assembleArray;
    };
    // 监测邀请状态变换
    watch(
      () => store.state.invitationIng,
      async (newValue) => {
        if (newValue.length > 0) {
          await getExpertsList();
          ExpertsList.value = await statusReplace(newValue, ExpertsList.value);
        }
      },
      {
        deep: true, //深度监听
      }
    );
    onMounted(() => {
      Pagination.pageSize =
        Math.ceil(scrollbarWarp.value.offsetHeight / 40) + 3;
      getExpertsList();
    });

    return {
      Pagination, // 分页
      // 专家列表
      ExpertsList,
      getExpertsList,
      loadingList, // 列表刷新

      // 刷新
      refreshLoading,
      refreshFn,

      // 下拉
      PullRefresh,
      scrollbarWarp,
      scrollbarList,
      // 下拉加载中
      moreLoading,
      // 是否下拉
      // disabledScroll,
      // 加载完成
      noMore,

      // 发起邀请
      sendCallFn,
    };
  },
});
</script>

<style></style>
