<template>
  <router-view />
  <!-- 本地网络 -->
  <local-network v-if="!onLine" />
</template>
<script>
import LocalNetwork from "@/components/LocalNetwork.vue";
import { insertUserOnlineInfo } from "@/api/home.js";
import { InitRTM } from "@/anyrtc/rtm.js";
import { InitRTC } from "@/anyrtc/rtc.js";
let HeartBeat = null;
export default {
  components: { LocalNetwork },
  data() {
    return {
      // 网络状态
      onLine: navigator.onLine,
      // rtm 是否初始化
      initrtm: true,
    };
  },
  watch: {
    "$route.path": {
      handler(newInfo) {
        // 心跳包
        HeartBeat && clearInterval(HeartBeat);
        if (newInfo && newInfo != "/signin") {
          if (newInfo != "/videoplay") {
            this.initRTMFn();
          }
          HeartBeat = setInterval(() => {
            let oInfo = JSON.parse(
              sessionStorage.getItem("SynergyLogin_UserInfo")
            );
            insertUserOnlineInfo(
              Object.assign(oInfo, {
                optTs: Date.parse(new Date()) / 1000,
              })
            );
          }, 10000);
        } else if (newInfo == "/signin") {
          this.initrtm = true;
        }
      },
      immediate: true,
      // deep: true,
    },
  },
  mounted() {
    window.addEventListener("online", this.updateOnlineStatus);
    window.addEventListener("offline", this.updateOnlineStatus);
  },
  beforeUnmount() {
    window.removeEventListener("online", this.updateOnlineStatus);
    window.removeEventListener("offline", this.updateOnlineStatus);
  },
  methods: {
    // 本地网络状态
    updateOnlineStatus(e) {
      const { type } = e;
      this.onLine = type === "online";
      this.$store.commit("upDataNetwork", this.onLine);
    },
    async initRTMFn() {
      const oInfo = await JSON.parse(
        sessionStorage.getItem("SynergyLogin_UserInfo")
      );
      if (oInfo && this.initrtm) {
        this.initrtm = false;
        InitRTM();
        InitRTC();
      }
    },
  },
};
</script>
<style lang="scss">
html,
body,
#app {
  @apply h-full relative;
}
</style>
