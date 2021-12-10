// 获取 rtc/rtm 所需的相关信息
export const getServeInfo = () => {
  return JSON.parse(sessionStorage.getItem("SynergyLogin_UserInfo"));
};
