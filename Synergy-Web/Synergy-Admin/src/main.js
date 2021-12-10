import { createApp } from "vue";
import ElementPlus from "element-plus";
import zhCn from "element-plus/es/locale/lang/zh-cn";
import "../theme/index.css";
import "vue-video-player/src/custom-theme.css";
import "video.js/dist/video-js.css";
import "./assets/css/index.css";
// import "tailwindcss/tailwind.css";
import App from "./App.vue";
import router from "./router";
import store from "./store";
// import "@/assets/css/element-variables.scss";
// import "element-plus/theme-chalk/src/index.scss";
import VideoPlayer from "vue-video-player";

import "videojs-contrib-hls";

createApp(App)
  .use(VideoPlayer)
  .use(store)
  .use(router)
  .use(ElementPlus, {
    locale: zhCn,
  })
  .mount("#app");
