<template>
  <el-container class="h-full">
    <el-header class="p-0">
      <Heater />
    </el-header>
    <el-main class="bg-gray-100">
      <div ref="containerHeight" class="container m-auto bg-white p-5 h-full">
        <!-- 管理员端出现 -->
        <div class="flex justify-between pb-6">
          <div class="flex items-center">
            <!-- 刷新列表 -->
            <el-button type="primary" size="mini">
              <div
                class="flex items-center"
                @click="$store.commit('upDataRefresh', true)"
              >
                刷新
                <img
                  :class="[
                    'ml-2',
                    $store.state.indexRefresh ? 'animate-spin' : '',
                  ]"
                  src="@/assets/img/refresh.svg"
                  alt=""
                />
              </div>
            </el-button>
          </div>
          <div v-if="administration == 3">
            <el-button-group>
              <el-button
                round
                :type="adminSeletType == 0 ? 'primary' : 'info'"
                @click="adminSeletType = 0"
                :class="[
                  'font-bold text-sm',
                  adminSeletType == 0
                    ? ''
                    : 'text-gray-900 bg-gray-200 border-gray-200 hover:bg-gray-100 hover:border-gray-100 hover:text-gray-900',
                ]"
              >
                进行中的任务
              </el-button>
              <el-button
                round
                :type="adminSeletType == 1 ? 'primary' : 'info'"
                @click="adminSeletType = 1"
                :class="[
                  'font-bold text-sm',
                  adminSeletType == 1
                    ? ''
                    : 'text-gray-900 bg-gray-200 border-gray-200 hover:bg-gray-100 hover:border-gray-100 hover:text-gray-900',
                ]"
              >
                已结束任务
              </el-button>
            </el-button-group>
          </div>
          <div></div>
        </div>
        <!-- table 列表 -->
        <specia-list :setHeight="oSetHeight" v-if="adminSeletType == 0" />
        <admin-list v-else />
      </div>
    </el-main>
  </el-container>
</template>

<script>
import Heater from "@/components/Header.vue";
// 专家端
import SpeciaList from "@/components/specialist/SpeciaList.vue";
// 管理员端
import AdminList from "@/components/administrator/AdminList.vue";
import { useRoute } from "vue-router";
import { defineComponent, ref, onMounted } from "vue";
export default defineComponent({
  components: { Heater, SpeciaList, AdminList },
  setup() {
     // 获取路由信息
    const oGetRoute = useRoute();
    // 获取用户信息
    const UserInfo = JSON.parse(
      sessionStorage.getItem("SynergyLogin_UserInfo")
    );
    // 管理权限
    const administration = ref(UserInfo.userType);
    // 管理员权限选项
    const adminSeletType = ref(0);
    // 刷新
    const refresh = ref(false);
    if (oGetRoute.params.type) {
      adminSeletType.value = oGetRoute.params.type;
    }
    // 获取高度
    const containerHeight = ref(null);
    const oSetHeight = ref(0);
    onMounted(() => {
      oSetHeight.value = containerHeight.value.offsetHeight - 60 - 76 - 60;
    });

    return {
      administration,
      adminSeletType,
      refresh,
      containerHeight,
      oSetHeight,
    };
  },
});
</script>
