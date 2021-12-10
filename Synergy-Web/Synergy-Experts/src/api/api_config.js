/**
 * 配置项目权限
 * 2: 专家权限
 * 3: 管理员权限
 */
export const ProjectPermissions = 2;

/**
 * 接口配置
 */
export const PortConfig = {
  // 线上地址
  baseURL: "",
  timeout: 30000,
  withCredentials: true,
};

/**
 * loading 显示配置
 */
// 黑名单(不显示 loading)接口
export const BlankloadingPort = [
  "/teamview/insertUserOnlineInfo",
  "/teamview/joinRoom",
  "/teamview/leaveRoom",
  "/teamview/getSpecialist",
  "/teamview/getUserInfo",
];
