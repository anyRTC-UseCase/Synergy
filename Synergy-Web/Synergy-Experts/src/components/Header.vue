<template>
  <div
    class="flex items-center justify-between h-full bg-gradient-to-r from-synergy-blue_400 to-synergy-blue_300 pr-8"
  >
    <!-- <h2 class="text-2xl font-bold">智慧协同平台</h2> -->
    <img src="@/assets/img/logo_header.svg" draggable="false" alt="" />
    <el-dropdown class="text-white">
      <span class="el-dropdown-link hover:text-blue-300 flex items-center">
        {{ userInfo.userName }}
        <img class="ml-2" src="@/assets/img/id.svg" draggable="false" alt="" />
      </span>
      <template #dropdown>
        <el-dropdown-menu class="relative m-0">
          <el-dropdown-item @click="logout">退出</el-dropdown-item>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script>
import { ElMessage, ElLoading } from "element-plus";
// 路由跳转
// RTC 相关
import { LeaveRTCChannel, DestroyRTC } from "@/anyrtc/rtc.js";
import { LogoutRTM } from "@/anyrtc/rtm.js";
import { useRouter } from "vue-router";
import { defineComponent, ref } from "vue";
export default defineComponent({
  setup() {
    // 路由
    const oRoute = useRouter();
    // 用户信息
    const userInfo = ref(
      JSON.parse(sessionStorage.getItem("SynergyLogin_UserInfo"))
    );
    // 退出
    const logout = async () => {
      const loadingInstance = ElLoading.service({
        text: "退出中",
      });
      await LeaveRTCChannel();
      await DestroyRTC();
      await LogoutRTM();
      // 清空 sessionStorage
      await sessionStorage.clear();
      ElMessage.closeAll();
      loadingInstance.close();
      // 返回登录
      await oRoute.push("/signin");
    };
    return {
      userInfo,
      logout,
    };
  },
});
</script>

<style></style>
