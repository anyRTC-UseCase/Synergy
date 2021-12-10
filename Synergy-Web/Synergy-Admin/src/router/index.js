import { createRouter, createWebHistory } from "vue-router";
import Home from "../views/Home.vue";

import { InitRTM } from "@/anyrtc/rtm.js";
// // // RTC 初始化
import { InitRTC } from "@/anyrtc/rtc.js";

const routes = [
  {
    path: "/signin",
    name: "Signin",
    component: () => import("../views/Signin.vue"),
  },
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/home",
    redirect: "/", //重定向
  },
  // 协同页面
  {
    path: "/synergy",
    name: "Synergy",
    component: () => import("../views/specialist/Synergy.vue"),
  },
  // 播放页面
  {
    path: "/videoplay",
    name: "Videoplay",
    component: () => import("../views/administrator/Videoplay.vue"),
  },
  {
    path: "/about",
    name: "About",
    // route level code-splitting
    // this generates a separate chunk (about.[hash].js) for this route
    // which is lazy-loaded when the route is visited.
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/About.vue"),
  },
  // 没有匹配404页面
  {
    path: "/:pathMatch(.*)",
    name: "NotFound",
    component: () =>
      import(/* webpackChunkName: "about" */ "../views/NotFound.vue"),
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

/**
 * 路由守卫设置
 * to表示将要访问的路径
 * form表示从那个页面跳转而来
 * next表示允许跳转到指定位置
 *  */
router.beforeEach(async (to, from, next) => {
  // 跳转的就是登录
  if (to.path === "/signin") return next();
  // 本地token以及用户信息存在

  if (
    sessionStorage.getItem("token") &&
    sessionStorage.getItem("SynergyLogin_UserInfo")
  ) {
    if (to.path == "/" && from.path == "/signin") {
      const oInfo = await JSON.parse(
        sessionStorage.getItem("SynergyLogin_UserInfo")
      );
      if (oInfo.userType == 2) {
        await InitRTM();
      }
      await InitRTC();
    }

    next();
  } else {
    return next("/signin");
  }
});

export default router;
