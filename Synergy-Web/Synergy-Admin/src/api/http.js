import axios from "axios";
import router from "@/router";
// 接口动画
import { ElLoading, ElMessage } from "element-plus";
// 接口相关配置
import { PortConfig, BlankloadingPort } from "./api_config";
var oLoading = null;
// 创建axios
const oAxios = axios.create(PortConfig);
// 添加请求拦截器
oAxios.interceptors.request.use(
  function (config) {
    // 在发送请求之前做些什么
    // 显示 loading
    oLoading && oLoading.close();
    const oM = BlankloadingPort.filter((item) => {
      return item == config.url;
    });
    if (oM.length == 0) {
      // 显示loading
      oLoading = ElLoading.service({
        lock: true,
        text: "加载中，请稍后",
        spinner: "el-icon-loading",
      });
    }
    // 部分接口需要token
    let token = sessionStorage.getItem("token");
    if (token) {
      config.headers = {
        Authorization: "Bearer " + token,
      };
    }
    return config;
  },
  function (error) {
    if (error.message.includes("timeout")) {
      // 判断请求异常信息中是否含有超时timeout字符串
      ElMessage.error("请求超时，请稍后再试'");
      return Promise.reject(function () {}); // reject这个错误信息
    } else if (error.response.status == 401) {
      // router.push("/signin");
      return Promise.reject(function () {}); // reject这个错误信息
    } else {
      ElMessage.error("网络连接失败，请稍后再试");
      return Promise.reject(error);
    }
  }
);

// 添加响应拦截器
oAxios.interceptors.response.use(
  function (response) {
    // 对响应数据处理
    oLoading && oLoading.close();
    if (
      response.headers["artc-token"] &&
      !sessionStorage.getItem("artc-token")
    ) {
      sessionStorage.setItem("token", response.headers["artc-token"]);
    }
    return response;
  },
  function (error) {
    oLoading && oLoading.close();
    if (error.message.includes("timeout")) {
      // 判断请求异常信息中是否含有超时timeout字符串
      ElMessage.error("请求超时，请稍后再试'");
      return Promise.reject({
        code: error.response.status,
        msg: error.response.data,
      });
    } else if (error.response.status === 401) {
      router.push("/signin");
      return Promise.reject({
        code: error.response.status,
        msg: error.response.data,
      });
    } else {
      ElMessage.error("网络连接失败，请稍后再试");
      return Promise.reject(error);
    }
  }
);

// 封装 get
export function get(url, params) {
  return new Promise((resolve, reject) => {
    oAxios
      .get(url, {
        params: params,
      })
      .then((res) => {
        resolve(res.data);
      })
      .catch((err) => {
        reject(err.data);
      });
  });
}
// 封装 post
export function post(url, params) {
  return new Promise((resolve, reject) => {
    oAxios
      .post(url, params)
      .then((res) => {
        resolve(res.data);
      })
      .catch((err) => {
        reject(err);
      });
  });
}
